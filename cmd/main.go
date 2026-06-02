package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Ajulll22/payment-ai-assistant/internal/version"
	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"github.com/gin-gonic/gin"
	"github.com/judwhite/go-svc"
	"gorm.io/gorm"
)

type program struct {
	db   *gorm.DB
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

	if p.db != nil {
		database.CloseDB(p.dbApp)
	}

	if p.logFile != nil {
		p.logFile.Close()
	}

	log.Println("Stop: complete")
	return nil
}
