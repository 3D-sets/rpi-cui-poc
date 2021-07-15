package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

type GPIODevice struct {
	Id    string `json:"Id"`
	Desc  string `json:"Desc"`
	State bool   `json:"State"`
}

var GPIODevices []GPIODevice
var ledChannel1 chan bool

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
	err = ts.Execute(w, GPIODevices)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func returnAllDevices(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Return All Devices")
	json.NewEncoder(w).Encode(GPIODevices)
}

func toggleLED1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Toggle LED1")
	GPIODevices[0].State = !GPIODevices[0].State
	fmt.Println("State Changed: Toggle LED1")
	json.NewEncoder(w).Encode(GPIODevices[0])
	ledChannel1 <- GPIODevices[0].State
}

func handleRequests() {
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/devices", returnAllDevices)
	http.HandleFunc("/devices/led1", toggleLED1)
	log.Fatal(http.ListenAndServe("0.0.0.0:10000", nil))
}

func main() {
	fileServer := http.FileServer(http.Dir("./ui"))
	http.Handle("/", fileServer)

	GPIODevices = []GPIODevice{
		GPIODevice{Id: "LED1", Desc: "RPI GPIO4 LED", State: false},
		GPIODevice{Id: "LED2", Desc: "RPI GPIO24 LED", State: false},
	}

	ledChannel1 = make(chan bool)

	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	values := map[int]string{0: "inactive", 1: "active"}
	offset := rpi.GPIO4
	v := 0
	l, err := c.RequestLine(offset, gpiod.AsOutput(v))
	if err != nil {
		panic(err)
	}
	defer func() {
		l.Reconfigure(gpiod.AsInput)
		l.Close()
	}()
	fmt.Printf("Set pin %d %s\n", offset, values[v])

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	go handleRequests()
	for {
		select {
		case <-ledChannel1:
			v ^= 1
			l.SetValue(v)
			fmt.Printf("Set pin %d %s\n", offset, values[v])
		case <-quit:
			return
		}
	}

}
