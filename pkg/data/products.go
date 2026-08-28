package data

type Lang string

const (
	LangVI Lang = "vi"
	LangEN Lang = "en"
)

type Product struct {
	Name       string
	URL        string
	DisplayURL string
	DescVI     string
	DescEN     string
}

type CompanyInfo struct {
	NameVI   string
	NameEN   string
	TaxID    string
	Phone    string
	PhoneURL string
	Year     string
}

var Company = CompanyInfo{
	NameVI:   "CÔNG TY TNHH BANANA SOFTWARE",
	NameEN:   "BANANA SOFTWARE CO., LTD.",
	TaxID:    "2803238388",
	Phone:    "0335581402",
	PhoneURL: "https://zalo.me/0335581402",
	Year:     "2026",
}

var Products = []Product{
	{
		Name:       "TaiHoaDon",
		URL:        "https://taihoadon.online",
		DisplayURL: "taihoadon.online",
		DescVI:     "Tải hóa đơn điện tử hàng loạt, tạo bản kê tự động, tải Pdf gốc",
		DescEN:     "Batch download e-invoices, auto summary sheets, download original PDFs",
	},
	{
		Name:       "TaiToKhai",
		URL:        "https://taitokhai.online",
		DisplayURL: "taitokhai.online",
		DescVI:     "Tải hồ sơ đã nộp dichvucong và thuedientu",
		DescEN:     "Download filings from national public service & e-tax portals",
	},
	{
		Name:       "KeToanONE",
		URL:        "https://ketoan.one",
		DisplayURL: "ketoan.one",
		DescVI:     "Phần mềm kế toán",
		DescEN:     "Accounting software",
	},
	{
		Name:       "HoaDonGoc",
		URL:        "https://hoadongoc.com",
		DisplayURL: "hoadongoc.com",
		DescVI:     "Tìm và tải hóa đơn PDF gốc từ nhà cung cấp",
		DescEN:     "Lookup and download original vendor PDF invoices",
	},
	{
		Name:       "QuanLyHoaDon",
		URL:        "https://adsclone.com",
		DisplayURL: "quanlyhoadon.com",
		DescVI:     "Quản lý hóa đơn",
		DescEN:     "Invoice management",
	},
}
