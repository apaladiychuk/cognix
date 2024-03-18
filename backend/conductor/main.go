package main

import (
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/zap"
)

func main() {
	cfg, err := ReadConfig()
	if err != nil {
		zap.S().Errorf("read log %s:", err.Error())
		return
	}
	utils.InitLogger(cfg.Debug)
	db, err := repository.NewDatabase(cfg.DB)
	if err != nil {
		utils.Logger.Errorf("Init database connection: %s", err.Error())
		return
	}
	// repositories
	connectorRepo := repository.NewConnectorRepository(db)
	conductor := NewConductor(connectorRepo)
	conductor.Start()
}
