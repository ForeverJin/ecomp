package main

import "time"

// Component represents an electronic component.
type Component struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Model     string    `json:"model"`
	Package   string    `json:"package"`
	Quantity  int       `json:"quantity"`
	Location  string    `json:"location"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Transaction represents an inbound or outbound transaction.
type Transaction struct {
	ID          int       `json:"id"`
	ComponentID int       `json:"component_id"`
	Type        string    `json:"type"` // "inbound" or "outbound"
	Quantity    int       `json:"quantity"`
	Operator    string    `json:"operator"`
	Reason      string    `json:"reason"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}

// TransactionWithComponent joins transaction with component info.
type TransactionWithComponent struct {
	Transaction
	CompName  string `json:"comp_name"`
	CompModel string `json:"comp_model"`
}

// Stats holds dashboard statistics.
type Stats struct {
	TotalTypes    int `json:"total_types"`
	TotalQuantity int `json:"total_quantity"`
	LowStockCount int `json:"low_stock_count"`
	OutStockCount int `json:"out_stock_count"`
}

// TransactionFilter holds query parameters for transaction history.
type TransactionFilter struct {
	Type        string `json:"type"`
	ComponentID int    `json:"component_id"`
	DateFrom    string `json:"date_from"`
	DateTo      string `json:"date_to"`
}

// LCSCItem represents a parsed LCSC order line item.
type LCSCItem struct {
	Index      int    `json:"index"`
	LCSCNumber string `json:"lcsc_number"`
	Name       string `json:"name"`
	Model      string `json:"model"`
	Package    string `json:"package"`
	Quantity   int    `json:"quantity"`
	MatchedID  int    `json:"matched_id"`
	MatchType  string `json:"match_type"` // "existing" or "new"
	CurrentQty int    `json:"current_qty"`
}

// TemplateData is passed to every HTML template.
type TemplateData struct {
	Title      string
	Configured bool
	CSRFToken  string
	ActivePage string
	Flash      string
	FlashType  string
	Data       interface{}
}

// TypeLabel returns Chinese label for transaction type.
func TypeLabel(t string) string {
	switch t {
	case "inbound":
		return "入库"
	case "outbound":
		return "出库"
	default:
		return t
	}
}
