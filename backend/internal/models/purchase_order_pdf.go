package models

// PDFData holds all the data needed to generate a purchase order PDF.
// This structure is based on the cash-bill-template-golang library.
type PDFData struct {
	CompanyName    string      `json:"company_name"`
	RegNo          string      `json:"reg_no"`
	TinNo          string      `json:"tin_no"`
	MsicCode       string      `json:"msic_code"`
	CompanyAddress string      `json:"company_address"`
	CompanyPhones  string      `json:"company_phones"`
	CompanyEmail   string      `json:"company_email"`
	CustomerName   string      `json:"customer_name"`      // This will be the Vendor's name
	CustomerPhone  string      `json:"customer_phone"`     // Vendor's phone
	ReceiptNo      string      `json:"receipt_no"`         // This will be the PO Number
	ReceiptDate    string      `json:"receipt_date"`       // PO Date
	PaymentMethod  string      `json:"payment_method"`     // e.g., "30-day term"
	Items          []PDFItem   `json:"items"`
	RoundingAdj    float64     `json:"rounding_adj"`
	AmountInWords  string      `json:"amount_in_words"`    // Library can auto-generate
	BankName       string      `json:"bank_name"`
	BankAccount    string      `json:"bank_account"`
	CustomerAddress string     `json:"customer_address"`   // Vendor's address
	CustomerEmail  string      `json:"customer_email"`     // Vendor's email
}

// PDFItem represents a single item in the purchase order.
type PDFItem struct {
	Desc     string  `json:"desc"`
	Qty      float64 `json:"qty"`
	Uom      string  `json:"uom"`
	UPrice   float64 `json:"u_price"`
	Discount float64 `json:"discount"`
}
