package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/EdisonTantra/lemonPajak/internal/core/cons"
	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// overwrite the write func and assign as new writer
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (r *Router) MiddlewareLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		c.Writer = blw
		requestBodyStr := getRequestBody(c)
		headerStr := ""

		headersRaw := c.Request.Header
		for k, v := range headersRaw {
			headerStr += fmt.Sprintf("\"%s\":\"%s\", ", k, v)
		}

		params := c.Request.URL.Query()
		paramList := make([]string, 0)
		for k, v := range params {
			t := fmt.Sprintf("%s: %s", k, v)
			paramList = append(paramList, t)
		}

		paramsStr := strings.Join(paramList, ",")
		dataRequest := map[string]string{
			"url":        c.Request.URL.String(),
			"headers":    headerStr,
			"parameters": paramsStr,
			"body":       requestBodyStr,
		}
		r.logger.Info(ctx, "incoming request", "http_request", dataRequest)
		c.Next()

		dataResp := map[string]string{
			"url":           c.Request.URL.String(),
			"http_status":   fmt.Sprintf("%d", c.Writer.Status()),
			"response_body": truncateBody(blw.body.String()),
		}

		r.logger.Info(ctx, "incoming request", "http_response", dataResp)
	}
}

var sensitiveFields = map[string]struct{}{
	"password": {},
	"token":    {},
	"secret":   {},
}

func maskSensitiveFields(data map[string]interface{}) {
	for key, val := range data {
		if _, found := sensitiveFields[strings.ToLower(key)]; found {
			data[key] = cons.MaskLogText
			continue
		}

		if nestedMap, ok := val.(map[string]interface{}); ok {
			maskSensitiveFields(nestedMap)
		}

		if arrayVal, ok := val.([]interface{}); ok {
			for _, item := range arrayVal {
				if itemMap, ok := item.(map[string]interface{}); ok {
					maskSensitiveFields(itemMap)
				}
			}
		}
	}
}

func truncateBody(s string) string {
	if len(s) > cons.MaxLengthBodyLog {
		return fmt.Sprintf("%s...[TRUNCATED]", s[:cons.MaxLengthBodyLog])
	}
	return s
}

func getRequestBody(c *gin.Context) (requestBodyStr string) {
	contentType := c.GetHeader(cons.HeaderNameContentType)
	bodyCopy, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return cons.BodyLogUnreadable
	}

	switch contentType {
	case cons.ContentTypeText:
		requestBodyStr = strings.TrimSpace(string(bodyCopy))
	case cons.ContentTypeJSON:
		var jsonMap map[string]interface{}
		if err := json.Unmarshal(bodyCopy, &jsonMap); err == nil {
			maskSensitiveFields(jsonMap)
			maskedBytes, _ := json.Marshal(jsonMap)
			requestBodyStr = string(maskedBytes)
		}
	default:
		// Check if the content type is multipart (usually file upload)
		if strings.HasPrefix(contentType, cons.ContentTypeMultipartData) {
			// 👇 Don't log file content
			// Still must restore body so Gin can read it
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))
			requestBodyStr = cons.BodyLogFileUpload
		} else {
			requestBodyStr = strings.TrimSpace(string(bodyCopy))
		}
	}

	return truncateBody(requestBodyStr)
}
