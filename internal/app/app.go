package app

import (
	"errors"
	"regexp"

	"github.com/SashaMelva/car_catalog/internal/storage/memory"
	model "github.com/SashaMelva/car_catalog/internal/storage/models"
	"github.com/SashaMelva/car_catalog/server/client"
	"github.com/SashaMelva/car_catalog/server/filter"
	"go.uber.org/zap"
)

type App struct {
	HostClientApi string
	storage       *memory.Storage
	Logger        *zap.SugaredLogger
	RegexRegNum   *regexp.Regexp
}

func New(logger *zap.SugaredLogger, storage *memory.Storage, host *string) *App {
	r, _ := regexp.Compile(`^[ABEKMHOPCTYXАВЕКМНОРСТУХ][0-9][0-9][0-9][ABEKMHOPCTYXАВЕКМНОРСТУХ][ABEKMHOPCTYXАВЕКМНОРСТУХ][0-9][0-9][0-9]$`)
	return &App{
		HostClientApi: *host,
		storage:       storage,
		Logger:        logger,
		RegexRegNum:   r,
	}
}

func (a *App) GetCars(option filter.Option) (*model.CarCatalog, error) {
	var err error
	catalog := &model.CarCatalog{}
	a.Logger.Debug("Get cars with plugin: Limit" + string(option.Limit) + " Offset " + string(option.Offset))

	if len(option.Fileds) == 0 {
		a.Logger.Info("Get cars with not filter. Limit" + string(option.Limit) + " Offset " + string(option.Offset))
		catalog, err = a.storage.GetAllCars(option.Limit, option.Offset)
	} else {
		a.Logger.Info("Get cars with filter. Limit" + string(option.Limit) + " Offset " + string(option.Offset))
		query := ""

		for i := range option.Fileds {
			a.Logger.Debug(option.Fileds[i])
			if i != 0 {
				if option.Fileds[i].Group != filter.GroupNil && option.Fileds[i].Group != filter.GroupStart {
					query += " or "
				} else {
					query += " and "
				}
			}

			if option.Fileds[i].Group == filter.GroupStart {
				query += " ( "
			}

			if option.Fileds[i].DataType == filter.DateStr {
				query += option.Fileds[i].Param + option.Fileds[i].Operator + "'" + option.Fileds[i].Value + "'"
			} else {
				query += option.Fileds[i].Param + option.Fileds[i].Operator + option.Fileds[i].Value
			}

			if option.Fileds[i].Group == filter.GroupEnd {
				query += " ) "
			}
		}
		a.Logger.Debug("Parce filter", query)
		catalog, err = a.storage.GetCarsByFilter(query, option.Limit, option.Offset)
	}

	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	return catalog, nil
}

func (a *App) AddCarByRegNums(regNums []string) error {
	var err error
	a.Logger.Info("Run validate reg nums cars")
	for i := range regNums {
		err = a.validRegNum(regNums[i])
		if err != nil {
			a.Logger.Error("Valid error: ", err)
			return err
		}
	}

	for i := range regNums {
		req, err := client.GetInfoCarByRegNum(regNums[i], a.HostClientApi, a.Logger)
		if err != nil {
			a.Logger.Error("Error: ", err)
			return err
		}

		if req.Code == "200" {
			err = a.storage.AddCarCatalog(*req.Content)

		} else {
			return errors.New(req.Description)
		}
	}

	return nil
}

func (a *App) UpdateCars(cars []*model.Car) error {
	if len(cars) == 1 {
		a.Logger.Info("Update one car info ", cars[0].RegNum)
		err := a.storage.UpdateCarFromCatalog(cars[0])

		if err != nil {
			a.Logger.Error(err)
			return err
		}
	} else if len(cars) > 1 {
		a.Logger.Info("Update many cars info ")
		err := a.storage.UpdateCarsFromCatalog(cars)

		if err != nil {
			a.Logger.Error(err)
			return err
		}
	} else {
		return errors.New("Данные для обновления не могут быть пустыми")
	}

	return nil
}

func (a *App) DeleteCarByRegNum(regNum []string) error {
	var err error
	a.Logger.Info("Run validate reg nums cars")
	if len(regNum) > 0 {
		for i := range regNum {
			errNew := a.validRegNum(regNum[i])

			if errNew != nil {
				err = errors.Join(errNew, err)
			}
		}

		if err != nil {
			a.Logger.Error("Valid error: ", err)
			return err
		}
	} else {
		err = errors.New("Регистрационные номера машин пусты")
		a.Logger.Error(err)
		return err
	}

	a.Logger.Info("Delete cars run method sql")
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
