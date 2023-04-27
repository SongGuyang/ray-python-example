package utils

import (
	"fmt"
	automlv1 "github.com/ray-automl/apis/automl/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"strconv"
)

func NewDeploymentInstanceTrainer(instance *automlv1.Trainer) *appsv1.Deployment {
	BuildTrainerContainer(instance)
	deployment := &appsv1.Deployment{}
	deployment.Name = instance.Name
	deployment.Namespace = instance.Namespace
	deployment.Spec = instance.Spec.DeploySpec
	return deployment
}

func BuildTrainerContainer(instance *automlv1.Trainer) {
	if instance.Spec.StartParams != nil {

		proxyPort := GetServicePort(instance.Spec.StartParams, automlv1.ProxyContainerPortNumber, automlv1.ProxyContainerPortNumberDefault)

		operatorAddress := fmt.Sprintf("%s.%s.svc.%s:%v",
			OperatorName,
			instance.Namespace,
			GetClusterDomainName(),
			OperatorPort)

		proxySvcAddr := fmt.Sprintf("%s.%s.svc.%s:%v",
			instance.Labels[ProxyLabelSelector],
			instance.Namespace,
			GetClusterDomainName(),
			proxyPort)

		trainerSvcAddr := fmt.Sprintf("%s.%s.svc.%s",
			instance.Name,
			instance.Namespace,
			GetClusterDomainName())

		instance.Spec.StartParams[OperatorAddress] = operatorAddress
		instance.Spec.StartParams[automlv1.ProxyAddress] = proxySvcAddr
		instance.Spec.StartParams[automlv1.TrainerId] = trainerSvcAddr
		instance.Spec.StartParams[automlv1.TaskId] = "$MY_POD_NAME"
		if instance.Spec.StartParams[automlv1.TrainerGrpcPortNumber] == "" {
			instance.Spec.StartParams[automlv1.TrainerGrpcPortNumber] = strconv.FormatInt(int64(automlv1.TrainerContainerPortNumberDefault), 10)
		}
	}
	index, container := GetMainContainer(instance.Spec.DeploySpec.Template, automlv1.TrainerContainerName)

	if index < 0 {
		instance.Spec.DeploySpec.Template.Spec.Containers = append(instance.Spec.DeploySpec.Template.Spec.Containers, corev1.Container{Name: automlv1.TrainerContainerName})
	} else {
		paramMap := convertParamMap(instance.Spec.StartParams)
		container.Command = []string{"python", "-m"}
		container.Args = []string{"automl.trainer " + paramMap}
		container.Image = instance.Spec.Image

		container.Env = append(container.Env, corev1.EnvVar{
			Name: "MY_POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		})
		cpu := getOrDefault(instance.Spec.StartParams, Cpu, "1")
		memory := getOrDefault(instance.Spec.StartParams, Memory, "1Gi")
		disk := getOrDefault(instance.Spec.StartParams, Disk, "1Gi")
		container.Resources.Limits = BuildResourceList(cpu, memory, disk)
		container.Resources.Requests = BuildResourceList(cpu, memory, disk)
		instance.Spec.DeploySpec.Template.Spec.Containers[index] = *container
	}

	instance.Spec.DeploySpec.Template.Labels = instance.Labels
	instance.Spec.DeploySpec.Template.Annotations = instance.Annotations
	instance.Spec.DeploySpec.Selector.MatchLabels = instance.Labels
}
