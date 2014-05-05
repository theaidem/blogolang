package common

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
)

func GetFrontHandler(user sessionauth.User, r render.Render) {
	if user.IsAuthenticated() {
		r.HTML(200, "home", user)
		return
	}
	r.HTML(200, "home", nil)
}

func GetAboutHandler() string {
	return "About page"
}
