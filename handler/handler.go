package handler

import (
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-ini/ini"
	"github.com/iris-contrib/middleware/cors"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/iris-contrib/middleware/logger"
	"github.com/iris-contrib/middleware/recovery"
	"github.com/kataras/iris"
)

var (
	signKey = []byte("XXXXXXXXXXXXXXXX") // JWT sign key
)

var jwtmid = jwtmiddleware.New(jwtmiddleware.Config{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func init() {
	conf, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Load config.ini error:", err)
	}
	server := iris.New()
	server.Use(cors.New(cors.Options{
		AllowedHeaders: []string{
			"Origin",
			"X-Requested-With",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		AllowedMethods: []string{
			"Get",
			"Post",
			"Delete",
		},
	}))
	server.Use(logger.New())
	server.Use(recovery.New())
	server.Get("/", web)
	server.Post("/login", userLogin)
	server.Post("/token", jwtmid.Serve, token)
	server.API("/user", userHandler{}, jwtmid.Serve)
	server.API("/exchange", exchangeHandler{}, jwtmid.Serve)
	server.API("/strategy", strategyHandler{}, jwtmid.Serve)
	server.API("/trader", traderHandler{}, jwtmid.Serve)
	server.Post("/run", jwtmid.Serve, traderRun)
	server.Post("/stop", jwtmid.Serve, traderStop)
	server.Post("/logs", jwtmid.Serve, logs)
	server.Get("/dist/:filename", web)
	server.Listen(":" + conf.Section("").Key("ServerPort").String())
}

func web(c *iris.Context) {
	filename := c.Param("filename")
	if filename == "" {
		filename = "index.html"
	}
	c.ServeFile("web/dist/"+filename, true)
}
