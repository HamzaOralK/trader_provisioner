package main

import (
	"context"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	capiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/utils/pointer"
	"log"
	"net/http"
)

func ProvisionHandler(w http.ResponseWriter, r *http.Request) {
	pr := ProvisionRequest{}
	_ = json.NewDecoder(r.Body).Decode(&pr)

	deploymentId := uuid.NewV4().String()
	resourceIdentifier := config.TraderPrefix + deploymentId

	deploymentsInterface, configMapInterface, serviceInterface := createClientSets()

	trader := Trader{UserId: pr.UserId, TraderId: deploymentId, Config: pr.Config}
	_ = db.instance.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&trader).Error; err != nil {
			log.Printf(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		deployment, dErr := createDeployment(resourceIdentifier, deploymentsInterface)
		configMap, cErr := createConfigMap(resourceIdentifier, pr.Config, configMapInterface)
		service, sErr := createService(resourceIdentifier, serviceInterface)
		if dErr != nil || cErr != nil || sErr != nil {
			deleteAll(resourceIdentifier, deploymentsInterface, configMapInterface, serviceInterface)
		}
		log.Printf("deployment %q with config map %q and service %q has been created", deployment.GetObjectMeta().GetName(), configMap.GetObjectMeta().GetName(), service.GetObjectMeta().GetName())
		response, _ := json.Marshal(ProvisionResponse{Id: deploymentId})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		return nil
	})
}

func createDeployment(resourceIdentifier string, deploymentInterface cappsv1.DeploymentInterface) (*appsv1.Deployment, error) {
	deploymentTemplate := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceIdentifier,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"trader": resourceIdentifier,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"trader": resourceIdentifier,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: config.TraderImage,
							Lifecycle: &apiv1.Lifecycle{
								PreStop: &apiv1.Handler{
									Exec: &apiv1.ExecAction{
										Command: []string {"python3", "scripts/rest_client.py", "--config", "user_data/config.json", "forcesell", "all"},
									},
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "trader-config",
									MountPath: "/freqtrade/user_data/config.json",
									SubPath:   "config.json",
								},
								{
									Name:      "strategies-pvc",
									MountPath: "/freqtrade/user_data/strategies",
								},
								{
									Name:      "notebooks-pvc",
									MountPath: "/freqtrade/user_data/notebooks",
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: config.TraderPort,
								},
							},
							Command: []string{"freqtrade"},
							Args:    []string{"trade", "--logfile", "/freqtrade/user_data/logs/freqtrade.log", "--db-url", "sqlite:////freqtrade/user_data/tradesv3.sqlite", "--config", "/freqtrade/user_data/config.json"},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "trader-config",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: resourceIdentifier,
									},
									DefaultMode: pointer.Int32Ptr(0777),
								},
							},
						},
						{
							Name: "notebooks-pvc",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: "notebooks-pvc",
								},
							},
						},
						{
							Name: "strategies-pvc",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: "strategies-pvc",
								},
							},
						},
					},
				},
			},
		},
	}
	return deploymentInterface.Create(context.TODO(), deploymentTemplate, metav1.CreateOptions{})
}

func createConfigMapTemplate(resourceIdentifier string, config string) *apiv1.ConfigMap {
	return &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceIdentifier,
		},
		Data: map[string]string{
			"config.json": config,
		},
	}
}

func createConfigMap(resourceIdentifier string, config string, configMapInterface capiv1.ConfigMapInterface) (*apiv1.ConfigMap, error) {
	configMapTemplate := createConfigMapTemplate(resourceIdentifier, config)
	return configMapInterface.Create(context.TODO(), configMapTemplate, metav1.CreateOptions{})
}

func createService(resourceIdentifier string, serviceInterface capiv1.ServiceInterface) (*apiv1.Service, error) {
	serviceTemplate :=  &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceIdentifier,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"trader": resourceIdentifier,
			},
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port:     config.TraderPort,
				},
			},
		},
	}
	return serviceInterface.Create(context.TODO(), serviceTemplate, metav1.CreateOptions{})
}

func deleteAll(resourceIdentifier string, deploymentInterface cappsv1.DeploymentInterface, configMapInterface capiv1.ConfigMapInterface, serviceInterface capiv1.ServiceInterface) {
	deletePolicy := metav1.DeletePropagationForeground
	_ = deploymentInterface.Delete(context.TODO(), resourceIdentifier, metav1.DeleteOptions{ PropagationPolicy: &deletePolicy })
	_ = configMapInterface.Delete(context.TODO(), resourceIdentifier, metav1.DeleteOptions{ PropagationPolicy: &deletePolicy })
	_ = serviceInterface.Delete(context.TODO(), resourceIdentifier, metav1.DeleteOptions{ PropagationPolicy: &deletePolicy })
}