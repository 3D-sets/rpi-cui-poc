package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/googolgl/go-i2c"
	"github.com/googolgl/go-pca9685"
)

var servo *pca9685.Servo
var motor *pca9685.Servo
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

func setupServo() {
	i2c, err := i2c.New(pca9685.Address, 1)
	if err != nil {
		log.Fatal(err)
	}

	pca0, err := pca9685.New(i2c, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Sets a single PWM channel 0
	pca0.SetChannel(0, 0, 130)

	// Servo on channel 0
	servo = pca0.ServoNew(0, nil)

	pca0.SetChannel(1, 0, 130)
	motor = pca0.ServoNew(1, nil)

}

func setAngle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/servo" {
		http.NotFound(w, r)
		return
	}
	fmt.Println("Endpoint Hit: setAngle")
	keys, ok := r.URL.Query()["angle"]
	angle, _ := strconv.Atoi(keys[0])
	angle += 45

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'angle' is missing")
		return
	}
	servo.Angle(angle)

	fmt.Printf("New Angle Set: %s\n", keys[0])
}

func setSpeed(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/motor" {
		http.NotFound(w, r)
		return
	}
	fmt.Println("Endpoint Hit: setSpeed")
	keys, ok := r.URL.Query()["speed"]
	speed, _ := strconv.Atoi(keys[0])

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'speed' is missing")
		return
	}
	motor.Angle(speed)

	fmt.Printf("New Speed Set: %s\n", keys[0])
}

func handleRequests() {
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/servo", setAngle)
	http.HandleFunc("/motor", setSpeed)
	log.Fatal(http.ListenAndServe("0.0.0.0:10000", nil))
}

func main() {
	fileServer := http.FileServer(http.Dir("./ui"))
	http.Handle("/", fileServer)
	setupServo()
	servo.Fraction(0.5)
	motor.Fraction(0.5)

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
