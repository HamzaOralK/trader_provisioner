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
	Name         string `gorm:"primaryKey,index"`
	TraderId     string
	TradingModel string
}

func (u *Trader) BeforeCreate(tx *gorm.DB) (err error) {
	var temp = Trader{}
	tx.Model(&Trader{}).Where("deleted_at IS NULL and name = ?", u.Name).First(&temp)
	if temp != (Trader{}) {
		err = errors.New(fmt.Sprintf("Can't save trader for user %s exists with id of %s", temp.Name, temp.TraderId))
	}
	return
}

type ProvisionRequest struct {
	Name         string `json:"name"`
	TradingModel string `json:"tradingModel"`
	Config        string `json:"config"`
}

type ProvisionResponse struct {
	Id string `json:"id"`
}

type DeletionRequest struct {
	Name string `json:"name"`
}

type Config struct {
	TraderImage string
	TraderPort  int32
	TraderPrefix string
}

func initializeConfig() Config {
	log.Printf("configuration initialized")
	port, _ := strconv.ParseInt(os.Getenv("TRADER_PORT"), 10, 32)
	return Config{
		TraderImage: os.Getenv("TRADER_IMAGE"),
		TraderPrefix: os.Getenv("TRADER_PREFIX"),
		TraderPort:  int32(port),
	}
}
