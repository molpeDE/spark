package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/molpeDE/spark/cmd/testapp/app"
	"github.com/molpeDE/spark/cmd/testapp/cfg"
	"github.com/molpeDE/spark/cmd/testapp/frontend"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	configPath := flag.String("config", "./config.ini", "config file path")
	flag.Parse()

	err := cfg.Parse(*configPath)

	if err != nil {
		log.Fatalln(err)
	}

	server := echo.New()
	server.HideBanner = true

	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "[${status}] ${method} ${uri} (${latency_human})\n"}))
	server.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
		log.Printf("Caught Error (%s): %s\n%s", c.Path(), err, string(stack))
		return c.String(http.StatusInternalServerError, "internal server error")
	}}))

	server.GET("/*", echo.WrapHandler(frontend.SPA), middleware.Gzip()) // send out assets with gzip compression (if available)

	backend := server.Group("/rpc")
	/** Example JWT Validation

	backend.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    &globals.JWTKey.PublicKey,
		ContextKey:    "jwt",
		SigningMethod: "ES256",
		TokenLookup:   "cookie:session",
		Skipper: func(c echo.Context) bool {
			return slices.Contains(unauthenticatedRoutes, c.Path())
		},
	}))

	*/
	app.Instance.Attach(backend)

	server.Logger.Fatal(server.Start(fmt.Sprintf("%s:%d", cfg.Get().Server.Host, cfg.Get().Server.Port)))

}
