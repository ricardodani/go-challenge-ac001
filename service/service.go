package service

import (
	"errors"

	"../db"
)

type City struct {
	ID      int64   `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Borders []int64 `json:"borders,omitempty"`
}

type Cities []City

type Path struct {
	Path []int64 `json:path,omitempty`
}

func NewCity() City {
	var city City
	return city
}

func getCityBorders(city *City) error {
	var borders []int64
	stmt, _ := db.DB.Prepare("SELECT `to` FROM borders WHERE `from` = ?")
	rows, err := stmt.Query(city.ID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var border int64
		_ = rows.Scan(&border)
		borders = append(borders, border)
	}
	city.Borders = borders
	return nil
}

func GetCities() (Cities, error) {
	rows, err := db.DB.Query("SELECT * FROM cities")
	if err != nil {
		return nil, err
	}
	var cities Cities
	for rows.Next() {
		var city City
		_ = rows.Scan(&city.ID, &city.Name)

		err := getCityBorders(&city)
		if err != nil {
			return nil, err
		}

		cities = append(cities, city)
	}
	return cities, nil
}

func GetCity(cityId int64) (City, error) {
	stmt, _ := db.DB.Prepare("SELECT * FROM cities where id = ?")
	rows, err := stmt.Query(cityId)
	var city City
	if err != nil {
		return city, err
	}
	for rows.Next() {
		_ = rows.Scan(&city.ID, &city.Name)

		err := getCityBorders(&city)
		if err != nil {
			return city, err
		}
		return city, nil
	}
	return city, errors.New("City not found")
}

func CreateCity(city *City) error {
	stmt, _ := db.DB.Prepare("INSERT INTO cities(name) values (?)")
	result, err := stmt.Exec(city.Name)
	if err != nil {
		return err
	}
	city.ID, _ = result.LastInsertId()
	return nil
}

func UpdateCity(city *City) error {
	stmt, _ := db.DB.Prepare("UPDATE cities SET name = ? WHERE id = ?")
	_, err := stmt.Exec(city.Name, city.ID)
	return err
}

func InsertCityBorders(city *City) error {
	var borders []int64
	for _, value := range city.Borders {
		_, err := db.DB.Exec(
			"INSERT INTO borders VALUES ($1, $2)",
			city.ID,
			value,
		)
		if err != nil {
			return err
		}
		borders = append(borders, value)
	}
	city.Borders = borders
	return nil
}

func RemoveCityBorders(city *City) error {
	stmt, _ := db.DB.Prepare("DELETE FROM borders WHERE `from` = ?")
	_, err := stmt.Exec(city.ID)
	if err != nil {
		return err
	}
	var borders []int64
	city.Borders = borders
	return nil
}

func RemoveCity(cityID int64) error {
	stmt, _ := db.DB.Prepare("DELETE FROM cities WHERE id = ?")
	_, err := stmt.Exec(cityID)
	return err
}

func RemoveCities() error {
	_, err := db.DB.Exec("DELETE FROM cities")
	if err != nil {
		return err
	}
	_, errBorders := db.DB.Exec("DELETE FROM borders")
	return errBorders
}
