package http

import (
	"context"
	"errors"
	automlv1 "github.com/ray-automl/pkg/client/clientset/versioned/typed/automl/v1"
	"k8s.io/client-go/rest"
	"net/http"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	_ "github.com/ray-automl/docs"
)

var serverLog = logf.Log.WithName("server")

type RestServer struct {
	server            *http.Server
	automlV1Interface automlv1.AutomlV1Interface
}

// Start will start the rest server
func (r *RestServer) Start(errorChan chan<- error) {
	go func() {
		err := r.server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			errorChan <- err
		}
	}()
}

// Stop will shut down the rest server
func (r *RestServer) Stop(ctx context.Context) error {
	serverLog.Info("server shutting down")
	return r.server.Shutdown(ctx)
}

// New will create a rest api server: debug link: http://localhost:7070/swagger/index.html
func New(config *rest.Config) (*RestServer, error) {

	automlV1Client, err := automlv1.NewForConfig(config)
	if err != nil {
		serverLog.Error(err, "failed to create automlV1Client")
	}

	restServer := &RestServer{
		server:            nil,
		automlV1Interface: automlV1Client,
	}
	server := &http.Server{
		Addr:         ":" + DefaultRayOperatorServerPort,
		Handler:      restServer.setupRoute(DefaultRayOperatorServerStaticPath),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	server.SetKeepAlivesEnabled(true)
	restServer.server = server
	return restServer, nil
}
