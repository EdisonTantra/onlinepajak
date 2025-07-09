package application

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/EdisonTantra/lemonPajak/internal/core/cons"
	"github.com/EdisonTantra/lemonPajak/internal/core/domain"
	lemonPort "github.com/EdisonTantra/lemonPajak/internal/core/port"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/logat"
	"github.com/dslipak/pdf"
	"github.com/gin-gonic/gin"
)

var _ lemonPort.AppsHandler = (*Handler)(nil)

const (
	RequestParamEFakturFile = "file"
)

type Handler struct {
	svcApp lemonPort.ApplicationService
}

type HandlerOpts struct {
	SvcApp lemonPort.ApplicationService
}

func New(opts *HandlerOpts) *Handler {
	return &Handler{
		svcApp: opts.SvcApp,
	}
}

func (h *Handler) Mount(router *gin.Engine) {
	router.POST("/v1/e-faktur/validation", h.EFakturValidation())
}

func (h *Handler) EFakturValidation() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		ctx := gCtx.Request.Context()
		file, err := gCtx.FormFile(RequestParamEFakturFile)
		if err != nil {
			logat.GetLogger().Error(
				ctx,
				"error FormFile",
				cons.EventLogNameEFakturValidation,
				err,
			)
			gCtx.JSON(http.StatusBadRequest, HTTPError{
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Message: err.Error(),
			})
			return
		}

		pdfReader, err := h.getPDFReader(file)
		if err != nil {
			logat.GetLogger().Error(
				ctx,
				"error failed to read pdf",
				cons.EventLogNameEFakturValidation,
				err,
			)
			gCtx.JSON(http.StatusInternalServerError, HTTPError{
				Code:    fmt.Sprintf("%d", http.StatusInternalServerError),
				Message: err.Error(),
			})
		}

		rawRequest, err := h.convertPDFtoRequest(pdfReader)
		if err != nil {
			logat.GetLogger().Error(ctx, "error convert pdf into request", cons.EventLogNameEFakturValidation, err)
			gCtx.JSON(http.StatusInternalServerError, HTTPError{
				Code:    fmt.Sprintf("%d", http.StatusInternalServerError),
				Message: err.Error(),
			})
			return
		}

		//TODO remove debug valid request
		//rawRequest = &EFakturRequest{
		//	SellerTaxID:           "012345678012000",
		//	SellerTaxName:         "PT ABC",
		//	BuyerTaxID:            "023456789217000",
		//	BuyerTaxName:          "PT XYZ",
		//	DocumentEFakturNumber: "0700002212345678",
		//	DocumentEFakturDate:   "01/04/2022",
		//	TotalTaxBaseValue:     "15000000",
		//	TotalVATValue:         "1650000",
		//}

		validFaktur, err := h.svcApp.EFakturValidation(ctx, &domain.EFakturValidationRequest{
			SellerTaxID:           rawRequest.SellerTaxID,
			SellerTaxName:         rawRequest.SellerTaxName,
			BuyerTaxID:            rawRequest.BuyerTaxID,
			BuyerTaxName:          rawRequest.BuyerTaxName,
			DocumentEFakturNumber: rawRequest.DocumentEFakturNumber,
			DocumentEFakturDate:   rawRequest.DocumentEFakturDate,
			TotalTaxBaseValue:     rawRequest.TotalTaxBaseValue,
			TotalVATValue:         rawRequest.TotalVATValue,
		})
		if err != nil {
			logat.GetLogger().Error(ctx, "error invalid efaktur", cons.EventLogNameEFakturValidation, err)
			gCtx.JSON(http.StatusBadRequest, HTTPError{
				Code:    fmt.Sprintf("%d", http.StatusInternalServerError),
				Message: err.Error(),
			})
			return
		}

		deviations := make([]Deviations, 0)
		var validData *ValidatedData
		if validFaktur.ValidatedData != nil {
			validData = &ValidatedData{
				SellerTaxID:           validFaktur.ValidatedData.SellerTaxID,
				SellerTaxName:         validFaktur.ValidatedData.SellerTaxName,
				BuyerTaxID:            validFaktur.ValidatedData.BuyerTaxID,
				BuyerTaxName:          validFaktur.ValidatedData.BuyerTaxName,
				DocumentEFakturNumber: validFaktur.ValidatedData.DocumentEFakturNumber,
				DocumentEFakturDate:   validFaktur.ValidatedData.DocumentEFakturDate,
				TotalTaxBaseValue:     validFaktur.ValidatedData.TotalTaxBaseValue,
				TotalVATValue:         validFaktur.ValidatedData.TotalVATValue,
			}
		}

		if len(validFaktur.Deviations) != 0 {
			for _, entry := range validFaktur.Deviations {
				temp := Deviations{
					FieldName:     entry.Field,
					PdfValue:      entry.PdfValue,
					DjpAPIValue:   entry.DjpAPIValue,
					DeviationType: entry.DeviationType,
				}

				deviations = append(deviations, temp)
			}
		}

		gCtx.JSON(http.StatusOK, EFakturResponse{
			Status:  validFaktur.Status,
			Message: validFaktur.Message,
			ValidationResults: &ValidationResults{
				Deviations:    deviations,
				ValidatedData: validData,
			},
		})
		return
	}
}

func (h *Handler) getPDFReader(file *multipart.FileHeader) (*pdf.Reader, error) {
	if file == nil {
		err := errors.New("file request required")
		return nil, err
	}

	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buf.Bytes())
	pdfReader, err := pdf.NewReader(reader, int64(reader.Len()))
	if err != nil {
		return nil, err
	}

	return pdfReader, nil
}

func (h *Handler) convertPDFtoRequest(pdfReader *pdf.Reader) (*EFakturRequest, error) {
	var resultText string
	var resultList []string

	numPages := pdfReader.NumPage()
	for i := 1; i <= numPages; i++ {
		p := pdfReader.Page(i)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		resultList = make([]string, len(rows))
		for i, row := range rows {
			var rowText string
			for _, word := range row.Content {
				// each word
				rowText += fmt.Sprintf("%s", word.S)
			}
			// each line
			resultText += fmt.Sprintf("%s\n", rowText)
			resultList[i] = rowText
		}
	}

	var flagSeller, flagBuyer bool
	rawRequest := &EFakturRequest{}
	for _, r := range resultList {
		if r == "" {
			continue
		}

		split := strings.Split(r, ":")
		keywordPDF := normalizeKey(split[0])

		var valuePDF string
		if len(split) >= 2 {
			valuePDF = strings.Join(split[1:], ":")
			valuePDF = sanitizeString(valuePDF)
		}

		switch keywordPDF {
		case normalizeKey(cons.KeywordPDFDocumentNumber):
			rawRequest.DocumentEFakturNumber = valuePDF
		case normalizeKey(cons.KeywordPDFSeller):
			flagSeller = true
			flagBuyer = false
		case normalizeKey(cons.KeywordPDFBuyer):
			flagBuyer = true
			flagSeller = false
		case normalizeKey(cons.KeywordPDFName):
			if flagSeller {
				rawRequest.SellerTaxName = valuePDF
			} else if flagBuyer {
				rawRequest.BuyerTaxName = valuePDF
			}
		case normalizeKey(cons.KeywordPDFNPWP):
			if flagSeller {
				rawRequest.SellerTaxID = valuePDF
			} else if flagBuyer {
				rawRequest.BuyerTaxID = valuePDF
			}
		case normalizeKey(cons.KeywordPDFAddress):
			if flagSeller {
				rawRequest.SellerAddress = valuePDF
			} else if flagBuyer {
				rawRequest.BuyerAddress = valuePDF
			}
		default:
			// regex format: "{place}, {date} {month-IDN} {year}"
			isDate, _ := regexp.MatchString(`^[A-Za-z]+, \d{2} [A-Za-z]+ \d{4}$`, r)
			if isDate {
				splitList := strings.Split(r, ",")
				timestamp := splitList[1]
				rawRequest.DocumentEFakturDate = sanitizeString(timestamp)
			}

			r = normalizeKey(r)
			if strings.HasPrefix(r, normalizeKey(cons.KeywordPDFVAT)) {
				nr := strings.Replace(r, normalizeKey(cons.KeywordPDFVAT), "", 1)
				rawRequest.TotalVATValue = nr
			} else if strings.HasPrefix(r, normalizeKey(cons.KeywordPDFTaxBaseValue)) {
				nr := strings.Replace(r, normalizeKey(cons.KeywordPDFTaxBaseValue), "", 1)
				rawRequest.TotalTaxBaseValue = nr
			}
		}
	}

	return rawRequest, nil
}

func sanitizeString(s string) string {
	result := strings.TrimSpace(s)
	return result
}

func normalizeKey(s string) string {
	result := strings.ToLower(s)
	result = sanitizeString(
		strings.ReplaceAll(result, " ", ""),
	)
	return result
}
