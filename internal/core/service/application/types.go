package application

const (
	statusValidWDeviations = "validated_with_deviations"
	statusValidSuccess     = "validated_successfully"
	statusValidError       = "error"

	deviationFieldNpwpPenjual   = "npwpPenjual"
	deviationFieldNamaPenjual   = "namaPenjual"
	deviationFieldNpwpPembeli   = "npwpPembeli"
	deviationFieldNamaPembeli   = "namaPembeli"
	deviationFieldNomorFaktur   = "nomorFaktur"
	deviationFieldTanggalFaktur = "tanggalFaktur"
	deviationFieldJumlahDpp     = "jumlahDpp"
	deviationFieldJumlahPpn     = "jumlahPpn"
)

var (
	lookupMonthIDN = map[string]string{
		"januari":   "01",
		"febuari":   "02",
		"maret":     "03",
		"april":     "04",
		"mei":       "05",
		"juni":      "06",
		"juli":      "07",
		"agustus":   "08",
		"september": "09",
		"oktober":   "10",
		"november":  "11",
		"desember":  "12",
	}
)
