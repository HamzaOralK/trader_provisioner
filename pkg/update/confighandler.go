package update

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Coinoner/trader_provisioner/pkg/config"
	"github.com/Coinoner/trader_provisioner/pkg/models"
	"github.com/Coinoner/trader_provisioner/pkg/provision"
	"github.com/Coinoner/trader_provisioner/pkg/utility"
	"log"
	"net/http"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	ucr := models.UpdateConfigRequest{}
	_ = json.NewDecoder(r.Body).Decode(&ucr)
	dbFindResult := config.ApplicationConfig.GetDbInstance().Where("user_id = ? AND trader_id = ?", ucr.UserId, ucr.TraderId).Model(&models.Trader{})
	if dbFindResult.Error != nil {
		log.Println(dbFindResult.Error.Error())
		http.Error(w, dbFindResult.Error.Error(), http.StatusBadRequest)
	} else {
		dbFindResult.Update("config", ucr.Config)
		resourceIdentifier := config.ApplicationConfig.GetTraderPrefix() + ucr.TraderId
		deploymentsClient, configMapClient, _, _ := utility.CreateClientSets()
		configMapClient.Update(context.TODO(), provision.CreateConfigMapTemplate(resourceIdentifier, ucr.Config), metav1.UpdateOptions{})
		data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().String())
		deploymentsClient.Patch(context.TODO(), resourceIdentifier, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-recreate"})
	}
}
