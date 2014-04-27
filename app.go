package main

import (
	"app"
	"app/config"
)

var (
	cfg = config.Config{
		SiteName:      "Blogolang",
		SessionSecret: "ultra_mega_secret",
		DB: config.DB{
			Address:  "localhost:28015",
			Database: "blogolang",
			Tables:   []string{"users", "posts", "comments"},
		},
	}
)

func main() {
	app := app.GetApp(&cfg)
	app.Run()
}
