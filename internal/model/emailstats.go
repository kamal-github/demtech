package model

type EmailStats struct {
	TotalEmailsSent int            `json:"totalEmailsSent"`
	SuccessCount    int            `json:"successCount"`
	TotalErrCount   int            `json:"totalErrCount"`
	Errors          map[string]int `json:"errors,omitempty"` // typeOfErr -> Count
}
