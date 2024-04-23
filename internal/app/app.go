package app

import (
	"errors"
	"regexp"

	"github.com/SashaMelva/car_catalog/internal/storage/memory"
	model "github.com/SashaMelva/car_catalog/internal/storage/models"
	"github.com/SashaMelva/car_catalog/server/filter"
	"go.uber.org/zap"
)

type App struct {
	storage     *memory.Storage
	Logger      *zap.SugaredLogger
	RegexRegNum *regexp.Regexp
}

func New(logger *zap.SugaredLogger, storage *memory.Storage) *App {
	r, _ := regexp.Compile(`^[АВЕКМНОРСТУХ][0-9][0-9][0-9][АВЕКМНОРСТУХ][АВЕКМНОРСТУХ][0-9][0-9][0-9]$`)
	return &App{
		storage:     storage,
		Logger:      logger,
		RegexRegNum: r,
	}
}

func (a *App) GetCars(option filter.Option) (*model.CarCatalog, error) {
	var err error
	catalog := &model.CarCatalog{}

	if len(option.Fileds) == 0 {
		catalog, err = a.storage.GetAllCars()
	} else {
		query := ""

		for i := range option.Fileds {
			query += option.Fileds[i].Param + " " + option.Fileds[i].Operator + " " + option.Fileds[i].Value
			if i != len(option.Fileds)-1 {
				query += " and "
			}
		}

		catalog, err = a.storage.GetCarsByFilter(query)
	}

	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	return catalog, nil
}

func (a *App) AddCarByRegNums(regNums []string) error {
	var err error

	for i := range regNums {
		err = a.validRegNum(regNums[i])
		if err != nil {
			a.Logger.Error(err)
			return err
		}
	}

	return nil
}

func (a *App) UpdateCars(cars []*model.Car) error {
	if len(cars) == 1 {
		err := a.storage.UpdateCarFromCatalog(cars[0])

		if err != nil {
			a.Logger.Error(err)
			return err
		}
	} else if len(cars) > 1 {
		err := a.storage.UpdateCarsFromCatalog(cars)

		if err != nil {
			a.Logger.Error(err)
			return err
		}
	}

	return errors.New("Данные для обновления не могут быть пустыми")
}

func (a *App) DeleteCarByRegNum(regNum []string) error {
	var err error

	if len(regNum) > 0 {
		for i := range regNum {
			errNew := a.validRegNum(regNum[i])

			if errNew != nil {
				err = errors.Join(errNew, err)
			}
		}

		if err != nil {
			a.Logger.Error(err)
			return err
		}
	} else {
		err = errors.New("Регистрационные номера машин пусты")
		a.Logger.Error(err)
		return err
	}

	err = a.storage.DeleteCarByRegNum(regNum)

	if err != nil {
		a.Logger.Error(err)
		return err
	}

	return nil
}

func (a *App) validRegNum(regNum string) error {
	matched := a.RegexRegNum.MatchString(regNum)

	if !matched {
		return errors.New("Регистрационный номер машины " + regNum + " не соответсвует стандарту ГОСТ РФ")
	}

	return nil
}

func (a *App) valiCarInfo(carInfo model.Car) error {

	matched := a.RegexRegNum.MatchString(carInfo.RegNum)

	if !matched {
		return errors.New("Регистрационный номер машины не соответсвует стандарту ГОСТ РФ")
	}

	if carInfo.Year > 1886 && carInfo.Year < 3000 {
		return errors.New("Год машины не относится к реальности")
	}

	return nil
}
