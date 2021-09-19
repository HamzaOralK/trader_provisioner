package update

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Coinoner/trader_provisioner/pkg/config"
	"github.com/Coinoner/trader_provisioner/pkg/models"
	"github.com/Coinoner/trader_provisioner/pkg/utility"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"net/http"
)

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	uir := models.UpdateImageRequest{}
	_ = json.NewDecoder(r.Body).Decode(&uir)
	dbFindResult := config.ApplicationConfig.GetDbInstance().Where("user_id = ? AND trader_id = ?", uir.UserId, uir.TraderId)
	if dbFindResult.Error != nil {
		log.Println(dbFindResult.Error.Error())
		http.Error(w, dbFindResult.Error.Error(), http.StatusBadRequest)
	} else {
		deploymentsInterface, _, _, _ := utility.CreateClientSets()
		resourceIdentifier := config.ApplicationConfig.GetTraderPrefix() + uir.TraderId
		data := fmt.Sprintf(`{"spec":{"template":{"spec":{"containers":[{"name":"web", "image":"%s"}]}}}}`, config.ApplicationConfig.GetTraderImage())
		_, err := deploymentsInterface.Patch(context.TODO(), resourceIdentifier, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-recreate"})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}