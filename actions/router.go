package actions

import (
	"github.com/jhq0113/gochat/lib/constants"
	"github.com/jhq0113/gochat/lib/router"
)

func LoadRouter() *router.Router {
	route := router.NewRouter()
	route.Add(constants.Login, Login)
	route.Add(constants.Join, Join)
	return route
}
