package api

import "go.uber.org/fx"

var Constructors = fx.Options(
	fx.Provide(NewSubtaskServer),
	fx.Provide(NewRequestHandler),
	fx.Provide(NewSubtaskRoutes),
)
