package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
)

var db *gorm.DB
var err error
var kubernetesConfig *rest.Config
var kubernetesClientSet *kubernetes.Clientset

func main() {
	dsn := "host=postgres-postgresql user=postgres password=12345aaa dbname=trader port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	_ = db.AutoMigrate(&Trader{})

	kubernetesConfig, _ = rest.InClusterConfig()
	kubernetesClientSet, err = kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter()
	r.HandleFunc("/provision", ProvisionHandler).Methods("POST")
	r.HandleFunc("/deletion", DeletionHandler).Methods("POST")
	http.ListenAndServe(":8080", r)
}
