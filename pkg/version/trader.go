package version

import (
	"encoding/json"
	"github.com/Coinoner/trader_provisioner/pkg/config"
	"github.com/Coinoner/trader_provisioner/pkg/models"
	"net/http"
	"strings"
)

func GetTraderVersion(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(models.TraderVersionResponse{
		Version:  strings.Split(config.ApplicationConfig.GetTraderImage(),":")[1],
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}