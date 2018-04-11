package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"../types"
)

var mainDB *sql.DB

func InitDatabase(databaseName string) error {
	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return err
	}
	mainDB = db
	_, errCities := mainDB.Exec(
		"CREATE TABLE IF NOT EXISTS cities (id INTEGER PRIMARY KEY, name TEXT)",
	)
	if errCities != nil {
		return errCities
	}
	_, errBorders := mainDB.Exec(
		"CREATE TABLE IF NOT EXISTS borders (`from` integer, `to` integer)",
	)
	return errBorders
}

func getCityBorders(city *types.City) error {
	var borders []int64
	stmt, _ := mainDB.Prepare("SELECT `to` FROM borders WHERE `from` = ?")
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

func GetCities() (types.Cities, error) {
	rows, err := mainDB.Query("SELECT * FROM cities")
	if err != nil {
		return nil, err
	}
	var cities types.Cities
	for rows.Next() {
		var city types.City
		_ = rows.Scan(&city.ID, &city.Name)

		err := getCityBorders(&city)
		if err != nil {
			return nil, err
		}

		cities = append(cities, city)
	}
	return cities, nil
}

func GetCity(cityId int64) (types.City, error) {
	stmt, _ := mainDB.Prepare("SELECT * FROM cities where id = ?")
	rows, err := stmt.Query(cityId)
	var city types.City
	if err != nil {
		return city, err
	}
	for rows.Next() {
		_ = rows.Scan(&city.ID, &city.Name)

		err := getCityBorders(&city)
		if err != nil {
			return types.City{}, err
		}
	}
	return city, nil
}

func CreateCity(city *types.City) error {
	stmt, _ := mainDB.Prepare("INSERT INTO cities(name) values (?)")
	result, err := stmt.Exec(city.Name)
	if err != nil {
		return err
	}
	city.ID, _ = result.LastInsertId()
	return nil
}

func UpdateCity(city *types.City) error {
	stmt, _ := mainDB.Prepare("UPDATE cities SET name = ? WHERE id = ?")
	_, err := stmt.Exec(city.Name, city.ID)
	return err
}

func InsertCityBorders(city *types.City) error {
	var borders []int64
	for _, value := range city.Borders {
		_, err := mainDB.Exec(
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

func RemoveCityBorders(city *types.City) error {
	stmt, _ := mainDB.Prepare("DELETE FROM borders WHERE from = ?")
	_, err := stmt.Exec(city.ID)
	if err != nil {
		return err
	}
	var borders []int64
	city.Borders = borders
	return nil
}

func RemoveCity(cityID int64) error {
	stmt, _ := mainDB.Prepare("DELETE FROM cities WHERE id = ?")
	_, err := stmt.Exec(cityID)
	return err
}

func RemoveCities() error {
	_, err := mainDB.Exec("DELETE FROM cities")
	return err
}
