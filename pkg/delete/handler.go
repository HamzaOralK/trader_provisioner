package delete

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Coinoner/trader_provisioner/pkg/config"
	"github.com/Coinoner/trader_provisioner/pkg/models"
	"github.com/Coinoner/trader_provisioner/pkg/utility"
	"log"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var tm models.Trader
	dr := models.DeletionRequest{}
	_ = json.NewDecoder(r.Body).Decode(&dr)
	dbFindResult := config.ApplicationConfig.GetDbInstance().Where("user_id = ? AND trader_id = ?", dr.UserId, dr.TraderId).First(&tm)
	if dbFindResult.Error != nil {
		log.Println(dbFindResult.Error.Error())
		http.Error(w, dbFindResult.Error.Error(), http.StatusBadRequest)
	} else {
		dbDeleteResult := config.ApplicationConfig.GetDbInstance().Delete(&tm)
		if dbDeleteResult.Error != nil {
			msg := fmt.Sprintf("could not delete trader for user %s, with ID of %s", tm.UserId, tm.TraderId)
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
		} else {
			log.Println("deleting deployment")
			deletePolicy := metav1.DeletePropagationForeground
			deploymentsInterface, configMapInterface, serviceInterface, ingressInterface := utility.CreateClientSets()

			if err := deploymentsInterface.Delete(context.TODO(), config.ApplicationConfig.GetTraderPrefix()+tm.TraderId, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				log.Println(err)
			}
			if err := configMapInterface.Delete(context.TODO(), config.ApplicationConfig.GetTraderPrefix()+tm.TraderId, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				log.Println(err)
			}
			if err := serviceInterface.Delete(context.TODO(), config.ApplicationConfig.GetTraderPrefix()+tm.TraderId, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				log.Println(err)
			}
			if err := ingressInterface.Delete(context.TODO(), config.ApplicationConfig.GetTraderPrefix()+tm.TraderId, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				log.Println(err)
			}

			log.Println("deleted deployment")
		}
	}

}
