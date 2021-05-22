package main

import (
	v1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	apiv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func createClientSets() (appsv1.DeploymentInterface, apiv1.ConfigMapInterface, apiv1.ServiceInterface) {
	deploymentsClient := kubernetesClientSet.AppsV1().Deployments(v1.NamespaceDefault)
	configMapClient := kubernetesClientSet.CoreV1().ConfigMaps(v1.NamespaceDefault)
	serviceClient := kubernetesClientSet.CoreV1().Services(v1.NamespaceDefault)
	return deploymentsClient, configMapClient, serviceClient
}
