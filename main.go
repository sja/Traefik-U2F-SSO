package main

import (
	"embed"
	"fmt"
	. "github.com/Tedyst/Traefik-U2F-SSO/config"
	"github.com/Tedyst/Traefik-U2F-SSO/internal"
	"github.com/Tedyst/Traefik-U2F-SSO/storage"
	"github.com/Tedyst/Traefik-U2F-SSO/web"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

//go:embed static/*.html
var statics embed.FS

func main() {
	if err := InitConfig(); err != nil {
		panic(fmt.Errorf("could not init config. %w", err))
	}

	logger, err := internal.InitLogger()
	if err != nil {
		log.Fatalf("could not init logger. %v", err)
	}

	s, err := storage.InitStorage(logger)
	if err != nil {
		log.Fatalf("could not init storage. %v", err)
	}
	defer s.CloseDb()

	mux := http.NewServeMux()
	handler := web.NewHandler(logger, statics, s)
	handler.Register(mux)

	logger.Info("Started on :", viper.GetString(ConfPort))

	if err := http.ListenAndServe(":"+viper.GetString(ConfPort), internal.RequestLogger(logger, mux)); err != nil {
		logger.Fatalf("Error in ListenAndServe: %s", err)
	}
}
