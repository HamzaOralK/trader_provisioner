package main

import (
	"github.com/Coinoner/trader_provisioner/pkg/config"
	"github.com/Coinoner/trader_provisioner/pkg/delete"
	"github.com/Coinoner/trader_provisioner/pkg/provision"
	"github.com/Coinoner/trader_provisioner/pkg/update"
	"github.com/Coinoner/trader_provisioner/pkg/version"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var err error

func init() {
	config.InitializeConfig()
}

func main() {

	var kubernetesConfig, _ = rest.InClusterConfig()
	config.ApplicationConfig.KubernetesClientSet, err = kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter()
	r.HandleFunc("/traderVersion", version.GetTraderVersion).Methods("GET")
	r.HandleFunc("/provision", provision.Handler).Methods("POST")
	r.HandleFunc("/deletion", delete.Handler).Methods("POST")
	r.HandleFunc("/updateImage", update.ImageHandler).Methods("POST")
	r.HandleFunc("/updateConfig", update.ConfigHandler).Methods("PUT")

	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"HEAD", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(":8080", handlers.CORS(originsOk, methodsOk)(r))
}
