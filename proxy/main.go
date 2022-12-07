// Запуск обратного прокси
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// адрес прокси
const proxy string = "localhost:8800"

// адреса серверов
var (
	counter int    = 0              //переключатель
	host1   string = "http://:8080" //1-й сервер
	host2   string = "http://:8081" //2-й сервер
)

func main() {

	http.HandleFunc("/", reverseProxy)
	log.Printf("   PROXY по адресу:  %s", proxy)

	log.Fatalln(http.ListenAndServe(proxy, nil))
}

// обработчик обратного прокси для 2-х серверов
func reverseProxy(w http.ResponseWriter, r *http.Request) {
	if counter == 0 {
		urlFirst, err := url.Parse(host1)
		if err != nil {
			log.Fatalln(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(urlFirst)
		proxy.ServeHTTP(w, r)
		log.Printf("доступен: %s%s", host1, r.URL)
		counter++
		return
	}

	urlSecond, err := url.Parse(host2)
	if err != nil {
		log.Fatalln(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(urlSecond)
	proxy.ServeHTTP(w, r)
	log.Printf("доступен: %s%s", host2, r.URL)
	counter--
}
