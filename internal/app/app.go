package app

import (
	"errors"
	"regexp"

	"github.com/SashaMelva/car_catalog/internal/storage/memory"
	model "github.com/SashaMelva/car_catalog/internal/storage/models"
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
	return nil
}

func (a *App) DeleteCarByRegNum(regNum string) error {
	err := a.validRegNum(regNum)

	if err != nil {
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

func (a *App) DeleteCarByRegNums(regNums []string) error {
	var err error
	for i := range regNums {
		err = a.validRegNum(regNums[i])
		if err != nil {
			a.Logger.Error(err)
			return err
		}
	}

	err = a.storage.DeleteCarByRegNum(regNums[0])

	if err != nil {
		a.Logger.Error(err)
	}

	return nil
}

func (a *App) validRegNum(regNum string) error {
	matched := a.RegexRegNum.MatchString(regNum)

	if !matched {
		return errors.New("Регистрационный номер машины не соответсвует стандарту ГОСТ РФ")
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
