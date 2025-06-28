package application

type EFakturRequest struct {
	SellerTaxID           string `json:"sellerTaxID"`
	SellerTaxName         string `json:"sellerTaxName"`
	BuyerTaxID            string `json:"buyerTaxID"`
	BuyerTaxName          string `json:"buyerTaxName"`
	DocumentEFakturNumber string `json:"documentEFakturNumber"`
	DocumentEFakturDate   string `json:"documentEFakturDate"`
	TotalTaxBaseValue     string `json:"totalTaxBaseValue"`
	TotalVATValue         string `json:"totalVATValue"`
}

type EFakturResponse struct {
	Status            string             `json:"status"`
	Message           string             `json:"message"`
	ValidationResults *ValidationResults `json:"validation_results"`
}

type ValidationResults struct {
	Deviations    []Deviations   `json:"deviations"`
	ValidatedData *ValidatedData `json:"validated_data"`
}

type Deviations struct {
	FieldName     string `json:"field"`
	PdfValue      string `json:"pdf_value"`
	DjpAPIValue   string `json:"djp_api_value"`
	DeviationType string `json:"deviation_type"`
}

type ValidatedData struct {
	SellerTaxID           string `json:"npwpPenjual"`
	SellerTaxName         string `json:"namaPenjual"`
	BuyerTaxID            string `json:"npwpPembeli"`
	BuyerTaxName          string `json:"namaPembeli"`
	DocumentEFakturNumber string `json:"nomorFaktur"`
	DocumentEFakturDate   string `json:"tanggalFaktur"`
	TotalTaxBaseValue     int64  `json:"jumlahDpp"`
	TotalVATValue         int64  `json:"jumlahPpn"`
}

type HTTPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
