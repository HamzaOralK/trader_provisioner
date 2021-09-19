package utility

import (
	"github.com/Coinoner/trader_provisioner/pkg/config"
	v1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	apiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	networkingv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
)

func CreateClientSets() (appsv1.DeploymentInterface, apiv1.ConfigMapInterface, apiv1.ServiceInterface, networkingv1.IngressInterface) {
	deploymentsInterface := config.ApplicationConfig.KubernetesClientSet.AppsV1().Deployments(v1.NamespaceDefault)
	configMapInterface := config.ApplicationConfig.KubernetesClientSet.CoreV1().ConfigMaps(v1.NamespaceDefault)
	serviceInterface := config.ApplicationConfig.KubernetesClientSet.CoreV1().Services(v1.NamespaceDefault)
	networkingInterface := config.ApplicationConfig.KubernetesClientSet.NetworkingV1().Ingresses(v1.NamespaceDefault)
	return deploymentsInterface, configMapInterface, serviceInterface, networkingInterface
}
