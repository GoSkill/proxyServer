package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes() http.Handler {

	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"}) //доверенный IP-адрес

	//Загрузка шаблонов (Рендеринг HTML)
	router.LoadHTMLGlob("templates/*")

	//Маршруты:
	router.GET("/", IndexPage)
	router.GET("/users", getAllPerson)
	router.GET("/users/view/:id", getPerson)
	router.POST("/users/create", postPerson)         //$ curl -X POST -i http://localhost:8800/users/create -H "content-type: application/json" -d "{\"firstname\":\"Ivan\", \"lastname\":\"Pupkin\", \"age\":33}"
	router.DELETE("/users/delete/:id", deletePerson) //$ curl -X DELETE -i http://localhost:8800/users/delete/4
	router.PUT("/users/:id", putAge)                 //$ curl -X PUT -H "content-type: application/json" -d "22" -i http://localhost:8800/users/2
	router.POST("/users/friends", postFriends)       //$ curl -X POST -i http://localhost:8800/users/friends -H "content-type: application/json" -d "{\"sourceid\":1,\"targetid\":9}"
	router.GET("/users/friends/:id", getFriends)
	return router
}