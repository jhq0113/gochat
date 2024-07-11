package actions

import (
	"github.com/jhq0113/gochat/lib/constants"
	"github.com/jhq0113/gochat/lib/protocol"
)

func LoadRouter() *protocol.Router {
	router := protocol.NewRouter()
	router.Add(constants.Login, Login)
	return router
}
