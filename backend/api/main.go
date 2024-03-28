package main

import (
	"cognix.ch/api/v2/api/handler"
	"cognix.ch/api/v2/core/oauth"
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/server"
	"cognix.ch/api/v2/core/utils"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
	"net/http"
)

// @title Cognix API
// @version 1.0
// @description This is Cognix Golang API Documentation

// @contact.name API Support
// @contact.url
// @contact.email

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @BasePath /
// @query.collection.format multi

func main() {

	cfg, err := ReadConfig()
	if err != nil {
		zap.S().Errorf("read log %s:", err.Error())
		return
	}
	utils.InitLogger(cfg.Debug)
	db, err := repository.NewDatabase(cfg.DB)
	if err != nil {
		otelzap.S().Errorf("Init database connection: %s", err.Error())
		return
	}
	oauthProxy := oauth.NewGoogleProvider(cfg.OAuth, cfg.RedirectURL)
	jwtService := security.NewJWTService(cfg.JWTSecret, cfg.JWTExpiredTime)
	// repositories
	connectorRepo := repository.NewConnectorRepository(db)
	credentialRepo := repository.NewCredentialRepository(db)
	userRepo := repository.NewUserRepository(db)
	// handlers
	authMiddleware := server.NewAuthMiddleware(jwtService)
	authHandler := handler.NewAuthHandler(oauthProxy, jwtService, userRepo)
	connectorHandler := handler.NewCollectorHandler(connectorRepo, credentialRepo)
	credentialHandler := handler.NewCredentialHandler(credentialRepo)
	swaggerHandler := handler.NewSwaggerHandler()
	router := NewRouter()

	connectorHandler.Mount(router, authMiddleware.RequireAuth)
	credentialHandler.Mount(router, authMiddleware.RequireAuth)
	authHandler.Mount(router)
	swaggerHandler.Mount(router)

	RunServer(cfg, router)
}

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.Use(otelgin.Middleware("service-name"))
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
	otelzap.S().Infof("Start HTTP server %s ", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		otelzap.S().Errorf("HTTP server: %s", err.Error())
	}
}
