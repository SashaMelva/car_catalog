package hendler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	model "github.com/SashaMelva/car_catalog/internal/storage/models"
	"github.com/SashaMelva/car_catalog/server/filter"
)

func (s *Service) CarCatalogHendler(w http.ResponseWriter, req *http.Request) {
	s.Logger.Debug("Path responce ", req.Method, req.URL)
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

		if args.Get("limit") == "" {
			option.Limit = 100
		} else {
			limit, err := strconv.Atoi(args.Get("limit"))

			if err != nil {
				s.Logger.Error(err)
				returnError(&ErrorResponseBody{
					Status:  http.StatusBadRequest,
					Message: []byte(err.Error()),
				}, w)
				return
			}
			option.Limit = limit
		}

		if args.Get("offset") == "" {
			option.Offset = 0
		} else {
			offset, err := strconv.Atoi(args.Get("limit"))

			if err != nil {
				s.Logger.Error(err)
				returnError(&ErrorResponseBody{
					Status:  http.StatusBadRequest,
					Message: []byte(err.Error()),
				}, w)
				return
			}
			option.Offset = offset
		}

		if regNums != "" {
			splitRegNums := strings.Split(regNums, ",")
			lenI := len(splitRegNums) - 1

			if len(splitRegNums) > 1 {
				for i := range splitRegNums {
					if i == 0 {
						option.AddFileds(filter.ParamRegNum, filter.OperatorEq, splitRegNums[i], filter.DateStr, filter.GroupStart)
					} else if i == lenI {
						option.AddFileds(filter.ParamRegNum, filter.OperatorEq, splitRegNums[i], filter.DateStr, filter.GroupEnd)
					} else {
						option.AddFileds(filter.ParamRegNum, filter.OperatorEq, splitRegNums[i], filter.DateStr, filter.GroupElement)
					}
				}
			} else {
				option.AddFileds(filter.ParamRegNum, filter.OperatorEq, splitRegNums[0], filter.DateStr, filter.GroupNil)
			}
		}
		if mark != "" {
			splitMark := strings.Split(mark, ",")
			lenI := len(splitMark) - 1

			s.app.Logger.Info(lenI)
			if len(splitMark) > 1 {
				for i := range splitMark {
					s.app.Logger.Debug(splitMark)
					if i == 0 {
						option.AddFileds(filter.ParamMark, filter.OperatorEq, splitMark[i], filter.DateStr, filter.GroupStart)
					} else if i == lenI {
						option.AddFileds(filter.ParamMark, filter.OperatorEq, splitMark[i], filter.DateStr, filter.GroupEnd)
					} else {
						option.AddFileds(filter.ParamMark, filter.OperatorEq, splitMark[i], filter.DateStr, filter.GroupElement)
					}
				}
			} else {
				option.AddFileds(filter.ParamMark, filter.OperatorEq, splitMark[0], filter.DateStr, filter.GroupNil)
			}
		}
		if model != "" {
			splitModel := strings.Split(model, ",")
			lenI := len(splitModel) - 1

			if len(splitModel) > 1 {
				for i := range splitModel {
					if i == 0 {
						option.AddFileds(filter.ParamModel, filter.OperatorEq, splitModel[i], filter.DateStr, filter.GroupStart)
					} else if i == lenI {
						option.AddFileds(filter.ParamModel, filter.OperatorEq, splitModel[i], filter.DateStr, filter.GroupEnd)
					} else {
						option.AddFileds(filter.ParamModel, filter.OperatorEq, splitModel[i], filter.DateStr, filter.GroupElement)
					}
				}
			} else {
				option.AddFileds(filter.ParamModel, filter.OperatorEq, splitModel[0], filter.DateStr, filter.GroupNil)
			}
		}
		if year != "" {
			splitYear := strings.Split(year, ",")
			lenI := len(splitYear) - 1
			s.app.Logger.Info(len(splitYear), lenI, splitYear)
			if len(splitYear) > 1 {
				for i := range splitYear {
					if i == 0 {
						option.AddFileds(filter.ParamYear, filter.OperatorEq, splitYear[i], filter.DateStr, filter.GroupStart)
					} else if i == lenI {
						option.AddFileds(filter.ParamYear, filter.OperatorEq, splitYear[i], filter.DateStr, filter.GroupEnd)
					} else {
						option.AddFileds(filter.ParamYear, filter.OperatorEq, splitYear[i], filter.DateStr, filter.GroupElement)
					}
				}
			} else {
				option.AddFileds(filter.ParamYear, filter.OperatorEq, splitYear[0], filter.DateStr, filter.GroupNil)
			}
		}
		if periodYear != "" {
			split := strings.Split(periodYear, ":")

			s.Logger.Info(split[0], split[1])
			if split[0] == "" && split[1] == "" {
				s.Logger.Error("Годы периода не могут быть пустыми")
				returnError(&ErrorResponseBody{
					Status:  http.StatusBadRequest,
					Message: []byte("Годы периода не могут быть пустыми"),
				}, w)
				return
			} else if split[0] == "" {
				option.AddFileds(filter.ParamYear, filter.OperatorLowerThen, split[1], filter.DateInt, filter.GroupNil)
			} else if split[1] == "" {
				option.AddFileds(filter.ParamYear, filter.OperatorHigherThen, split[0], filter.DateInt, filter.GroupNil)
			} else {
				option.AddFileds(filter.ParamYear, filter.OperatorBetween, split[0]+" and "+split[1], filter.DateInt, filter.GroupNil)
			}
		}

		s.Logger.Debug("Create filter: ", option.Fileds)

		s.getAllCars(option, w, req, ctx)
	}

	if req.Method == http.MethodPut {
		s.Logger.Info("Method Put, run update cars")
		s.updateCars(w, req, ctx)
		return
	}

	if req.Method == http.MethodPost {
		s.Logger.Info("Method Post, run create cars")
		s.addCarByRegNums(w, req, ctx)
		return
	}

	if req.Method == http.MethodDelete {
		args := req.URL.Query()
		regNumsStr := args.Get("regNums")

		s.Logger.Info("Method Delete, run delete cars ", regNumsStr)
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

	s.Logger.Debug("Method implementation not found")
}

// ShowAccount godoc
// @Summary      Get all cars
// @Description  get cars by params
// @Tags         car-catalog
// @Accept       json
// @Produce      json
// @Param        regNum   path string
// @Param        mark   path string
// @Param        model   path string
// @Param        year   path integer
// @Param        periodYear   path string
// @Success      200  {object}  model.CarCatalog
// @Failure      400  {object}  ErrorResponseBody
// @Failure      500  {object}  ErrorResponseBody
// @Router       /car-catalog [get]

func (s *Service) getAllCars(option filter.Option, w http.ResponseWriter, req *http.Request, ctx context.Context) {
	cars, err := s.app.GetCars(option)

	if err != nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusBadRequest,
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

	s.Logger.Info("OK")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// AddCars godoc
// @Summary      Add cars
// @Description  Add cars by reg Num
// @Tags         car-catalog
// @Accept       json
// @Produce      json
// Param         regNums   body   model.RegNumsCatalog
// @Success      200
// @Failure      400  {object}  ErrorResponseBody
// @Failure      500  {object}  ErrorResponseBody
// @Router       /car-catalog [post]
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
			Status:  http.StatusBadRequest,
			Message: []byte("Регистрационные номера машин не найдены"),
		}, w)
		return
	}
	s.Logger.Info("Регистрационные номера машин для добавления: ", regNums.RegNums)

	err = s.app.AddCarByRegNums(regNums.RegNums)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusBadRequest)
		returnError(&ErrorResponseBody{
			Status:  http.StatusBadRequest,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	s.Logger.Info("OK")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// UpdateCars godoc
// @Summary      Update cars
// @Description  Update many cars by reg nums
// @Tags         car-catalog
// @Accept       json
// @Produce      json
// Param         cars   body      model.CarCatalog
// @Success      200
// @Failure      400  {object}  ErrorResponseBody
// @Failure      500  {object}  ErrorResponseBody
// @Router       /car-catalog [put]
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
			Status:  http.StatusBadRequest,
			Message: []byte("Данные пусте"),
		}, w)
		return
	}

	err = s.app.UpdateCars(cars.Cars)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusBadRequest)
		returnError(&ErrorResponseBody{
			Status:  http.StatusBadRequest,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	s.Logger.Info("OK")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// DeleteCars godoc
// @Summary      Delete cars
// @Description  delete cars by reg nums. One or many
// @Tags         car-catalog
// @Accept       json
// @Produce      json
// @Param        regNum path string true "The string reg nums" example(A777AA777, A777AA777)
// @Success      200
// @Failure      400  {object}  ErrorResponseBody
// @Failure      500  {object}  ErrorResponseBody
// @Router       /car-catalog [delete]
func (s *Service) deleteCarByRegNums(regNums []string, w http.ResponseWriter, req *http.Request, ctx context.Context) {
	err := s.app.DeleteCarByRegNum(regNums)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusBadRequest)
		returnError(&ErrorResponseBody{
			Status:  http.StatusBadRequest,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	s.Logger.Info("OK")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
