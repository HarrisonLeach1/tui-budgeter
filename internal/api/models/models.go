package models

// Connection - Represents an organisation as it appears in the /connections endpoint
type Connection struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenantId"`
	TenantType     string `json:"tenantType"`
	TenantName     string `json:"tenantName"`
	CreateDateUTC  string `json:"createdDateUtc"`
	UpdatedDateUTC string `json:"updatedDateUtc"`
}

type APIException struct {
	ErrorNumber int       `json:"ErrorNumber"`
	Type        string    `json:"Type"`
	Message     string    `json:"Message"`
	Elements    []Element `json:"Elements"`
}

type Element struct {
	ValidationErrors []ValidationError `json:"ValidationErrors"`
}
type ValidationError struct {
	Message string `json:"Message"`
}

type ReportResponse struct {
	Reports []Report `json:"Reports"`
}

type Report struct {
	ReportID       string      `json:"ReportID"`
	ReportName     string      `json:"ReportName"`
	ReportType     string      `json:"ReportType"`
	ReportTitles   []string    `json:"ReportTitles"`
	ReportDate     string      `json:"ReportDate"`
	UpdatedDateUTC string      `json:"UpdatedDateUTC"`
	Rows           []ReportRow `json:"Rows"`
}

type ReportRow struct {
	RowType string       `json:"RowType"`
	Title   string       `json:"Title"`
	Cells   []ReportCell `json:"Cells"`
	Rows    []ReportRow  `json:"Rows"`
}

type ReportCell struct {
	Value      string           `json:"Value"`
	Attributes []CellAttributes `json:"Attributes"`
}

type CellAttributes struct {
	Value string `json:"Value"`
	Id    string `json:"Id"`
}
