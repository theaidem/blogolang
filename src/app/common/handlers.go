package common

import (
	//"fmt"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
)

func GetFrontHandler(user sessionauth.User, r render.Render) {
	arg_map := map[string]interface{}{"authuser": user}
	r.HTML(200, "home", arg_map)
}

func GetAboutHandler(user sessionauth.User, r render.Render) {
	arg_map := map[string]interface{}{"authuser": user}
	r.HTML(200, "about", arg_map)
}
