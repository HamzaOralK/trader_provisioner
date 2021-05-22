package main

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
)

type Trader struct {
	gorm.Model
	UserId       string `gorm:"primaryKey,index"`
	TraderId     string
	Config        string
}

func (u *Trader) BeforeCreate(tx *gorm.DB) (err error) {
	var temp []Trader
	var count int64
	tx.Model(&Trader{}).Where("deleted_at IS NULL and user_id = ?", u.UserId).Find(&temp).Count(&count)
	if count >= config.MaxTraderPerUser {
		err = errors.New(fmt.Sprintf("Can't save trader for user %s, it has maximum pod capacity", temp[0].UserId))
	}
	return
}

type ProvisionRequest struct {
	UserId       string `json:"userId"`
	Config        string `json:"config"`
}

type ProvisionResponse struct {
	Id string `json:"id"`
}

type DeletionRequest struct {
	UserId   string `json:"userId"`
	TraderId string `json:"traderId"`
}

type UpdateConfigRequest struct {
	UserId string `json:"userId"`
	TraderId string `json:"traderId"`
	Config string `json:"config"`
}

type Config struct {
	TraderImage string
	TraderPort  int32
	TraderPrefix string
	MaxTraderPerUser int64
}

func initializeConfig() Config {
	log.Printf("configuration initialized")
	port, _ := strconv.ParseInt(os.Getenv("TRADER_PORT"), 10, 32)
	return Config{
		TraderImage: os.Getenv("TRADER_IMAGE"),
		TraderPrefix: os.Getenv("TRADER_PREFIX"),
		TraderPort:  int32(port),
		MaxTraderPerUser: 1,
	}
}
