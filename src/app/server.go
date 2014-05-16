package app

import (
	"app/common"
	"app/config"
	"app/user"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
)

func GetApp(cfg *config.Config) *martini.ClassicMartini {
	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))
	store.Options(sessions.Options{MaxAge: 0})

	server := martini.Classic()

	server.Use(render.Renderer(render.Options{Layout: "layout"}))
	server.Use(sessions.Sessions("sid", store))
	server.Use(martini.Static("static", martini.StaticOptions{SkipLogging: true}))

	sess := common.InitDB(cfg)
	server.Map(cfg)
	server.Map(sess)

	server.Use(sessionauth.SessionUser(user.GenerateAnonymousUser))
	sessionauth.RedirectUrl = "/signin"
	sessionauth.RedirectParam = "next"

	server.Get("/", common.GetFrontHandler)
	server.Get("/about", common.GetAboutHandler)

	server.Get("/signup", user.GetSignUpHandler)
	server.Post("/signup", binding.Bind(user.User{}), user.PostSignUpHandler)
	server.Get("/signin", user.GetSignInHandler)
	server.Post("/signin", binding.Bind(user.User{}), user.PostSignInHandler)
	server.Get("/signout", sessionauth.LoginRequired, user.GetSignOutHandler)

	server.Get("/users", binding.Bind(user.User{}), user.GetUserListHandler)
	server.Get("/users/:page", binding.Bind(user.User{}), user.GetUserListHandler)
	server.Get("/user/:id", binding.Bind(user.User{}), user.GetUserProfileHandler)
	server.Get("/user/:id/delete", binding.Bind(user.User{}), user.GetUserDeleteHandler)

	return server
}
