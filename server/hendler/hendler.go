package hendler

import (
	"net/http"
	"sync"
	"time"

	"github.com/SashaMelva/car_catalog/internal/app"
	"go.uber.org/zap"
)

type Service struct {
	Logger zap.SugaredLogger
	app    app.App
	sync.RWMutex
}

type ErrorResponseBody struct {
	Status  int
	Message []byte
}

func NewService(log *zap.SugaredLogger, application *app.App, timeout time.Duration) *Service {
	return &Service{
		Logger: *log,
		app:    *application,
	}
}

func returnError(errorResponse *ErrorResponseBody, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errorResponse.Status)
	w.Write(errorResponse.Message)
}
