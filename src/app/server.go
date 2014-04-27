package app

import (
	"app/common"
	"app/config"
	"app/user"
	"github.com/go-martini/martini"
)

func GetApp(cfg *config.Config) *martini.ClassicMartini {
	server := martini.Classic()

	sess := common.InitDB(cfg)
	server.Map(cfg)
	server.Map(sess)

	server.Get("/", common.FrontHandler)
	server.Get("/about", common.AboutHandler)

	server.Get("/signup", user.SignUpHandler)
	server.Get("/signin", user.SignInHandler)
	server.Get("/signout", user.SignOutHandler)

	return server
}
