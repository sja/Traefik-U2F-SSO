package handler

import (
	"embed"
	. "github.com/Tedyst/Traefik-U2F-SSO/config"
	"github.com/Tedyst/Traefik-U2F-SSO/internal"
	"github.com/Tedyst/Traefik-U2F-SSO/storage"
	"go.uber.org/zap"
)

func NewHandler(config Config, logger *zap.SugaredLogger, statics embed.FS, storage *storage.Storage) *Handler {
	webauthnService, err := internal.InitWebauthn(config, storage)
	if err != nil {
		logger.Fatalw("could not init webauthn", err)
		return nil
	}
	h := Handler{
		config:        config,
		logger:        logger,
		statics:       statics,
		sessionsStore: storage.GetSessionsStore(),
		webauth:       webauthnService,
	}
	return &h
}
