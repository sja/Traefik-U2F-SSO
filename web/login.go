package web

import (
	"encoding/json"
	"github.com/Tedyst/Traefik-U2F-SSO/models"
	"github.com/koesie10/webauthn/webauthn"
	"log"
	"net/http"
)

func (h *Handler) LoginStart(w http.ResponseWriter, r *http.Request) {
	u := &models.User{
		Name: r.URL.Query().Get("name"),
	}

	sess, err := h.sessionsStore.Get(r, "auth_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorw("Error getting a session",
			"Session", sess.ID,
			"User", u.Name,
		)
		return
	}

	h.logger.Debugw("Started logging in",
		"Session", sess.ID,
		"User", u.Name,
	)
	h.webauth.StartLogin(r, w, u, webauthn.WrapMap(sess.Values))
	err = sess.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) LoginFinish(w http.ResponseWriter, r *http.Request) {
	u := &models.User{
		Name: r.URL.Query().Get("name"),
	}

	sess, err := h.sessionsStore.Get(r, "auth_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Errorw("Error getting a session",
			"Session", sess.ID,
			"User", u.Name,
		)
		return
	}

	h.logger.Debugw("Finishing logging in",
		"Session", sess.ID,
		"User", u.Name,
	)
	authenticator := h.webauth.FinishLogin(r, w, u, webauthn.WrapMap(sess.Values))
	if authenticator == nil {
		h.logger.Debugw("Did not finish logging in",
			"Session", sess.ID,
			"User", u.Name,
		)
		return
	}

	_, ok := authenticator.(*models.Authenticator)
	if !ok {
		h.logger.Debugw("Help",
			"Session", sess.ID,
			"User", u.Name,
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Debugw("Logged in",
		"Session", sess.ID,
		"User", u.Name,
	)

	payload, _ := json.Marshal(u)
	sess.Values["logged"] = true
	err = sess.Save(r, w)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
