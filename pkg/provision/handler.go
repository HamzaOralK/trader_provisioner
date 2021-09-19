package provision

import (
	"context"
	"encoding/json"
	"github.com/Coinoner/trader_provisioner/pkg/config"
	"github.com/Coinoner/trader_provisioner/pkg/models"
	"github.com/Coinoner/trader_provisioner/pkg/utility"
	"log"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	capiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	cnetworkingv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
	"k8s.io/utils/pointer"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	pr := models.ProvisionRequest{}
	_ = json.NewDecoder(r.Body).Decode(&pr)

	deploymentId := uuid.NewV4().String()
	resourceIdentifier := config.ApplicationConfig.GetTraderPrefix() + deploymentId

	deploymentsInterface, configMapInterface, serviceInterface, ingressInterface := utility.CreateClientSets()

	trader := models.Trader{UserId: pr.UserId, TraderId: deploymentId, Config: pr.Config}
	_ = config.ApplicationConfig.GetDbInstance().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&trader).Error; err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		deployment, dErr := createDeployment(resourceIdentifier, deploymentsInterface)
		configMap, cErr := createConfigMap(resourceIdentifier, pr.Config, configMapInterface)
		service, sErr := createService(resourceIdentifier, serviceInterface)
		ingress, iErr := createIngress(resourceIdentifier, deploymentId, ingressInterface)

		if dErr != nil || cErr != nil || sErr != nil || iErr != nil {
			deleteAll(resourceIdentifier, deploymentsInterface, configMapInterface, serviceInterface, ingressInterface)
		}
		log.Printf("deployment %q with config map %q and service %q has been created, ingress inserted into %q",
			deployment.GetObjectMeta().GetName(), configMap.GetObjectMeta().GetName(), service.GetObjectMeta().GetName(), ingress.GetObjectMeta().GetName())
		response, _ := json.Marshal(models.ProvisionResponse{
			Id: deploymentId,
			Version: strings.Split(config.ApplicationConfig.GetTraderImage(),":")[1],
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		return nil
	})
}

func createDeployment(resourceIdentifier string, deploymentInterface cappsv1.DeploymentInterface) (*appsv1.Deployment, error) {
	deploymentTemplate := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceIdentifier,
			Annotations: map[string]string{
				"application": "trader",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Strategy: appsv1.DeploymentStrategy{Type: "Recreate"},
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
					Annotations: map[string]string{
						"cluster-autoscaler.kubernetes.io/safe-to-evict": "false",
					},
				},
				Spec: apiv1.PodSpec{
					ImagePullSecrets: []apiv1.LocalObjectReference{{Name: config.ApplicationConfig.GetImagePullSecrets()}},
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: config.ApplicationConfig.GetTraderImage(),
							Lifecycle: &apiv1.Lifecycle{
								PreStop: &apiv1.Handler{
									Exec: &apiv1.ExecAction{
										Command: []string{"python3", "scripts/rest_client.py", "--config", "user_data/config.json", "forcesell", "all"},
									},
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "trader-config",
									MountPath: "/freqtrade/user_data/config.json",
									SubPath:   "config.json",
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: config.ApplicationConfig.GetTraderPort(),
								},
							},
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									apiv1.ResourceCPU:    resource.MustParse("500m"),
									apiv1.ResourceMemory: resource.MustParse("500Mi"),
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
					},
				},
			},
		},
	}
	return deploymentInterface.Create(context.TODO(), deploymentTemplate, metav1.CreateOptions{})
}

func CreateConfigMapTemplate(resourceIdentifier string, config string) *apiv1.ConfigMap {
	return &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceIdentifier,
			Annotations: map[string]string{
				"application": "trader",
			},
		},
		Data: map[string]string{
			"config.json": config,
		},
	}
}

func createConfigMap(resourceIdentifier string, config string, configMapInterface capiv1.ConfigMapInterface) (*apiv1.ConfigMap, error) {
	configMapTemplate := CreateConfigMapTemplate(resourceIdentifier, config)
	return configMapInterface.Create(context.TODO(), configMapTemplate, metav1.CreateOptions{})
}

func createService(resourceIdentifier string, serviceInterface capiv1.ServiceInterface) (*apiv1.Service, error) {
	serviceTemplate := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceIdentifier,
			Annotations: map[string]string{
				"application": "trader",
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"trader": resourceIdentifier,
			},
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port:     config.ApplicationConfig.GetTraderPort(),
				},
			},
		},
	}
	return serviceInterface.Create(context.TODO(), serviceTemplate, metav1.CreateOptions{})
}

func createIngress(resourceIdentifier string, id string, ingressInterface cnetworkingv1.IngressInterface) (*networkingv1.Ingress, error) {
	pt := networkingv1.PathTypePrefix
	path := networkingv1.HTTPIngressPath{
		Path:     "/",
		PathType: &pt,
		Backend: networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: resourceIdentifier,
				Port: networkingv1.ServiceBackendPort{
					Number: 8080,
				},
			},
		},
	}

	om := metav1.ObjectMeta{
		Name: resourceIdentifier,
		Annotations: map[string]string{
			"kubernetes.io/ingress.class":   "nginx",
			"cert-manager.io/issuer": config.ApplicationConfig.GetClusterCertificate(),
			"application":                   "trader",
		},
	}

	specRules := []networkingv1.IngressRule{
		{
			Host: id + "." + config.ApplicationConfig.GetHostname(),
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						path,
					},
				},
			},
		},
	}

	specTLS := []networkingv1.IngressTLS{
		{
			Hosts: []string{
				id + "." + config.ApplicationConfig.GetHostname(),
			},
			SecretName: resourceIdentifier,
		},
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: om,
		Spec: networkingv1.IngressSpec{
			Rules: specRules,
			TLS:   specTLS,
		},
	}

	return ingressInterface.Create(context.TODO(), ingress, metav1.CreateOptions{})
}

func deleteAll(resourceIdentifier string, deploymentInterface cappsv1.DeploymentInterface, configMapInterface capiv1.ConfigMapInterface, serviceInterface capiv1.ServiceInterface, ingressInterface cnetworkingv1.IngressInterface) {
	deletePolicy := metav1.DeletePropagationForeground
	_ = deploymentInterface.Delete(context.TODO(), resourceIdentifier, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	_ = configMapInterface.Delete(context.TODO(), resourceIdentifier, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	_ = serviceInterface.Delete(context.TODO(), resourceIdentifier, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	_ = ingressInterface.Delete(context.TODO(), resourceIdentifier, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
}
