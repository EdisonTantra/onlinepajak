package application

import (
	"errors"
	"fmt"
	"github.com/EdisonTantra/lemonPajak/internal/core/cons"
	"github.com/EdisonTantra/lemonPajak/internal/core/domain"
	lemonPort "github.com/EdisonTantra/lemonPajak/internal/core/port"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/logat"
	"github.com/gin-gonic/gin"
	"net/http"
)

var _ lemonPort.AppsHandler = (*Handler)(nil)

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
		headerContentType := gCtx.GetHeader("Content-Type")
		if headerContentType != cons.ContentTypeMultipart {
			err := errors.New("invalid content type request")
			logat.GetLogger().Error(
				ctx,
				"invalid content type",
				cons.EventLogNameEFakturValidation,
				err,
			)

			gCtx.JSON(http.StatusBadRequest, HTTPError{
				Code:    fmt.Sprintf("%d", http.StatusBadRequest),
				Message: err.Error(),
			})
			return
		}

		//TODO parse pdf to golang struct

		//TODO remove debug
		rawRequest := EFakturRequest{
			SellerTaxID:           "atest",
			SellerTaxName:         "atest",
			BuyerTaxID:            "atest",
			BuyerTaxName:          "atest",
			DocumentEFakturNumber: "atest",
			DocumentEFakturDate:   "atest",
			TotalTaxBaseValue:     "10000",
			TotalVATValue:         "10000",
		}

		//rawRequest := EFakturRequest{
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
