package main

import (
	"log"
	"net/http"
	"os"

	"trde/redirector"
	"trde/store"

	"github.com/urfave/cli/v2"
)

func main() {
	var err error
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "listen-address",
				Aliases: []string{"l", "listen"},
				Value:   ":8080",
				Usage:   "address to listen on",
				EnvVars: []string{"LISTEN_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "routes-path",
				Aliases: []string{"r", "routes"},
				Value:   "routes.yaml",
				Usage:   "path to the routes file",
				EnvVars: []string{"ROUTES_PATH"},
			},
			&cli.StringFlag{
				Name:    "static-path",
				Aliases: []string{"s", "static"},
				Value:   "/app/portal/dist",
				Usage:   "path to the static files",
				EnvVars: []string{"STATIC_PATH"},
			},
			&cli.BoolFlag{
				Name:    "proxy-enabled",
				Aliases: []string{"p", "proxy"},
				Value:   false,
				Usage:   "enable proxy mode",
				EnvVars: []string{"PROXY_ENABLED"},
			},
			&cli.StringFlag{
				Name:    "proxy-url",
				Aliases: []string{"u", "url"},
				Value:   "http://localhost:5173",
				Usage:   "proxy URL",
				EnvVars: []string{"PROXY_URL"},
			},
		},
		Action: func(c *cli.Context) error {
			listenAddress := c.String("listen-address")
			routesPath := c.String("routes-path")
			staticPath := c.String("static-path")
			proxyEnabled := c.Bool("proxy-enabled")
			proxyURL := c.String("proxy-url")

			storage := store.LocalStore{
				Filepath: routesPath,
			}

			redirector := redirector.Redirector{
				RoutesFilepath: routesPath,
				Store:          &storage,
				StaticFilepath: staticPath,
				ProxyEnabled:   proxyEnabled,
				ProxyURL:       proxyURL,
			}

			err := redirector.ReloadRoutes()
			if err != nil {
				return err
			}

			mux := http.NewServeMux()
			mux.HandleFunc("GET /", redirector.HandleForward)
			mux.HandleFunc("GET /routes.json", redirector.HandleGetRoutes)
			mux.HandleFunc("POST /routes.json", redirector.HandleReloadRoutes)
			mux.HandleFunc("PUT /routes.json", redirector.HandlePutRoute)

			return http.ListenAndServe(listenAddress, mux)
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
