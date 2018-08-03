package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/motionwerkGmbH/msp-backend-api/configs"
	"github.com/motionwerkGmbH/msp-backend-api/tools"
	"net/http"
	"strconv"
	"time"
)

var router *gin.Engine

func main() {

	Config := configs.Load()

	if (Config.GetString("environment")) == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.New()
	router.Use(gin.Recovery())

	// allow all origins
	router.Use(cors.Default())

	InitializeRoutes()

	// Establish database connection
	tools.Connect("_theDb.db")
	tools.MySQLConnect("blockchain")

	log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	log.Println("~~~~~~~~~~~~~~~~~~ MSP BACKEND RUNNING ~~~~~~~~~~~~~~~~~")
	log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	// Serve 'em...
	server := &http.Server{
		Addr:           Config.GetString("hostname") + ":" + strconv.Itoa(Config.GetInt("port")),
		Handler:        router,
		ReadTimeout:    200 * time.Second,
		WriteTimeout:   200 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.SetKeepAlivesEnabled(false)
	server.ListenAndServe()

}
