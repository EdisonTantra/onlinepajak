package djp

import "encoding/xml"

type EFakturResponse struct {
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
