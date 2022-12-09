package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func IndexPage(c *gin.Context) { //смотрим
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Home Page",
		"Any":   "Что-то",
	})
	c.String(200, "Выполнено")
}

// 1. Вернуть весь список
func getAllPerson(c *gin.Context) {
	users, err := Latest()
	if err != nil {
		c.String(404, "ошибка извлечения данных")
		return
	}
		
	// Установить статус HTTP 200 (OK)
	c.HTML(http.StatusOK, "home.html", gin.H{
		"Payload": users,
	})
}

// 2. Вернуть по ID
func getPerson(c *gin.Context) {
	if ID, err := strconv.Atoi(c.Param("id")); err == nil {
		user := viewById(ID) //если есть ID?
		// визуализация шаблона HTML
		c.HTML(http.StatusOK, "user.html", gin.H{
			"payload": user, //Передать данные
		})
	} else { // Если не верный ID, то ошибка
		c.String(404, "Неверный ID")
	}
}

// 3. Вставить строку в MYSQL
func postPerson(c *gin.Context) {
	//получить данные "json" из запроса
	//var newuser NewUser
	if err := c.Bind(&newUser); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		Insert(newUser) //вставить в таблицу MYSQL
		c.IndentedJSON(http.StatusCreated, newUser)
	}
}

// 4. Удалить строку в MYSQL
func deletePerson(c *gin.Context) {
	//получить параметр из запроса
	if ID, err := strconv.Atoi(c.Param("id")); err == nil {
		deleteById(ID) //удалить строку таблицы MYSQL
		c.String(200, "ID удален")
	} else {
		c.String(400, "Неверный ID")
	}
}

// 5. Изменить поле столбца в MYSQL
func putAge(c *gin.Context) {
	if id, err := strconv.Atoi(c.Param("id")); err == nil {
		//получить данные "json" из запроса
		if err := c.BindJSON(&newAge); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			Change(newAge, (id)) //изменить поле в строке MYSQL
			c.String(200, "Age изменен")
		}
	} else {
		c.String(400, "Неверный ID")
	}
}

// 6. Делает друзей из двух пользователей
var union Friendship

func postFriends(c *gin.Context) {
	//получить данные "json" из запроса
	if err := c.BindJSON(&union); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		Friend(union) //изменить поле в строке MYSQL

		union.SourceId, union.TargetId = union.TargetId, union.SourceId
		Friend(union) //изменить поле в строке MYSQL
		c.String(201, "SourceId и TargetId теперь друзья")
	}
}

// 7. Вернуть друзей пользователя по ID
func getFriends(c *gin.Context) {
	if ID, err := strconv.Atoi(c.Param("id")); err == nil {
		//если есть ID, получаем строку из MYSQL
		f := viewFriends(ID)
		if err != nil {
			c.String(404, "ошибка извлечения данных")
			return
		}
		ID := f.SourceID
		SourceN := f.SourceName
		FriendsAll := string(f.AllFriends)
		wordsFriendsAll := strings.Split(FriendsAll, ",")
		// визуализация шаблона HTML
		c.HTML(http.StatusOK, "friends.html", gin.H{
			//Передать данные в форму HTML
			"Id":      ID,
			"Name":    SourceN,
			"Payload": wordsFriendsAll,
		})
	} else { // Если не верный ID, то ошибка
		c.String(404, "Неверный ID")
	}
}
