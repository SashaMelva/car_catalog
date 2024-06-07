package http

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/SashaMelva/car_catalog/internal/app"
	"github.com/SashaMelva/car_catalog/internal/config"
	"github.com/SashaMelva/car_catalog/server/hendler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	srv *http.Server
	log *zap.SugaredLogger
}

func NewServer(log *zap.SugaredLogger, app *app.App, config *config.ConfigHttpServer) *Server {
	log.Info("URL api " + config.Host + ":" + config.Port)
	log.Debug("URL api running " + config.Host + ":" + config.Port)

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	h := hendler.NewHendler(log, app)

	router.GET("/", func(ctx *gin.Context) {
		fmt.Println("Hellow world)")
		log.Debug("Test path working")
	})

	router.GET("/car-catalog/", h.GetCarsCatalogHendler)
	router.GET("/car-catalog/wb", h.HandleConnections)
	router.POST("/car-catalog/", h.AddCarsCatalogHendler)
	router.PUT("/car-catalog/", h.UpdateCarsCatalogHendler)
	router.DELETE("/car-catalog/", h.DeleteCarByRegNumsHendler)

	router.Run(":" + config.Port)
	return &Server{
		srv: &http.Server{
			Addr:    config.Host + ":" + config.Port,
			Handler: router,
		},
		log: log,
	}
}

func (s *Server) Start(ctx context.Context) {
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.Fatal("listen: ", err)
	}
}

func (s *Server) Stop(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Fatal("Server forced to shutdown: ", err)
	}

	os.Exit(1)
}
