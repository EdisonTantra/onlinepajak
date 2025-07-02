package cons

const (
	EventLogNameRoot = "root"

	EventLogNameUserDetail        = "user_detail"
	EventLogNameHealth            = "health_check"
	EventLogNameEFakturValidation = "efaktur_validation"
)

const (
	MaxLengthBodyLog = 200

	HeaderNameContentType    = "Content-Type"
	ContentTypeJSON          = "application/json"
	ContentTypeText          = "text/plain"
	ContentTypeMultipartData = "multipart/form-data"

	BodyLogUnreadable = "[UNREADABLE BODY]"
	BodyLogFileUpload = "[FILE UPLOADED]"

	MaskLogText = "***"
)

const (
	KeywordPDFDocumentNumber = "Kode dan Nomor Seri Faktur Pajak"
	KeywordPDFSeller         = "Pengusaha Kena Pajak"
	KeywordPDFBuyer          = "Pembeli Barang Kena Pajak / Penerima Jasa Kena Pajak"
	KeywordPDFName           = "Nama"
	KeywordPDFAddress        = "Alamat"
	KeywordPDFNPWP           = "NPWP"
	KeywordPDFTaxBaseValue   = "Dasar Pengenaan Pajak"
	KeywordPDFVAT            = "PPN = 10% x Dasar Pengenaan Pajak"
)
