package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

type Trader struct {
	gorm.Model
	Name string
}

var db *gorm.DB
var err error

func main() {

	dsn := "host=localhost user=postgres password=12345aaa dbname=trader port=55553 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	_ = db.AutoMigrate(&Trader{})

	trader := Trader { Name: "hamza" }
	db.Create(&trader)

	r := mux.NewRouter()
	r.HandleFunc("/provision", ProvisionHandler).Methods("POST")
	http.ListenAndServe(":8080", r)
}

func ProvisionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)
}
