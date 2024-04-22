package main

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	utils.InitLogger(true)
	app := fx.New(Module, fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: log}
	}),
	)

	app.Run()
}
