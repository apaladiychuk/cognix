package main

import (
	"cognix.ch/api/v2/api/handler"
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/utils"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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
	oauthProxy := oauth.NewGoogleProvider(cfg.OAuth, cfg.RedirectURL)
	jwtService := security.NewJWTService(cfg.JWTSecret, cfg.JWTExpiredTime)
	// repositories
	connectorRepo := repository.NewConnectorRepository(db)
	userRepo := repository.NewUserRepository(db)
	// handlers
	authHandler := handler.NewAuthHandler(oauthProxy, jwtService, userRepo)
	connectorHandler := handler.NewCollectorHandler(connectorRepo)

	router := NewRouter()

	connectorHandler.Mount(router, nil)
	authHandler.Mount(router)
	RunServer(cfg, router)
}

func NewRouter() *gin.Engine {
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))
	return router
}

func RunServer(cfg *Config, router *gin.Engine) {
	srv := http.Server{}
	srv.Addr = fmt.Sprintf(":%d", cfg.Port)
	srv.Handler = router
	utils.Logger.Infof("Start HTTP server %s ", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		utils.Logger.Errorf("HTTP server: %s", err.Error())
	}
}
