package main

import "gorm.io/gorm"

type Trader struct {
	gorm.Model
	Name string `gorm:"primaryKey"`
	TraderId string
	TradingModel string
}

type ProvisionRequest struct {
	Name string `json:"name"`
	TradingModel string `json:"tradingModel"`
}

type DeletionRequest struct {
	Name string `json:"name"`
}
