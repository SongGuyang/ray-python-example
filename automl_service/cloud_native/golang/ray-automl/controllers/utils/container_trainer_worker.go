package utils

import (
	"fmt"
	"github.com/go-logr/logr"
	automlv1 "github.com/ray-automl/apis/automl/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	WorkerId = "worker-id"
)

func NewDeploymentInstanceWorker(instance *automlv1.Trainer, log logr.Logger) []*appsv1.Deployment {
	workerDeployments := make([]*appsv1.Deployment, 0)
	if instance.Spec.Workers != nil && len(instance.Spec.Workers) > 0 {
		trainerPort := GetServicePort(instance.Spec.StartParams, automlv1.TrainerContainerPortNumber, automlv1.TrainerContainerPortNumberDefault)

		trainerSvcAddr := fmt.Sprintf("%s.%s.svc.%s:%v",
			instance.Name,
			instance.Namespace,
			GetClusterDomainName(),
			trainerPort)

		for group, params := range instance.Spec.Workers {

			log.Info("group and params", "group", group, "params", params)

			workerDeployment := &appsv1.Deployment{}
			workerDeployment.Name = instance.Name + "-" + group
			workerDeployment.Namespace = instance.Namespace

			container := &corev1.Container{Name: automlv1.TrainerWorkerContainerName}
			workerParams := map[string]string{}
			workerParams[automlv1.TrainerAddress] = trainerSvcAddr
			workerParams[WorkerId] = "$MY_POD_NAME"

			paramMap := convertParamMap(workerParams)

			container.Env = append(container.Env, corev1.EnvVar{
				Name: "MY_POD_NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			})

			container.Command = []string{"python", "-m"}
			container.Args = []string{"automl.worker " + paramMap}
			container.Image = instance.Spec.Image

			cpu := getOrDefault(params, Cpu, "1")
			memory := getOrDefault(params, Memory, "1Gi")
			disk := getOrDefault(params, Disk, "1Gi")
			gpuCard := getOrDefault(params, GpuCard, GpuCardDefault)
			gpu := getOrDefault(params, Gpu, "0")

			container.Resources.Limits = BuildResourceList(cpu, memory, disk)
			container.Resources.Requests = BuildResourceList(cpu, memory, disk)

			if gpu != "" && gpu != "0" {
				gpuResourceName := corev1.ResourceName(gpuCard)
				quantity := resource.MustParse(gpu)
				container.Resources.Limits[gpuResourceName] = quantity
			}

			workerDeployment.Spec.Template.Spec.Containers = append(workerDeployment.Spec.Template.Spec.Containers, *container)

			workerDeployment.Labels = instance.Labels
			workerDeployment.Labels[TrainerWorkerLabelSelector] = instance.Name + "-" + group
			workerDeployment.Annotations = instance.Annotations
			if workerDeployment.Spec.Selector == nil {
				workerDeployment.Spec.Selector = &metav1.LabelSelector{
					MatchLabels: workerDeployment.Labels,
				}
			}
			workerDeployment.Spec.Template.Labels = workerDeployment.Labels
			workerDeployment.Spec.Template.Annotations = workerDeployment.Annotations
			workerDeployments = append(workerDeployments, workerDeployment.DeepCopy())
		}
	}
	return workerDeployments
}
