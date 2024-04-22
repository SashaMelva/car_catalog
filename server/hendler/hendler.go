package hendler

import (
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

type ResponseBody struct {
	Message      string
	MessageError string
}

func NewService(log *zap.SugaredLogger, application *app.App, timeout time.Duration) *Service {
	return &Service{
		Logger: *log,
		app:    *application,
	}
}
