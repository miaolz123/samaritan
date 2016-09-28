package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/cors"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/iris-contrib/middleware/logger"
	"github.com/iris-contrib/middleware/recovery"
	"github.com/kataras/iris"
	"github.com/miaolz123/samaritan/config"
)

var (
	signKey = []byte("XXXXXXXXXXXXXXXX") // JWT sign key
	jwtmid  = jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(signKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
)

func web(c *iris.Context) {
	filename := c.Param("filename")
	if filename == "" {
		filename = "index.html"
	}
	if err := c.ServeFile("web/dist/"+filename, true); err != nil {
		c.Error("Not Found", iris.StatusNotFound)
	}
}

// Run ...
func Run() {
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
	server.Post("/login", userLogin)
	server.Post("/token", jwtmid.Serve, token)
	server.API("/user", userHandler{}, jwtmid.Serve)
	server.API("/exchange", exchangeHandler{}, jwtmid.Serve)
	server.API("/strategy", strategyHandler{}, jwtmid.Serve)
	server.API("/trader", traderHandler{}, jwtmid.Serve)
	server.Post("/run", jwtmid.Serve, traderRun)
	server.Post("/stop", jwtmid.Serve, traderStop)
	server.Post("/logs", jwtmid.Serve, logs)
	server.Delete("/logs", jwtmid.Serve, logsDelete)
	server.Get("/", web)
	server.Get("/dist/:filename", web)
	server.Listen(":" + config.String("port"))
}
