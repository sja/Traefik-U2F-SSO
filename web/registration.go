package web

import (
	"github.com/Tedyst/Traefik-U2F-SSO/models"
	"github.com/koesie10/webauthn/webauthn"
	"net/http"
)

func (h *Handler) RegistrationFinish(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	if !h.config.Registration.Allowed {
		http.Error(w, "Registration not allowed in config", http.StatusForbidden)
		logger.Debug("Registration attempt denied since the registration is disabled in config")
		return
	}
	u := &models.User{
		Name: r.URL.Query().Get("name"),
	}

	// TODO sja: if domain mismatches, we have an error here, log that!
	sess, err := h.sessionsStore.Get(r, "auth_session")
	logger = logger.With("Session", sess.ID, "User", u.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Errorw("Error getting a session")
		return
	}

	logger.Debugw("Finishing registration")
	h.webauth.FinishRegistration(r, w, u, webauthn.WrapMap(sess.Values))
}

func (h *Handler) RegistrationStart(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	if !h.config.Registration.Allowed {
		http.Error(w, "Registration not allowed in config", http.StatusForbidden)
		logger.Debug("Registration attempt denied since the registration is disabled in config")
		return
	}
	if h.config.Registration.Token != r.URL.Query().Get("token") {
		http.Error(w, "Wrong token", http.StatusForbidden)
		logger.Debug("Registration attempt denied since the token is wrong")
		return
	}
	u := &models.User{
		Name: r.URL.Query().Get("name"),
	}

	sess, err := h.sessionsStore.Get(r, "auth_session")
	logger = logger.With("Session", sess.ID, "User", u.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Errorw("Error getting a session")
		return
	}

	logger.Debugw("Started registration")
	h.webauth.StartRegistration(r, w, u, webauthn.WrapMap(sess.Values))

	if err := sess.Save(r, w); err != nil {
		h.logger.Error("error persisting registration: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
