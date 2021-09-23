package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var sleep chan bool

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}
	fmt.Println("Endpoint Hit: homePage")

	ts, err := template.ParseFiles("./ui/html-tmpl/home.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func setAngle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/servo" {
		http.NotFound(w, r)
		return
	}
	fmt.Println("Endpoint Hit: setAngle")
	keys, ok := r.URL.Query()["angle"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'angle' is missing")
		return
	}

	fmt.Printf("New Angle Set: %s\n", keys[0])
}

func handleRequests() {
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/servo", setAngle)
	log.Fatal(http.ListenAndServe("0.0.0.0:10000", nil))
}

func main() {
	fileServer := http.FileServer(http.Dir("./ui"))
	http.Handle("/", fileServer)

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	go handleRequests()
	for {
		select {
		case <-sleep:
			time.Sleep(2 * time.Millisecond)
		case <-quit:
			return
		}
	}

}
