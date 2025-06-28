package domain

import "encoding/xml"

type EFakturValidationRequest struct {
	SellerTaxID           string `json:"sellerTaxID"`
	SellerTaxName         string `json:"sellerTaxName"`
	BuyerTaxID            string `json:"buyerTaxID"`
	BuyerTaxName          string `json:"buyerTaxName"`
	DocumentEFakturNumber string `json:"documentEFakturNumber"`
	DocumentEFakturDate   string `json:"documentEFakturDate"`
	TotalTaxBaseValue     string `json:"totalTaxBaseValue"`
	TotalVATValue         string `json:"totalVATValue"`
}

type EFakturValidationResponse struct {
	Status        string       `json:"status"`
	Message       string       `json:"message"`
	Deviations    []Deviations `json:"deviations"`
	ValidatedData *EFakturData `json:"validatedData"`
}

type Deviations struct {
	Field         string `json:"field"`
	PdfValue      string `json:"pdf_value"`
	DjpAPIValue   string `json:"djp_api_value"`
	DeviationType string `json:"deviation_type"`
}

type EFakturData struct {
	SellerTaxID           string `json:"sellerTaxID"`
	SellerTaxName         string `json:"sellerTaxName"`
	BuyerTaxID            string `json:"buyerTaxID"`
	BuyerTaxName          string `json:"buyerTaxName"`
	DocumentEFakturNumber string `json:"documentEFakturNumber"`
	DocumentEFakturDate   string `json:"documentEFakturDate"`
	TotalTaxBaseValue     int64  `json:"totalTaxBaseValue"`
	TotalVATValue         int64  `json:"totalVATValue"`
}

type EFakturDJPResponse struct {
	XMLName              xml.Name        `xml:"resValidateFakturPm"`
	KdJenisTransaksi     string          `xml:"kdJenisTransaksi"`
	FgPengganti          string          `xml:"fgPengganti"`
	NomorFaktur          string          `xml:"nomorFaktur"`
	TanggalFaktur        string          `xml:"tanggalFaktur"`
	NpwpPenjual          string          `xml:"npwpPenjual"`
	NamaPenjual          string          `xml:"namaPenjual"`
	AlamatPenjual        string          `xml:"alamatPenjual"`
	NpwpLawanTransaksi   string          `xml:"npwpLawanTransaksi"`
	NamaLawanTransaksi   string          `xml:"namaLawanTransaksi"`
	AlamatLawanTransaksi string          `xml:"alamatLawanTransaksi"`
	JumlahDpp            int             `xml:"jumlahDpp"`
	JumlahPpn            int             `xml:"jumlahPpn"`
	JumlahPpnBm          int             `xml:"jumlahPpnBm"`
	StatusApproval       string          `xml:"statusApproval"`
	StatusFaktur         string          `xml:"statusFaktur"`
	Referensi            string          `xml:"referensi"`
	DetailTransaksi      DetailTransaksi `xml:"detailTransaksi"`
}

type DetailTransaksi struct {
	Nama         string `xml:"nama"`
	HargaSatuan  int    `xml:"hargaSatuan"`
	JumlahBarang int    `xml:"jumlahBarang"`
	HargaTotal   int    `xml:"hargaTotal"`
	Diskon       int    `xml:"diskon"`
	Dpp          int    `xml:"dpp"`
	Ppn          int    `xml:"ppn"`
	TarifPpnbm   int    `xml:"tarifPpnbm"`
	Ppnbm        int    `xml:"ppnbm"`
}
