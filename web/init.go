package web

import (
	"embed"
	. "github.com/Tedyst/Traefik-U2F-SSO/config"
	"github.com/Tedyst/Traefik-U2F-SSO/internal"
	"github.com/Tedyst/Traefik-U2F-SSO/storage"
	"github.com/Tedyst/sqlitestore"
	"github.com/koesie10/webauthn/webauthn"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	config        Config
	logger        *zap.SugaredLogger
	statics       embed.FS
	sessionsStore *sqlitestore.SqliteStore
	webauth       *webauthn.WebAuthn
}

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

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/webauthn/login/start", h.LoginStart)
	mux.HandleFunc("/webauthn/login/finish", h.LoginFinish)
	mux.HandleFunc("/webauthn/registration/start", h.RegistrationStart)
	mux.HandleFunc("/webauthn/registration/finish", h.RegistrationFinish)
	mux.HandleFunc("/verify", h.Verify)
}

func (h *Handler) render(w http.ResponseWriter, file string) {
	content, err := h.statics.ReadFile(file)
	if err != nil {
		h.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if _, err := w.Write(content); err != nil {
		h.logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
