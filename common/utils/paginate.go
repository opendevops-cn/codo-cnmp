package utils

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// K8sPaginate k8s资源对象分页
func K8sPaginate[T any](items []T, page, pageSize uint32) ([]T, uint32) {
	if len(items) == 0 {
		return nil, 0
	}
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > uint32(len(items)) {
		end = uint32(len(items))
	}
	return items[start:end], uint32(len(items))
}

type WorkloadType string

const (
	DeploymentType  WorkloadType = "Deployment"
	StatefulSetType WorkloadType = "StatefulSet"
	DaemonSetType   WorkloadType = "DaemonSet"
	PodType         WorkloadType = "Pod"
	ServiceType     WorkloadType = "Service"
	ConfigMapType   WorkloadType = "ConfigMap"
	SecretType      WorkloadType = "Secret"
)

func getWorkloadType(obj runtime.Object) WorkloadType {
	switch obj.(type) {
	case *appsv1.DeploymentList:
		return DeploymentType
	case *appsv1.StatefulSetList:
		return StatefulSetType
	case *appsv1.DaemonSetList:
		return DaemonSetType
	case *corev1.PodList:
		return PodType
	case *corev1.ServiceList:
		return ServiceType
	case *corev1.ConfigMapList:
		return ConfigMapType
	case *corev1.SecretList:
		return SecretType
	default:
		return ""
	}
}
