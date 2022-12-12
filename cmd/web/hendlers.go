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
			"payload": user, //Передать данные  шаблон
		})
	} else {
		c.String(404, "Неверный ID")
	}
}

// 3. Вставить строку в MYSQL
func postPerson(c *gin.Context) {
	if err := c.Bind(&newUser); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		Insert(newUser) //вставить в таблицу MYSQL
	}
	c.HTML(http.StatusOK, "create.html", gin.H{
		"new": newUser,
	})
}

// 4. Удалить строку в MYSQL
func deletePerson(c *gin.Context) {
	if id, err := strconv.Atoi(c.Param("id")); err == nil {
		deleteById(id) //удалить строку таблицы MYSQL
		c.HTML(http.StatusOK, "delete.html", gin.H{
			"del": id,
		})
	} else {
		c.String(400, "Неверный ID")
	}
}

// 5. Изменить поле столбца в MYSQL
func putAge(c *gin.Context) {
	if id, err := strconv.Atoi(c.Param("id")); err == nil {
		if err := c.ShouldBind(&newUser); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			Change(newUser.Age, id)
			c.HTML(http.StatusOK, "change.html", gin.H{
				"change": newUser,
				"id":     id,
			})
		}
	} else {
		c.String(400, "Неверный ID")
	}
}

// 6. Делает друзей из двух пользователей
var union Friendship

func postFriends(c *gin.Context) {
	if err := c.Bind(&union); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		Friend(union) //первая запись в строке MYSQL
		union.SourceId, union.TargetId = union.TargetId, union.SourceId
		Friend(union) //вторая запись в строке MYSQL
		c.HTML(http.StatusOK, "friendship.html", gin.H{
			"ID1": union.SourceId,
			"ID2": union.TargetId,
		})
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
		c.HTML(http.StatusOK, "friends.html", gin.H{
			"Id":      ID,
			"Name":    SourceN,
			"Payload": wordsFriendsAll,
		})
	} else {
		c.String(404, "Неверный ID")
	}
}
