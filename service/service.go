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

type Cities struct {
	Cities []City `json:"cities"`
}

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

// Return all Cities from database
func GetCities() (Cities, error) {
	var cities Cities
	rows, err := db.DB.Query("SELECT * FROM cities")
	if err != nil {
		return cities, err
	}
	for rows.Next() {
		var city City
		_ = rows.Scan(&city.ID, &city.Name)

		err := getCityBorders(&city)
		if err != nil {
			return cities, err
		}

		cities.Cities = append(cities.Cities, city)
	}
	return cities, nil
}

// Get a `City` object from database using `id`
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

// Create a `City` object on database
func CreateCity(city *City) error {
	stmt, _ := db.DB.Prepare("INSERT INTO cities(name) values (?)")
	result, err := stmt.Exec(city.Name)
	if err != nil {
		return err
	}
	city.ID, _ = result.LastInsertId()
	return nil
}

// Update a `City` object on database
func UpdateCity(city *City) error {
	stmt, _ := db.DB.Prepare("UPDATE cities SET name = ? WHERE id = ?")
	_, err := stmt.Exec(city.Name, city.ID)
	return err
}

// Insert the `Borders` of a `City` on database
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

// Remove all `Borders` of a `City` from the database
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

// Remove a specific `City` given the id
func RemoveCity(cityID int64) error {
	stmt, _ := db.DB.Prepare("DELETE FROM cities WHERE id = ?")
	_, err := stmt.Exec(cityID)
	return err
}

// Remove all city and border data from database
func RemoveCities() error {
	_, err := db.DB.Exec("DELETE FROM cities")
	if err != nil {
		return err
	}
	_, errBorders := db.DB.Exec("DELETE FROM borders")
	return errBorders
}

// Find recursively a valid path and update the given Path pointer
func findPath(path *Path, start int64, end int64) error {
	stmt, _ := db.DB.Prepare("SELECT `to` FROM borders WHERE `from` = ?")
	rows, err := stmt.Query(start)
	if err != nil {
		return err
	}
	first := true
	for rows.Next() {
		if first {
			first = false
			path.Borders = append(path.Borders, start)
		}
		var border int64
		_ = rows.Scan(&border)
		if border == end {
			path.Borders = append(path.Borders, end)
			return nil
		}
		return findPath(path, border, end)
	}
	return errors.New("No path found")
}

// Return a valid Path from a City to another
func GetPath(fromId int64, toId int64) (Path, error) {
	var path Path
	err := findPath(&path, fromId, toId)
	if err != nil {
		return path, err
	}
	return path, nil
}
