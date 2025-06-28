package djp

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/EdisonTantra/lemonPajak/internal/core/domain"
	"github.com/EdisonTantra/lemonPajak/internal/core/port"
	"io"
	"net/http"
	"time"
)

var _ port.DJPClient = (*Client)(nil)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) port.DJPClient {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // adjust as needed
		},
	}
}

func (c *Client) EFakturValidation(ctx context.Context, approvalCode string) (*domain.EFakturDJPResponse, error) {
	// Construct full URL
	url := fmt.Sprintf("%s/validasi/faktur/approvalCode/%s", c.baseURL, approvalCode) // adjust the path if needed

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send the request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read and decode response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	//TODO remove debug mock resp djp
	respBody = []byte(mockRespDJP)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode, string(respBody))
	}

	var efakturResp domain.EFakturDJPResponse
	if err := xml.Unmarshal(respBody, &efakturResp); err != nil {
		return nil, fmt.Errorf("failed to parse XML response: %w", err)
	}

	return &efakturResp, nil
}

var mockRespDJP = `<resValidateFakturPm>
	<kdJenisTransaksi>07</kdJenisTransaksi>
	<fgPengganti>0</fgPengganti>
	<nomorFaktur>0700002212345678</nomorFaktur>
	<tanggalFaktur>01/04/2022</tanggalFaktur>
	<npwpPenjual>012345678012000</npwpPenjual>
	<namaPenjual>PT ABC</namaPenjual>
	<alamatPenjual>Jalan Gatot Subroto No. 40A, Senayan, Kebayoran Baru,
Jakarta Selatan 12910</alamatPenjual>
	<npwpLawanTransaksi>023456789217000</npwpLawanTransaksi>
online-pajak.com

	<namaLawanTransaksi>PT XYZ</namaLawanTransaksi>
	<alamatLawanTransaksi>Jalan Kuda Laut No. 1, Sungai Jodoh, Batu Ampar,
Batam 29444</alamatLawanTransaksi>
	<jumlahDpp>15000000</jumlahDpp>
	<jumlahPpn>1650000</jumlahPpn>
	<jumlahPpnBm>0</jumlahPpnBm>
	<statusApproval>Faktur Valid, Sudah Diapprove oleh DJP</statusApproval>
	<statusFaktur>Faktur Pajak Normal</statusFaktur>
	<referensi>123/ABC/IV/2022</referensi>
	<detailTransaksi>
		<nama>KOMPUTER MERK ABC, HS Code 84714110</nama>
		<hargaSatuan>5000000</hargaSatuan>
		<jumlahBarang>3</jumlahBarang>
		<hargaTotal>15000000</hargaTotal>
		<diskon>0</diskon>
		<dpp>15000000</dpp>
		<ppn>1650000</ppn>
		<tarifPpnbm>0</tarifPpnbm>
		<ppnbm>0</ppnbm>
	</detailTransaksi>
</resValidateFakturPm>`
