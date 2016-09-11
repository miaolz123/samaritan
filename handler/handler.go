package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/cors"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/iris-contrib/middleware/logger"
	"github.com/iris-contrib/middleware/recovery"
	"github.com/kataras/iris"
)

var (
	signKey = []byte("XXXXXXXXXXXXXXXX") // JWT sign key
)

// Server ...
var Server = iris.New()
var jwtmid = jwtmiddleware.New(jwtmiddleware.Config{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func init() {
	Server.Use(cors.New(cors.Options{AllowedHeaders: []string{
		"Origin",
		"X-Requested-With",
		"Content-Type",
		"Accept",
		"Authorization"}}))
	Server.Use(logger.New(iris.Logger))
	Server.Use(recovery.New(iris.Logger))
	Server.Post("/login", userLogin)
	Server.API("/user", userHandler{}, jwtmid.Serve)
	Server.API("/exchange", exchangeHandler{}, jwtmid.Serve)
	Server.API("/strategy", strategyHandler{}, jwtmid.Serve)
	Server.API("/trader", traderHandler{}, jwtmid.Serve)
	Server.Post("/run", jwtmid.Serve, traderRun)
	Server.Post("/stop", jwtmid.Serve, traderStop)
}
