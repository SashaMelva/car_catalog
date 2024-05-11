package hendler

import (
	"net/http"
	"strconv"
	"strings"

	model "github.com/SashaMelva/car_catalog/internal/storage/models"
	"github.com/SashaMelva/car_catalog/server/filter"
	"github.com/gin-gonic/gin"
)

func (s *Service) GetCarsCatalogHendler(ctx *gin.Context) {
	limitStr := ctx.Params.ByName("limit")
	offsetStr := ctx.Params.ByName("offset")
	regNums := ctx.Params.ByName("regNums")
	mark := ctx.Params.ByName("mark")
	model := ctx.Params.ByName("model")
	year := ctx.Params.ByName("year")
	periodYear := ctx.Params.ByName("periodYear")

	option := filter.NewOption()

	if limitStr == "" {
		option.Limit = 100
	} else {
		limit, err := strconv.Atoi(limitStr)

		if err != nil {
			s.log.Error(err.Error(), http.StatusBadRequest)
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		option.Limit = limit
	}

	if offsetStr == "" {
		option.Offset = 0
	} else {
		offset, err := strconv.Atoi(offsetStr)

		if err != nil {
			s.log.Error(err.Error(), http.StatusBadRequest)
			ctx.JSON(http.StatusBadRequest, err.Error())
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

		s.log.Info(lenI)
		if len(splitMark) > 1 {
			for i := range splitMark {
				s.log.Debug(splitMark)
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
		s.log.Info(len(splitYear), lenI, splitYear)
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

		s.log.Info(split[0], split[1])
		if split[0] == "" && split[1] == "" {
			s.log.Error("Годы периода не могут быть пустыми", http.StatusBadRequest)
			ctx.JSON(http.StatusBadRequest, "Годы периода не могут быть пустыми")
			return
		} else if split[0] == "" {
			option.AddFileds(filter.ParamYear, filter.OperatorLowerThen, split[1], filter.DateInt, filter.GroupNil)
		} else if split[1] == "" {
			option.AddFileds(filter.ParamYear, filter.OperatorHigherThen, split[0], filter.DateInt, filter.GroupNil)
		} else {
			option.AddFileds(filter.ParamYear, filter.OperatorBetween, split[0]+" and "+split[1], filter.DateInt, filter.GroupNil)
		}
	}

	cars, err := s.app.GetCars(option)

	if err != nil {
		s.log.Error(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	s.log.Info("OK")
	ctx.JSON(http.StatusOK, cars)
}

func (s *Service) AddCarsCatalogHendler(ctx *gin.Context) {
	cars := model.CarCatalog{}

	if err := ctx.ShouldBindJSON(&cars); err != nil {
		s.log.Error(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	s.log.Debug(cars)

	if cars.Cars == nil {
		s.log.Error("Данные пусте", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, "Данные пусте")
		return
	}

	err := s.app.AddCarByRegNums(cars.Cars)

	if err != nil {
		s.log.Error(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	s.log.Info("OK")
	ctx.JSON(http.StatusOK, gin.H{})
}

func (s *Service) UpdateCarsCatalogHendler(ctx *gin.Context) {
	s.log.Info("Method Put, run update cars")

	cars := model.CarCatalog{}

	if err := ctx.ShouldBindJSON(&cars); err != nil {
		s.log.Error(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	s.log.Debug(cars)

	if cars.Cars == nil {
		s.log.Error("Данные пусте", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, "Данные пусте")
		return
	}

	err := s.app.UpdateCars(cars.Cars)

	if err != nil {
		s.log.Error(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	s.log.Info("OK")
	ctx.JSON(http.StatusOK, gin.H{})
}

func (s *Service) DeleteCarByRegNumsHendler(ctx *gin.Context) {
	regNums := strings.Split(ctx.Params.ByName("regNums"), ",")
	if len(regNums) == 0 {
		s.log.Error("Не указан регистрационные номер", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, "Для удаления машины необходим регистрационный номер")
		return
	}

	s.log.Debug("Reg nums: ", regNums)
	err := s.app.DeleteCarByRegNum(regNums)

	if err != nil {
		s.log.Error(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	s.log.Info("OK")
	ctx.JSON(http.StatusOK, gin.H{})
}
