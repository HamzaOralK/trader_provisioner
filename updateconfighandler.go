package main

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"net/http"
	"time"
)

func UpdateConfigHandler(w http.ResponseWriter, r *http.Request) {
	ucr := UpdateConfigRequest{}
	_ = json.NewDecoder(r.Body).Decode(&ucr)
	dbFindResult := db.instance.Where("user_id = ? AND trader_id = ?", ucr.UserId, ucr.TraderId)
	if dbFindResult.Error != nil {
		log.Printf(dbFindResult.Error.Error())
		http.Error(w, dbFindResult.Error.Error(), http.StatusBadRequest)
	} else {
		dbFindResult.Update("config", ucr.Config)
		resourceIdentifier := config.TraderPrefix + ucr.TraderId
		deploymentsClient, configMapClient, _ := createClientSets()
		configMapClient.Update(context.TODO(), createConfigMapTemplate(resourceIdentifier, ucr.Config), metav1.UpdateOptions{})
		data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().String())
		deploymentsClient.Patch(context.TODO(), resourceIdentifier, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	}
}
