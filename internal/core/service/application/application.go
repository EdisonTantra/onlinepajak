package application

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/EdisonTantra/lemonPajak/internal/core/cons"
	"github.com/EdisonTantra/lemonPajak/internal/core/domain"
	"github.com/EdisonTantra/lemonPajak/internal/core/port"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/logat"
	"github.com/asaskevich/govalidator"
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
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		logat.GetLogger().Error(ctx, "failed to get response from djp", cons.EventLogNameEFakturValidation, err)
		return nil, err
	}

	//TODO remove dummy approval code
	approvalCode := "527d5baf11452b2a424b8b899e549f99426cc89fe072d84cac822e58bdf8bb56"
	if req.DocumentEFakturNumber != "" {
		approvalCode = req.DocumentEFakturNumber
	}

	respDJP, err := s.djpClient.EFakturValidation(ctx, approvalCode)
	if err != nil {
		logat.GetLogger().Error(ctx, "failed to get response from djp", cons.EventLogNameEFakturValidation, err)
		return nil, err
	}

	ds, err := compareDJPData(ctx, req, respDJP)
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
			SellerTaxID:           respDJP.NpwpPenjual,
			SellerTaxName:         respDJP.NamaPenjual,
			BuyerTaxID:            respDJP.NpwpLawanTransaksi,
			BuyerTaxName:          respDJP.NamaLawanTransaksi,
			DocumentEFakturNumber: respDJP.NomorFaktur,
			DocumentEFakturDate:   respDJP.TanggalFaktur,
			TotalTaxBaseValue:     int64(respDJP.JumlahDpp),
			TotalVATValue:         int64(respDJP.JumlahPpn),
		},
	}

	return resp, nil
}

func compareDJPData(ctx context.Context, req *domain.EFakturValidationRequest, djpData *domain.EFakturDJPResponse) ([]domain.Deviations, error) {
	ds := make([]domain.Deviations, 0)

	compare := func(field, pdfValRaw, apiValRaw string) {
		pdfVal := normalize(pdfValRaw)
		apiVal := normalize(apiValRaw)

		switch field {
		case deviationFieldJumlahDpp, deviationFieldJumlahPpn:
			pdfVal = normalizeStringIDRFormat(pdfVal)
			apiVal = normalizeStringIDRFormat(apiVal)
		case deviationFieldTanggalFaktur:
			pdfVal = normalizeStringIDNDate(pdfVal)
		default:
			pdfVal = normalizeString(pdfVal)
			apiVal = normalizeString(apiVal)
		}

		switch {
		case pdfVal == "" && apiVal != "":
			ds = append(ds, domain.Deviations{
				Field:         field,
				PdfValue:      "",
				DjpAPIValue:   apiValRaw,
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
				DjpAPIValue:   apiValRaw,
				DeviationType: "mismatch",
			})
		}

		logat.GetLogger().Info(
			ctx, "compare function",
			cons.EventLogNameEFakturValidation, map[string]string{
				"field":     field,
				"pdf_value": pdfVal,
				"api_value": apiVal,
			},
		)
	}

	// Convert int to string for comparison
	intToStr := func(i int) string {
		if i == 0 {
			return ""
		}
		return fmt.Sprintf("%d", i)
	}

	// Compare fields
	compare(deviationFieldNpwpPenjual, req.SellerTaxID, djpData.NpwpPenjual)
	compare(deviationFieldNamaPenjual, req.SellerTaxName, djpData.NamaPenjual)
	compare(deviationFieldNpwpPembeli, req.BuyerTaxID, djpData.NpwpLawanTransaksi)
	compare(deviationFieldNamaPembeli, req.BuyerTaxName, djpData.NamaLawanTransaksi)
	compare(deviationFieldNomorFaktur, req.DocumentEFakturNumber, djpData.NomorFaktur)
	compare(deviationFieldTanggalFaktur, req.DocumentEFakturDate, djpData.TanggalFaktur)
	compare(deviationFieldJumlahDpp, req.TotalTaxBaseValue, intToStr(djpData.JumlahDpp))
	compare(deviationFieldJumlahPpn, req.TotalVATValue, intToStr(djpData.JumlahPpn))

	return ds, nil
}

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if isDateFormat(s) {
		t, err := time.Parse("02/01/2006", s)
		if err == nil {
			return t.Format("2006-01-02")
		}
	}
	return s
}

func isDateFormat(s string) bool {
	matched, _ := regexp.MatchString(`^\d{2}/\d{2}/\d{4}$`, s)
	return matched
}

func normalizeString(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	s = reg.ReplaceAllString(s, "")
	return s
}

func normalizeStringIDRFormat(s string) string {
	s = strings.ReplaceAll(s, ".", "")
	s = strings.Replace(s, ",00", "", 1)
	return s
}

func normalizeStringIDNDate(s string) string {
	split := strings.Split(s, " ")
	if len(split) < 3 {
		return s
	}

	date, month, year := split[0], split[1], split[2]
	dateInt, err := strconv.Atoi(date)
	if err != nil {
		return s
	}

	dateStr := fmt.Sprintf("%d", dateInt)
	if dateInt < 10 {
		dateStr = fmt.Sprintf("0%d", dateInt)
	}

	monthStr := lookupMonthIDN[month]
	return fmt.Sprintf("%s-%s-%s", year, monthStr, dateStr)
}
