// Запуск обратного прокси
package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//go run ./cmd/web -addr="127.0.0.1:8085" -dsn="web:1@/proxyserver?parseTime=true"
//go run ./cmd/web -addr="127.0.0.1:8086"
//go run ./proxy -host1="127.0.0.1:8085" -host2="http://:8086"

// адрес прокси
const proxy string = "localhost:8800"

var counter int = 0

func main() {

	host1 := flag.String("host1", "http://:8080", "Сетевой адрес HTTP")
	host2 := flag.String("host2", "http://:8081", "Сетевой адрес HTTP")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if counter == 0 {
			urlFirst, err := url.Parse(*host1)
			if err != nil {
				log.Fatalln(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(urlFirst)
			proxy.ServeHTTP(w, r)
			log.Printf("доступен: %s%s", *host1, r.URL)
			counter++
			return
		}

		urlSecond, err := url.Parse(*host2)
		if err != nil {
			log.Fatalln(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(urlSecond)
		proxy.ServeHTTP(w, r)
		log.Printf("доступен: %s%s", *host2, r.URL)
		counter--
	})

	log.Printf("   PROXY по адресу:  %s", proxy)

	log.Fatalln(http.ListenAndServe(proxy, nil))
}
