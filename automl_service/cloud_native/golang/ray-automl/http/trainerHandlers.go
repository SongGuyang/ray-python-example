package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	automlv1 "github.com/ray-automl/apis/automl/v1"
	"github.com/ray-automl/controllers/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// trainerCreateV1
// @Tags trainerCreateV1
// @ID trainerCreateV1
// @Summary trainerCreateV1 from ray-automl-operator
// @Param data body TrainerStartReq{} true "request for trainer create V1"
// @Success 200 {object} Response{success=bool,code=int,data=string}
// @Failure 500 {object} Response{success=bool,code=int,data=string}
// @Router /api/v1/trainer/create [post]
func (r *RestServer) trainerCreateV1(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		message := fmt.Sprintf("the request for pod recreate is failed: %v", err)
		serverLog.Error(err, message)
		JSONFailedResponse(c, err, message)
		return
	}

	trainerStartReq := &TrainerStartReq{}
	err = json.Unmarshal(data, trainerStartReq)
	if err != nil {
		JSONFailedResponse(c, err, fmt.Sprintf("invalid json :%v", string(data)))
		return
	}

	serverLog.Info("trainerCreateV1 received", "trainerStartReq", trainerStartReq)

	trainers := r.automlV1Interface.Trainers(trainerStartReq.Namespace)
	trainerInstance := trainerGenerator(trainerStartReq)
	serverLog.Info("trainerCreateV1 trainer", "instance", trainerInstance)
	trainer, err := trainers.Create(context.TODO(), trainerInstance, metav1.CreateOptions{})
	if err != nil {
		serverLog.Error(err, "trainerCreateV1 failed")
		JSONFailedResponse(c, err, fmt.Sprintf("invalid create trainer :%v", string(data)))
		return
	}

	serverLog.Info("success to submit trainer start request", "trainer", trainer)

	JSONSuccessResponse(c, nil, "success to submit trainer start request")
}

type TrainerStartReq struct {
	ProxyName   string                       `json:"proxyName,omitempty"`
	Name        string                       `json:"name,omitempty"`
	Namespace   string                       `json:"namespace,omitempty"`
	StartParams map[string]string            `json:"startParams,omitempty"`
	Image       string                       `json:"image,omitempty"`
	Workers     map[string]map[string]string `json:"workers,omitempty"`
}

func trainerGenerator(trainerStartReq *TrainerStartReq) *automlv1.Trainer {
	trainer := &automlv1.Trainer{}
	trainer.APIVersion = automlv1.GroupVersion.String()
	trainer.Kind = "Trainer"
	trainer.Name = trainerStartReq.Name
	if trainer.Labels == nil {
		trainer.Labels = map[string]string{}
	}
	trainer.Labels[utils.ProxyLabelSelector] = trainerStartReq.ProxyName
	trainer.Labels[utils.TrainerLabelSelector] = trainerStartReq.Name
	trainer.Namespace = trainerStartReq.Namespace
	trainer.Spec.StartParams = trainerStartReq.StartParams
	trainer.Spec.Image = trainerStartReq.Image
	trainer.Spec.Workers = trainerStartReq.Workers

	return trainer
}
