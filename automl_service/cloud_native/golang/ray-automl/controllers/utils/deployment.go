package utils

import (
	"bytes"
	"fmt"
	"github.com/go-logr/logr"
	automlv1 "github.com/ray-automl/apis/automl/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sort"
	"strings"
)

func NewDeployment(object interface{}, log logr.Logger) *appsv1.Deployment {
	deployment := &appsv1.Deployment{}
	switch object.(type) {
	case *automlv1.Trainer:
		deployment = NewDeploymentInstanceTrainer(object.(*automlv1.Trainer))
		log.Info("create deployment for NewDeploymentInstanceTrainer", "deployment", deployment)
	case *automlv1.Proxy:
		deployment = NewDeploymentInstanceProxy(object.(*automlv1.Proxy))
		log.Info("create deployment for NewDeploymentInstanceProxy", "deployment", deployment)
	}
	return deployment
}

func GetMainContainer(podSpec corev1.PodTemplateSpec, containerName string) (int, *corev1.Container) {
	for index, container := range podSpec.Spec.Containers {
		if strings.EqualFold(container.Name, containerName) {
			return index, container.DeepCopy()
		}
	}
	return -1, nil
}

func convertParamMap(startParams map[string]string) (s string) {
	keys := make([]string, 0)
	for k := range startParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	flags := new(bytes.Buffer)
	for _, k := range keys {
		_, _ = fmt.Fprintf(flags, " --%s=%s ", k, startParams[k])
	}
	return flags.String()
}

func getOrDefault(params map[string]string, key string, defaultValue string) string {
	if params[key] != "" {
		return params[key]
	}
	return defaultValue
}

func BuildResourceList(cpu string, memory string, ephemeralStorage string) corev1.ResourceList {
	return corev1.ResourceList{
		corev1.ResourceCPU:              resource.MustParse(cpu),
		corev1.ResourceMemory:           resource.MustParse(memory),
		corev1.ResourceEphemeralStorage: resource.MustParse(ephemeralStorage),
	}
}
