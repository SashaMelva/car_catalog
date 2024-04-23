package memory

import (
	model "github.com/SashaMelva/car_catalog/internal/storage/models"
)

func (s *Storage) DeleteCarByRegNum(regNums []string) error {
	tx, err := s.ConnectionDB.Begin()

	if err != nil {
		return err
	}

	for i := range regNums {

		query := `delete car_catalog where reg_num = $1`
		_, err := s.ConnectionDB.Exec(query, regNums[i])

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) UpdateCarFromCatalog(car *model.Car) error {
	query := `update car_catalog set mark=$2 model=$3 year=$4 owner=$5 where reg_num=$1`
	_, err := s.ConnectionDB.Exec(query, car.RegNum, car.Mark, car.Model, car.Year, car.Owner)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateCarsFromCatalog(cars []*model.Car) error {
	tx, err := s.ConnectionDB.Begin()

	if err != nil {
		return err
	}

	for i := range cars {
		query := `update car_catalog set mark=$2 model=$3 year=$4 owner=$5 where reg_num=$1`
		_, err := tx.Exec(query, cars[i].RegNum, cars[i].Mark, cars[i].Model, cars[i].Year, cars[i].Owner)

		if err != nil {
			tx.Rollback()
			return err
		}

	}

	return tx.Commit()
}

func (s *Storage) AddCarCatalog(car model.Car) (int, error) {
	var carId int
	query := `insert into car_catalog(reg_num, mark, model, year, owner) values($1, $2, $3, $4, $5) RETURNING id`
	result := s.ConnectionDB.QueryRow(query, car.RegNum, car.Mark, car.Model, car.Year, car.Owner) // sql.Result
	err := result.Scan(&carId)

	if err != nil {
		return 0, err
	}

	return carId, nil
}

func (s *Storage) AddCarsCatalog(cars []*model.Car) error {
	tx, err := s.ConnectionDB.Begin()

	if err != nil {
		return err
	}

	for i := range cars {
		query := `insert into car_catalog(reg_num, mark, model, year, owner) values($1, $2, $3, $4, $5)`
		_, err := tx.Exec(query, cars[i].RegNum, cars[i].Mark, cars[i].Model, cars[i].Year, cars[i].Owner) // sql.Result

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) GetAllCars() (*model.CarCatalog, error) {
	catalog := model.CarCatalog{}
	query := `select reg_num, mark, model, year, owner from car_catalog`
	rows, err := s.ConnectionDB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		car := model.Car{}

		if err := rows.Scan(
			&car.RegNum,
			&car.Mark,
			&car.Model,
			&car.Year,
			&car.Owner,
		); err != nil {
			return nil, err
		}

		catalog.Cars = append(catalog.Cars, &car)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &catalog, nil
}

func (s *Storage) GetCarsByFilter(options string) (*model.CarCatalog, error) {
	catalog := model.CarCatalog{}
	query := `select reg_num, mark, model, year from car_catalog where ` + options
	s.Logger.Info(query)
	rows, err := s.ConnectionDB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		car := model.Car{}

		if err := rows.Scan(
			&car.RegNum,
			&car.Mark,
			&car.Model,
			&car.Year,
		); err != nil {
			return nil, err
		}

		catalog.Cars = append(catalog.Cars, &car)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &catalog, nil
}
