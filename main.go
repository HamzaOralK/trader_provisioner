package main

import (
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
)

var db DB
var err error
var kubernetesClientSet *kubernetes.Clientset

const traderPrefix string = "trader-"

func main() {
	db = databaseInitialize()
	_ = db.instance.AutoMigrate(&Trader{})

	var kubernetesConfig, _ = rest.InClusterConfig()
	kubernetesClientSet, err = kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter()
	r.HandleFunc("/provision", ProvisionHandler).Methods("POST")
	r.HandleFunc("/deletion", DeletionHandler).Methods("POST")
	http.ListenAndServe(":8080", r)
}
