package main

import (
	"context"
	"encoding/json"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func DeletionHandler(w http.ResponseWriter, r *http.Request) {
	var tm Trader
	dr := DeletionRequest{}
	_ = json.NewDecoder(r.Body).Decode(&dr)
	db.Where("name = ?", dr.Name).First(&tm)

	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	deploymentsClient := kubernetesClientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
	if err := deploymentsClient.Delete(context.TODO(), "trader-"+tm.TraderId, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
	db.Delete(&tm)
}
