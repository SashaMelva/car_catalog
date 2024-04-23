package hendler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	model "github.com/SashaMelva/car_catalog/internal/storage/models"
	"github.com/SashaMelva/car_catalog/server/filter"
)

func (s *Service) CarCatalogHendler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if req.Method == http.MethodGet {
		args := req.URL.Query()
		regNums := args.Get("regNum")
		mark := args.Get("mark")
		model := args.Get("model")
		year := args.Get("year")
		periodYear := args.Get("periodYear")

		option := filter.NewOption()
		if regNums != "" {
			option.AddFileds(filter.ParamRegNum, filter.OperatorEq, regNums, filter.DateStr)
		}
		if mark != "" {
			option.AddFileds(filter.ParamMark, filter.OperatorEq, mark, filter.DateStr)
		}
		if model != "" {
			option.AddFileds(filter.ParamModel, filter.OperatorEq, model, filter.DateStr)
		}
		if year != "" {
			option.AddFileds(filter.ParamYear, filter.OperatorEq, year, filter.DateInt)
		}
		if periodYear != "" {
			split := strings.Split(periodYear, ":")

			if split[0] == "" && split[1] == "" {
				returnError(&ErrorResponseBody{
					Status:  http.StatusBadRequest,
					Message: []byte("Годы периода не могут быть пустыми"),
				}, w)
				return
			} else if split[0] == "" {
				option.AddFileds(filter.ParamYear, filter.OperatorLowerThen, split[1], filter.DateInt)
			} else if split[1] == "" {
				option.AddFileds(filter.ParamYear, filter.OperatorHigherThen, split[0], filter.DateInt)
			} else {
				option.AddFileds(filter.ParamYear, filter.OperatorBetween, split[0]+" and "+split[1], filter.DateInt)
			}
		}

		s.getAllCars(option, w, req, ctx)
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
			s.deleteCarByRegNums(regNum, w, req, ctx)
			return
		}
		returnError(&ErrorResponseBody{
			Status:  http.StatusBadRequest,
			Message: []byte("Для удаления машины необходим регистрационный номер"),
		}, w)
		return
	}
}

func (s *Service) getAllCars(option filter.Option, w http.ResponseWriter, req *http.Request, ctx context.Context) {
	cars, err := s.app.GetCars(option)

	if err != nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusInternalServerError,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	json, err := json.Marshal(cars)

	if err != nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusInternalServerError,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
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

func (s *Service) deleteCarByRegNums(regNums []string, w http.ResponseWriter, req *http.Request, ctx context.Context) {
	err := s.app.DeleteCarByRegNum(regNums)

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
