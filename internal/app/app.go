package app

import (
	"errors"
	"regexp"

	"github.com/SashaMelva/car_catalog/internal/storage/memory"
	"go.uber.org/zap"
)

type App struct {
	storage     *memory.Storage
	Logger      *zap.SugaredLogger
	RegexRegNum *regexp.Regexp
}

func New(logger *zap.SugaredLogger, storage *memory.Storage) *App {
	r, _ := regexp.Compile(`^[АВЕКМНОРСТУХ][0-9][0-9][0-9][АВЕКМНОРСТУХ][АВЕКМНОРСТУХ]$`)
	return &App{
		storage:     storage,
		Logger:      logger,
		RegexRegNum: r,
	}
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
