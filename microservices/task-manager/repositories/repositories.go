package repositories

import "go.uber.org/fx"

var Constructors = fx.Options(
	fx.Provide(NewTaskRepository),
	fx.Provide(NewStorageRepository),
	fx.Provide(NewMessageBrokerRepository),
)
