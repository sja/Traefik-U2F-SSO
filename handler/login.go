package handler

import (
	"encoding/json"
	"github.com/Tedyst/Traefik-U2F-SSO/models"
	"github.com/koesie10/webauthn/webauthn"
	"net/http"
)

func (h *Handler) LoginStart(w http.ResponseWriter, r *http.Request) {
	u := &models.User{
		Name: r.URL.Query().Get("name"),
	}
	logger := h.logger.With("User", u.Name)

	sess, err := h.sessionsStore.Get(r, "auth_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Errorw("Error getting a session", "Session", sess.ID)
		return
	}

	logger = logger.With("Session", sess.ID)
	logger.Debugw("Started logging in")

	h.webauth.StartLogin(r, w, u, webauthn.WrapMap(sess.Values))
	err = sess.Save(r, w)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) LoginFinish(w http.ResponseWriter, r *http.Request) {
	u := &models.User{
		Name: r.URL.Query().Get("name"),
	}
	logger := h.logger.With("User", u.Name)

	sess, err := h.sessionsStore.Get(r, "auth_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Errorw("Error getting a session", "Session", sess.ID)
		return
	}

	logger = logger.With("Session", sess.ID)
	logger.Debug("Finishing logging in")

	authenticator := h.webauth.FinishLogin(r, w, u, webauthn.WrapMap(sess.Values))
	if authenticator == nil {
		logger.Debug("Did not finish logging in")
		return
	}

	_, ok := authenticator.(*models.Authenticator)
	if !ok {
		logger.Error("Error casting authenticator")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug("Logged in")

	payload, _ := json.Marshal(u)
	sess.Values["logged"] = true
	sess.Values["username"] = u.Name
	err = sess.Save(r, w)
	if err != nil {
		logger.Error(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(payload)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
