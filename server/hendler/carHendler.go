package hendler

import (
	"context"
	"net/http"
	"strings"
	"time"
)

func (s *Service) CarHendler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodDelete {
		args := req.URL.Query()
		regNumsStr := args.Get("regNums")
		if regNumsStr != "" {
			regNum := strings.Split(regNumsStr, ",")

			if len(regNum) == 1 {
				s.deleteCarByRegNum(regNum[0], w, req, ctx)
				return
			}
			if len(regNum) > 1 {
				s.deleteCarByRegNums(regNum, w, req, ctx)
				return
			}
		}
		returnError(&ErrorResponseBody{
			Status:  http.StatusBadRequest,
			Message: []byte("Для удаления машины необходим регистрационный номер"),
		}, w)
	}
}

func (s *Service) deleteCarByRegNum(regNum string, w http.ResponseWriter, req *http.Request, ctx context.Context) {
	err := s.app.DeleteCarByRegNum(regNum)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusInternalServerError)
		returnError(&ErrorResponseBody{
			Status:  http.StatusInternalServerError,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Service) deleteCarByRegNums(regNums []string, w http.ResponseWriter, req *http.Request, ctx context.Context) {
	err := s.app.DeleteCarByRegNums(regNums)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusInternalServerError)
		returnError(&ErrorResponseBody{
			Status:  http.StatusInternalServerError,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
