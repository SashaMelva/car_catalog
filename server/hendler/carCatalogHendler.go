package hendler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	model "github.com/SashaMelva/car_catalog/internal/storage/models"
)

func (s *Service) CarCatalogHendler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == http.MethodGet {
	}
	if req.Method == http.MethodPut {
		s.updateCars(w, req, ctx)
		return
	}
	if req.Method == http.MethodPost {
		s.addCarByRegNums(w, req, ctx)
		return
	}
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
		return
	}
}

func (s *Service) addCarByRegNums(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	regNums := model.RegNumsCatalog{}
	body, err := io.ReadAll(req.Body)

	if err != nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusInternalServerError,
			Message: []byte(err.Error()),
		}, w)
		return
	} else {
		err = json.Unmarshal(body, &regNums)
		if err != nil {
			returnError(&ErrorResponseBody{
				Status:  http.StatusInternalServerError,
				Message: []byte(err.Error()),
			}, w)
			return
		}
	}

	if regNums.RegNums == nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusInternalServerError,
			Message: []byte("Регистрационные номера машин не найдены"),
		}, w)
		return
	}

	err = s.app.AddCarByRegNums(regNums.RegNums)

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

func (s *Service) updateCars(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	cars := model.CarCatalog{}
	body, err := io.ReadAll(req.Body)

	if err != nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusInternalServerError,
			Message: []byte(err.Error()),
		}, w)
		return
	} else {
		err = json.Unmarshal(body, &cars)
		if err != nil {
			returnError(&ErrorResponseBody{
				Status:  http.StatusInternalServerError,
				Message: []byte(err.Error()),
			}, w)
			return
		}
	}

	if cars.Cars == nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusInternalServerError,
			Message: []byte("Данные пусте"),
		}, w)
		return
	}

	err = s.app.UpdateCars(cars.Cars)

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
