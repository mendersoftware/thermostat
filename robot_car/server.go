package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
)

func CarInitiator() gin.HandlerFunc {
	socketServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	testMode := os.Getenv("CAR_TEST")
	test := false
	if testMode != "" {
		test = true
	}
	car := NewKunmanCar(test)

	return func(c *gin.Context) {
		c.Set("socketServer", socketServer)
		c.Set("car", car)
		c.Next()
	}
}

func serve() {
	router := gin.Default()

	router.Use(CarInitiator())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	router.Use(cors.New(config))

	router.StaticFile("/", "./frontend")
	router.Static("/js", "./frontend/dist")
	// router.Static("/modules", "./frontend/node_modules")

	router.GET("/socket.io", socketHandler)
	router.POST("/socket.io", socketHandler)
	router.Handle("WS", "/socket.io", socketHandler)
	router.Handle("WSS", "/socket.io", socketHandler)

	router.Run(":8080")
}

func socketHandler(c *gin.Context) {
	socketServer := c.MustGet("socketServer").(*socketio.Server)
	car := c.MustGet("car").(*KunmanCar)

	socketServer.On("move", func(data string) {
		type Move struct {
			Force     float32
			Direction string
			Interval  int
		}

		var step Move
		err := json.Unmarshal([]byte(data), &step)
		if err != nil {
			log.Printf("can not parse data: %s", err)
		}
		fmt.Printf("got data: %v\n", step)

		if step.Direction == "" {
			car.stop()
		}
		if step.Force > 0.1 {
			car.changeSpeed(int(step.Force * 50))
			car.move(step.Direction, step.Interval)
		}
	})

	socketServer.ServeHTTP(c.Writer, c.Request)
}
