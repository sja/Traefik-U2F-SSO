package main

import (
	"embed"
	"fmt"
	. "github.com/Tedyst/Traefik-U2F-SSO/config"
	"github.com/Tedyst/Traefik-U2F-SSO/internal"
	"github.com/Tedyst/Traefik-U2F-SSO/storage"
	"github.com/Tedyst/Traefik-U2F-SSO/web"
	"log"
	"net/http"
)

//go:embed static/*.html
var statics embed.FS

func main() {
	config := InitConfig()
	if err := config.Validate(); err != nil {
		log.Fatal(fmt.Errorf("config invalid: %w", err))
	}

	logger, err := internal.InitLogger(config)
	if err != nil {
		log.Fatalf("could not init logger. %v", err)
	}

	s, err := storage.InitStorage(config, logger)
	if err != nil {
		log.Fatalf("could not init storage. %v", err)
	}
	defer s.CloseDb()

	mux := http.NewServeMux()
	handler := web.NewHandler(config, logger, statics, s)
	handler.Register(mux)

	logger.Info("Started listening on %v", config.Serve)

	if err := http.ListenAndServe(config.Serve, internal.RequestLogger(logger, mux)); err != nil {
		logger.Fatalf("Error in ListenAndServe: %s", err)
	}
}
