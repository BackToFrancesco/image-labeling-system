package datasources

import (
	"go.uber.org/fx"
)

var Constructors = fx.Options(
	fx.Provide(NewDatabase),
	fx.Provide(NewStorage),
	fx.Provide(NewMessageBroker),
)
