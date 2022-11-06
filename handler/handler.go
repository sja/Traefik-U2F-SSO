package handler

import (
	"embed"
	. "github.com/Tedyst/Traefik-U2F-SSO/config"
	"github.com/Tedyst/sqlitestore"
	"github.com/koesie10/webauthn/webauthn"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
)

type Handler struct {
	config        Config
	logger        *zap.SugaredLogger
	statics       embed.FS
	sessionsStore *sqlitestore.SqliteStore
	webauth       *webauthn.WebAuthn
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/webauthn/login/start", h.LoginStart)
	mux.HandleFunc("/webauthn/login/finish", h.LoginFinish)
	mux.HandleFunc("/webauthn/registration/start", h.RegistrationStart)
	mux.HandleFunc("/webauthn/registration/finish", h.RegistrationFinish)
	mux.HandleFunc("/verify", h.Verify)
	mux.HandleFunc("/logout", h.Logout)

	mux.HandleFunc("/.well-known/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
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

func (h *Handler) getUrlFromForwardedHeaders(r *http.Request) url.URL {
	u := url.URL{
		Scheme: r.Header.Get("X-Forwarded-Proto"),
		Host:   r.Header.Get("X-Forwarded-Host"),
		Path:   r.Header.Get("X-Forwarded-Uri"),
	}
	if !strings.Contains(u.Host, h.config.Session.Domain) {
		h.logger.Warnf("x-forwarded header contain domain which may not belong to us: %q", u.Host)
	}
	return u
}
