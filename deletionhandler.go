package main

import (
	"context"
	"encoding/json"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
)

func DeletionHandler(w http.ResponseWriter, r *http.Request) {
	var tm Trader
	dr := DeletionRequest{}
	_ = json.NewDecoder(r.Body).Decode(&dr)
	db.instance.Where("name = ?", dr.Name).First(&tm)
	dbResult := db.instance.Delete(&tm)

	if dbResult.Error != nil {
		log.Printf("Could not delete trader for user %s, with ID of %s", tm.Name, tm.TraderId)
		http.Error(w, fmt.Sprintf("Could not delete trader for user %s, with ID of %s", tm.Name, tm.TraderId), http.StatusBadRequest)
	} else {
		log.Println("Deleting deployment...")
		deletePolicy := metav1.DeletePropagationForeground
		deploymentsClient := kubernetesClientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
		if err := deploymentsClient.Delete(context.TODO(), "trader-"+tm.TraderId, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			panic(err)
		}
		log.Println("Deleted deployment.")
	}
}
