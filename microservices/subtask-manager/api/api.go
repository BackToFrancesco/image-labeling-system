package api

import "go.uber.org/fx"

var Constructors = fx.Options(
	fx.Provide(NewTaskServer),
	fx.Provide(NewRequestHandler),
	fx.Provide(NewTaskRoutes),
)
