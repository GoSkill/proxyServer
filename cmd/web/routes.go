package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes() http.Handler {

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"}) //доверенный IP-адрес

	//Загрузка шаблонов (Рендеринг HTML)
	r.LoadHTMLGlob("templates/*")

	//Маршруты:
	r.GET("/", IndexPage)
	r.GET("/users", getAllPerson)
	r.GET("/users/view/:id", getPerson)
	r.POST("/create", postPerson)
	r.POST("/delete/:id", deletePerson)
	r.POST("/change/:id", putAge)
	r.POST("/friendship", postFriends)
	r.GET("/friends/:id", getFriends)
	return r
}
