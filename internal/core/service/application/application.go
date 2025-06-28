package application

import (
	"context"
	"fmt"
	"github.com/EdisonTantra/lemonPajak/internal/core/cons"
	"github.com/EdisonTantra/lemonPajak/internal/core/domain"
	"github.com/EdisonTantra/lemonPajak/internal/core/port"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/logat"
	"github.com/asaskevich/govalidator"
	"regexp"
	"strings"
	"time"
)

var _ port.ApplicationService = (*Service)(nil)

type Service struct {
	djpClient port.DJPClient
}

func New(djpClient port.DJPClient) *Service {
	return &Service{
		djpClient: djpClient,
	}
}

func (s *Service) EFakturValidation(
	ctx context.Context,
	req *domain.EFakturValidationRequest,
) (*domain.EFakturValidationResponse, error) {
	//TODO implement validation
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		logat.GetLogger().Error(ctx, "failed to get response from djp", cons.EventLogNameEFakturValidation, err)
		return nil, err
	}

	//dummy approval code
	approvalCode := "527d5baf11452b2a424b8b899e549f99426cc89fe072d84cac822e58bdf8bb56"
	if req.DocumentEFakturNumber != "" {
		approvalCode = req.DocumentEFakturNumber
	}

	respDJP, err := s.djpClient.EFakturValidation(ctx, approvalCode)
	if err != nil {
		logat.GetLogger().Error(ctx, "failed to get response from djp", cons.EventLogNameEFakturValidation, err)
		return nil, err
	}

	ds, err := compareDJPData(req, respDJP)
	const (
		statusValidWDeviations = "validated_with_deviations"
		statusValidSuccess     = "validated_successfully"
		statusValidError       = "error"
	)

	status := statusValidSuccess
	if err != nil {
		status = statusValidError
	}

	if len(ds) != 0 {
		status = statusValidWDeviations
	}

	resp := &domain.EFakturValidationResponse{
		Status:     status,
		Message:    "something",
		Deviations: ds,
		ValidatedData: &domain.EFakturData{
			SellerTaxID:           req.SellerTaxID,
			SellerTaxName:         req.SellerTaxName,
			BuyerTaxID:            req.BuyerTaxID,
			BuyerTaxName:          req.BuyerTaxName,
			DocumentEFakturNumber: req.DocumentEFakturNumber,
			DocumentEFakturDate:   req.DocumentEFakturDate,
			TotalTaxBaseValue:     int64(respDJP.JumlahDpp),
			TotalVATValue:         int64(respDJP.JumlahPpn),
		},
	}

	return resp, nil
}

func compareDJPData(req *domain.EFakturValidationRequest, djpData *domain.EFakturDJPResponse) ([]domain.Deviations, error) {
	ds := make([]domain.Deviations, 0)

	// Helper to reduce repetition
	compare := func(field, pdfValRaw, apiValRaw string) {
		pdfVal := normalize(pdfValRaw)
		apiVal := normalize(apiValRaw)

		switch {
		case pdfVal == "" && apiVal != "":
			ds = append(ds, domain.Deviations{
				Field:         field,
				PdfValue:      "",
				DjpAPIValue:   apiVal,
				DeviationType: "missing_in_pdf",
			})
		case pdfVal != "" && apiVal == "":
			ds = append(ds, domain.Deviations{
				Field:         field,
				PdfValue:      pdfVal,
				DjpAPIValue:   "",
				DeviationType: "missing_in_api",
			})
		case pdfVal != apiVal:
			ds = append(ds, domain.Deviations{
				Field:         field,
				PdfValue:      pdfVal,
				DjpAPIValue:   apiVal,
				DeviationType: "mismatch",
			})
		}
	}

	// Convert int to string for comparison
	intToStr := func(i int) string {
		if i == 0 {
			return ""
		}
		return fmt.Sprintf("%d", i)
	}

	// Compare fields
	compare("sellerTaxID", req.SellerTaxID, djpData.NpwpPenjual)
	compare("sellerTaxName", req.SellerTaxName, djpData.NamaPenjual)
	compare("buyerTaxID", req.BuyerTaxID, djpData.NpwpLawanTransaksi)
	compare("buyerTaxName", req.BuyerTaxName, djpData.NamaLawanTransaksi)
	compare("documentEFakturNumber", req.DocumentEFakturNumber, djpData.NomorFaktur)
	compare("documentEFakturDate", req.DocumentEFakturDate, djpData.TanggalFaktur)
	compare("totalTaxBaseValue", req.TotalTaxBaseValue, intToStr(djpData.JumlahDpp))
	compare("totalVATValue", req.TotalVATValue, intToStr(djpData.JumlahPpn))

	return ds, nil
}

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if isDateFormat(s) {
		t, err := time.Parse("02/01/2006", s)
		if err == nil {
			return t.Format("2006-01-02") // normalized to YYYY-MM-DD
		}
	}
	return s
}

func isDateFormat(s string) bool {
	matched, _ := regexp.MatchString(`^\d{2}/\d{2}/\d{4}$`, s)
	return matched
}
