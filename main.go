package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var err error
var config Config

func init() {
	config = initializeConfig()
}

func main() {
	config.db = databaseInitialize()
	_ = config.db.instance.AutoMigrate(&Trader{})

	var kubernetesConfig, _ = rest.InClusterConfig()
	config.kubernetesClientSet, err = kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter()
	r.HandleFunc("/traderVersion", GetTraderVersion).Methods("GET")
	r.HandleFunc("/provision", ProvisionHandler).Methods("POST")
	r.HandleFunc("/deletion", DeletionHandler).Methods("POST")
	r.HandleFunc("/updateImage", UpdateImageHandler).Methods("POST")
	r.HandleFunc("/updateConfig", UpdateConfigHandler).Methods("PUT")

	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"HEAD", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(":8080", handlers.CORS(originsOk, methodsOk)(r))
}
