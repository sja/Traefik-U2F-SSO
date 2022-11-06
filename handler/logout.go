package handler

import (
	"github.com/Tedyst/Traefik-U2F-SSO/models"
	"net/http"
	"time"
)

// Logout will delete the auth_session cookie and redirects to a configured url or to referer.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	u := &models.User{
		Name: r.URL.Query().Get("name"),
	}
	logger := h.logger.With("User", u.Name)

	_, err := h.sessionsStore.Get(r, h.config.Session.CookieName)
	if err != nil {
		logger.Infof("session not found, either expired or invalid")
	}

	deleteCookie(w, h.config.Session.CookieName)

	redirectUrl := h.config.URL
	if referrer := r.Header.Get("referer"); len(referrer) > 0 {
		redirectUrl = referrer
	}

	logger.Infow("redirect after logout", "location", redirectUrl)
	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func deleteCookie(w http.ResponseWriter, name string) {
	cookie := http.Cookie{
		Name:    name,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		Secure:  true,
	}
	http.SetCookie(w, &cookie)
}
