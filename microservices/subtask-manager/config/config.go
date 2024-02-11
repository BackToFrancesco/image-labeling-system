package config

import "go.uber.org/fx"

var Constructors = fx.Options(
	fx.Provide(NewEnv),
)
