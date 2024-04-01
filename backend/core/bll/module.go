package bll

import (
	"go.uber.org/fx"
)

var BLLModule = fx.Options(
	fx.Provide(
		NewConnectorBL,
		NewCredentialBL,
		NewAuthBL,
		NewPersonaBL,
	),
)
