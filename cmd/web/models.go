package main

import (
	_ "database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct { //поля структуры должны соответствовать таблице MYSQL
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
	//	NewAge int `json:"newage"`
}

type Friends struct { //дубль-таблица Person для MYSQL
	SourceID   int     `json:"sourceid"`
	SourceName string  `json:"sourcename"`
	AllFriends []uint8 `json:"allfriends"`
}

type Friendship struct { // Структура для запроса дружбы
	FriendID int `json:"friendid"`
	SourceId int `json:"sourceid"` //ID инициатора дружбы
	TargetId int `json:"targetid"` //ID принявшего запрос
}

var (
	newUser Person
	users   []Person
)

// 1. "GET" запрос нескольких (10)строк
func Latest() ([]Person, error) {
	//создаем оператор MYSQL
	stmt := `SELECT PersonID, FirstName, LastName, Age FROM Persons ORDER BY PersonID LIMIT 10`
	//Получаем строки для их итерации
	rows, err := DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//инициализируем обнуленную структуру и слайс
	s := Person{}
	users = nil
	fmt.Println(users)
	for rows.Next() { //итерация строк...
		//копируем значения каждой записи базы в поля структуры "s"
		err = rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Age)
		if err != nil {
			return nil, err
		}
		users = append(users, s) //добавляем их в хранилище
	}
	//после завершения цикла вызываем rows.Err() для получения любой ошибки, возникшей во время итерации.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil //возвращаем полученные данные
}

// 2. "GET" запрос одной строки
func viewById(id int) Person {
	stmt := `SELECT PersonID, FirstName, LastName, Age FROM Persons WHERE PersonID = ?`
	row := DB.QueryRow(stmt, id)
	s := Person{}
	//копируем значения в поля структуры "Person"
	_ = row.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Age)
	return s
}

// 3. "POST" запрос вставки одной строки
func Insert(u Person) {
	stmtP := `INSERT INTO Persons (FirstName, LastName, Age) VALUES(?, ?, ?)`
	_, err := DB.Exec(stmtP, &u.FirstName, &u.LastName, &u.Age)
	if err != nil {
		log.Println(err)
	}
	//метод DB.Exec() используют только для вставки (INSERT) или удаления (DELETE)строк
	//дублируем таблицу "Persons == Users" (требуется для работы с друзьями)
	stmtU := `INSERT INTO Users (FirstN, LastN, AgeU) VALUES(?, ?, ?)`
	_, err = DB.Exec(stmtU, &u.FirstName, &u.LastName, &u.Age)
	if err != nil {
		log.Println(err)
	}
}

// 4. "DELETE" запрос удаления одной строки
func deleteById(id int) {
	stmtF := `DELETE FROM Friendship WHERE SourceID = ? OR TargetID = ?`
	_, err := DB.Exec(stmtF, id, id)
	if err != nil {
		log.Println(err)
	}
	stmt := `DELETE FROM Persons, Users USING Persons INNER JOIN Users ON Persons.PersonID = Users.UserID WHERE Persons.PersonID = ?`
	_, err = DB.Exec(stmt, id)
	if err != nil {
		log.Println(err)
	}
}

// 5. "PUT" запрос изменения поля строки
func Change(newAge int, id int) {
	stmtP := `UPDATE Persons SET Age = ? WHERE PersonID = ?`
	_, err := DB.Exec(stmtP, newAge, id)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("models %v %v\n", newAge, id)
	//корректируем таблицу "Users" в соответствии с "Persons"
	stmtU := `UPDATE Users SET AgeU = ? WHERE UserID = ?`
	_, err = DB.Exec(stmtU, newAge, id)
	if err != nil {
		log.Println(err)
	}
}

// 6. "POST" запрос создания дружбы
func Friend(f Friendship) {
	//создаем оператор MYSQL
	stmt := `INSERT INTO Friendship (SourceID, TargetID) VALUES(?, ?)`
	_, err := DB.Exec(stmt, &f.SourceId, &f.TargetId)
	if err != nil {
		log.Println(err)
	}
}

// 7. "GET" запрос друзей по ID пользователя
func viewFriends(id int) Friends { //[]uint8 {
	//SQL-представление (виртуальная таблица из 3-х источников)
	//Удалить таблицу, если есть
	_ = `DROP VIEW IF EXISTS Persons_Friendship_Summary`
	//создать виртуальную таблицу
	_ = `CREATE VIEW Persons_Friendship_Summary AS SELECT PersonID AS pfs_ID, max(FirstName) AS pfs_FirstName, 
	group_concat(LastN ORDER BY LastN SEPARATOR ',') AS pfs_Friend_array FROM Persons INNER JOIN Friendship ON 
	Persons.PersonID = Friendship.SourceID INNER JOIN Users ON Friendship.TargetID = Users.UserID GROUP BY Persons.PersonID`
	//создаем оператор MYSQL получения друзей по ID
	stmt := `SELECT pfs_ID, pfs_FirstName, pfs_Friend_array FROM Persons_Friendship_Summary WHERE pfs_ID = ?`
	row := DB.QueryRow(stmt, id)
	f := Friends{}
	_ = row.Scan(&f.SourceID, &f.SourceName, &f.AllFriends)
	return f
}
