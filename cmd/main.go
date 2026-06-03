package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Ajulll22/payment-ai-assistant/internal/middleware"
	"github.com/Ajulll22/payment-ai-assistant/internal/route"
	"github.com/Ajulll22/payment-ai-assistant/internal/version"
	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"github.com/Ajulll22/payment-ai-assistant/pkg/database"
	"github.com/Ajulll22/payment-ai-assistant/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/judwhite/go-svc"
	"gorm.io/gorm"
)

type program struct {
	dbApp   *gorm.DB
	dbDWH   *gorm.DB
	cfg     *constant.Config
	server  *gin.Engine
	logFile *os.File
}

func main() {
	showVersion := flag.Bool("version", false, "Show build version info")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\nCommit: %s\nBuilt: %s\n", version.Version, version.Commit, version.BuildTime)
		os.Exit(0)
	}

	p := &program{}
	if err := svc.Run(p); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	appDebugFile, err := os.OpenFile("./untracked.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	p.logFile = appDebugFile

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(appDebugFile)

	log.Println("Init: starting service setup")

	cfg := constant.GetEnv()
	p.cfg = cfg

	dbApp, err := database.SQLConnect(database.SQLConfig{
		User:            cfg.DBApp.User,
		PasswordEnc:     cfg.DBApp.Password,
		Host:            cfg.DBApp.Host,
		Port:            cfg.DBApp.Port,
		Name:            cfg.DBApp.Name,
		Timeout:         cfg.DBApp.Timeout,
		LogDir:          cfg.Log.Path + "logs/db",
		LogMaxFile:      cfg.DBApp.LogMaxFile,
		FallbackLogFile: p.logFile,
		AppKey:          cfg.App.Key,
	})
	if err != nil {
		log.Printf("database connection failed: %v", err)
		return err
	}
	dbDWH, err := database.SQLConnect(database.SQLConfig{
		User:            cfg.DBDWH.User,
		PasswordEnc:     cfg.DBDWH.Password,
		Host:            cfg.DBDWH.Host,
		Port:            cfg.DBDWH.Port,
		Name:            cfg.DBDWH.Name,
		Timeout:         cfg.DBDWH.Timeout,
		LogDir:          cfg.Log.Path + "logs/db",
		LogMaxFile:      cfg.DBDWH.LogMaxFile,
		FallbackLogFile: p.logFile,
		AppKey:          cfg.App.Key,
	})
	if err != nil {
		log.Printf("database connection failed: %v", err)
		return err
	}

	p.dbApp = dbApp
	p.dbDWH = dbDWH

	validation.RegisterCustomValidation()

	app := gin.Default()
	app.Use(middleware.LoggingMiddleware(cfg.Log))
	app.Use(middleware.RecoveryMiddleware(cfg.Log))
	app.Use(middleware.SetIPMiddleware())

	route.Register(app, dbApp, dbDWH, cfg)
	p.server = app

	log.Println("Init: complete")
	return nil
}

func (p *program) Start() error {
	log.Printf("Start: launching gin HTTP server on :%s", p.cfg.App.Port)

	go func() {
		if err := p.server.Run(":" + p.cfg.App.Port); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()
	return nil
}

func (p *program) Stop() error {
	log.Println("Stop: closing all connection")

	if p.dbApp != nil {
		database.CloseDB(p.dbApp)
	}
	if p.dbDWH != nil {
		database.CloseDB(p.dbDWH)
	}

	if p.logFile != nil {
		p.logFile.Close()
	}

	log.Println("Stop: complete")
	return nil
}
