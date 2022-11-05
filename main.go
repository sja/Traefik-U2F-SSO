package main

import (
	"embed"
	. "github.com/Tedyst/Traefik-U2F-SSO/config"
	"github.com/Tedyst/Traefik-U2F-SSO/handler"
	"github.com/Tedyst/Traefik-U2F-SSO/internal"
	"github.com/Tedyst/Traefik-U2F-SSO/storage"
	"log"
	"net/http"
)

//go:embed static/*.html
var statics embed.FS

func main() {
	config := InitConfig()
	if err := config.Validate(); err != nil {
		log.Fatal("config invalid", err)
	}

	logger, err := internal.InitLogger(config)
	if err != nil {
		log.Fatal("could not init logger", err)
	}

	s, err := storage.InitStorage(config, logger)
	if err != nil {
		log.Fatal("could not init storage", err)
	}
	defer s.CloseDb()
	defer logger.Sync()

	mux := http.NewServeMux()
	h := handler.NewHandler(config, logger, statics, s)
	h.Register(mux)

	logger.Infof("Started listening on %q", config.Serve)

	if err := http.ListenAndServe(config.Serve, internal.RequestLogger(logger, mux)); err != nil {
		logger.Fatal("Error in ListenAndServe", err)
	}
}
