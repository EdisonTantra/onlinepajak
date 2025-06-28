package http

import (
	"bytes"
	"fmt"
	"io"
	"strings"

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

		headersRaw := c.Request.Header
		jsonData, err := io.ReadAll(c.Request.Body)
		if err != nil {

		}

		requestBodyStr := string(jsonData)
		requestBodyStr = strings.TrimSpace(requestBodyStr)

		headerStr := ""
		params := c.Request.URL.Query()
		paramList := make([]string, 0)
		for k, v := range params {
			t := fmt.Sprintf("%s: %s", k, v)
			paramList = append(paramList, t)
		}
		paramsStr := strings.Join(paramList, ",")

		for k, v := range headersRaw {
			headerStr += fmt.Sprintf("\"%s\":\"%s\", ", k, v)
		}

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
			"response_body": blw.body.String(),
		}

		r.logger.Info(ctx, "incoming request", "http_response", dataResp)
	}
}
