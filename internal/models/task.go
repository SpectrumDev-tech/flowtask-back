// Package models
package models

type EnergyLevel string

const (
	EnergyLow    EnergyLevel = "low"
	EnergyMedium EnergyLevel = "medium"
	EnergyHigh   EnergyLevel = "high"
)

type Task struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Energy      EnergyLevel `json:"energy"`
	Category    string      `json:"category"`
	Completed   bool        `json:"completed"`
	CreatedAt   string      `json:"created_at"`
}
