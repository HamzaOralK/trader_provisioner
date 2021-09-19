package models

import (
	"gorm.io/gorm"
)

type Trader struct {
	gorm.Model
	UserId   string `gorm:"primaryKey,index"`
	TraderId string
	Config   string
}

type ProvisionRequest struct {
	UserId string `json:"user_id"`
	Config string `json:"configuration"`
}

type ProvisionResponse struct {
	Id string `json:"id"`
	Version string `json:"version"`
}

type DeletionRequest struct {
	UserId   string `json:"user_id"`
	TraderId string `json:"trader_id"`
}

type UpdateConfigRequest struct {
	UserId   string `json:"user_id"`
	TraderId string `json:"trader_id"`
	Config   string `json:"configuration"`
}

type UpdateImageRequest struct {
	UserId   string `json:"user_id"`
	TraderId string `json:"trader_id"`
}

type TraderVersionResponse struct {
	Version string `json:"version"`
}
