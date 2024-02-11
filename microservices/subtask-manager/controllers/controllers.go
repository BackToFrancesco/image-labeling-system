package controllers

import "go.uber.org/fx"

var Constructors = fx.Options(
	fx.Provide(NewTaskController),
)
