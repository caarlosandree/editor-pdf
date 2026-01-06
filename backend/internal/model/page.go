package model

// Page representa metadados de uma página de PDF
type Page struct {
	Number int     `json:"number"`
	Width  float64 `json:"width"`  // em PDF points (72 DPI)
	Height float64 `json:"height"` // em PDF points (72 DPI)
}
