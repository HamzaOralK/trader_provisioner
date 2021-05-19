package main

import (
	"context"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"log"
	"net/http"
)

func ProvisionHandler(w http.ResponseWriter, r *http.Request) {
	pr := ProvisionRequest{}
	_ = json.NewDecoder(r.Body).Decode(&pr)
	deploymentId := uuid.NewV4().String()
	deploymentsClient := kubernetesClientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "trader-" + deploymentId,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"trader": "trader-" + deploymentId,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"trader": "trader-" + deploymentId,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	trader := Trader { Name: pr.Name, TraderId: deploymentId, TradingModel: pr.TradingModel}
	dbResult := db.instance.Create(&trader)
	if dbResult.Error != nil {
		log.Printf(dbResult.Error.Error())
		http.Error(w, dbResult.Error.Error(), http.StatusBadRequest)
	} else {
		log.Println("Record has been created")
		log.Println("Creating deployment...")
		result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}
		log.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	}


}