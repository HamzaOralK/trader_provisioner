package main

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
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
	Config       string `json:"config"`
}

type DeletionRequest struct {
	Name string `json:"name"`
}

type Config struct {
}
