package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Smer4k/slow-web-calculator/internal/datatypes"
	_ "github.com/mattn/go-sqlite3"
)

// создает базу данных если ее нету
func InitDataBase() {
	db, err := sql.Open("sqlite3", "../../internal/database/database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Создаем таблицу expressions если ее нету
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS expressions (
			id TEXT PRIMARY KEY,
			JsonData TEXT,
			timeForSolve INTEGER,
			answer INTEGER,
			status TEXT
		)
	`)
	if err != nil {
		panic(err)
	}
	// Создаем таблицу timeexecution если ее нету
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS timeexecution (
		id TEXT PRIMARY KEY,
		timeExec INTEGER
	);
`)
	if err != nil {
		panic(err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM timeexecution").Scan(&count)
	if err != nil {
		panic(err)
	}

	if count == 0 { // Добавляются значения по умолчанию если таблица новая
		_, err = db.Exec(`
		INSERT INTO timeexecution (id, timeExec) VALUES ('time_sum', 20);
		INSERT INTO timeexecution (id, timeExec) VALUES ('time_subtraction', 20);
		INSERT INTO timeexecution (id, timeExec) VALUES ('time_multi', 20);
		INSERT INTO timeexecution (id, timeExec) VALUES ('time_division', 20);
		INSERT INTO timeexecution (id, timeExec) VALUES ('time_out', 20);
	`)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Таблицы успешно созданы в базе данных или уже существуют.")
}

// добавляет выражение в базу данных
func AddExpression(id string, data *datatypes.Expression, timeForSolve int, status string) error {
	db, err := sql.Open("sqlite3", "../../internal/database/database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	if data == nil {
		_, err = db.Exec("INSERT INTO expressions (id, JsonData, timeForSolve, answer, status) VALUES (?, ?, ?, ?, ?)",
			id, "", timeForSolve, 0, status)
		if err != nil {
			return err
		}
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO expressions (id, JsonData, timeForSolve, answer, status) VALUES (?, ?, ?, ?, ?)",
		id, jsonData, timeForSolve, 0, status)

	if err != nil {
		return err
	}
	return nil
}

// делает запрос в базу данных и возвращает все выражения которые нужно решить
func GetWorkExpressionsData() (map[string]datatypes.Expression, error) {
	db, err := sql.Open("sqlite3", "../../internal/database/database.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM expressions WHERE status = ?", "work").Scan(&count)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT id, JsonData FROM expressions WHERE status = ?", "work")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	newList := make(map[string]datatypes.Expression)

	for rows.Next() {
		var id, jsonData string
		err := rows.Scan(&id, &jsonData)
		if err != nil {
			return nil, err
		}

		var expr datatypes.Expression
		err = json.Unmarshal([]byte(jsonData), &expr)
		if err != nil {
			return nil, err
		}
		newList[id] = expr
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return newList, nil
}

// проверяет есть ли данное выражение в базе данных
func ContainsExpression(expr string) (bool, error) {
	db, err := sql.Open("sqlite3", "../../internal/database/database.db")
	if err != nil {
		return false, err
	}
	defer db.Close()

	var idFromDB string
	err = db.QueryRow("SELECT id FROM expressions WHERE id = ?", expr).Scan(&idFromDB)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}


// Обновляет данные о настройка в базе данных
func UpdateSettingsData(data map[datatypes.NameTimeExec]int) error {
	db, err := sql.Open("sqlite3", "../../internal/database/database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("REPLACE INTO timeexecution (id, timeExec) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, val := range data {
		_, err := stmt.Exec(string(key), val)
		if err != nil {
			return err
		}
	}
	return nil
}

// Запрашивает и возвращает данные о настройках
func GetSettingsData() (map[datatypes.NameTimeExec]int, error) {
	db, err := sql.Open("sqlite3", "../../internal/database/database.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM timeexecution")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dataSettings := make(map[datatypes.NameTimeExec]int)

	for rows.Next() {
		var name string
		var time int
	
		err := rows.Scan(&name, &time)
		if err != nil {
			return nil, err
		}
		dataSettings[datatypes.NameTimeExec(name)] = time
	}
	return dataSettings, nil
}
