package handler

import (
	"net/http"
	"net/url"
)

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	sess, err := h.sessionsStore.Get(r, "auth_session")
	logger := h.logger.With("Session", sess.ID)
	if err != nil {
		http.Redirect(w, r, h.config.URL, http.StatusSeeOther)
		logger.Debugw("Error getting the session")
		return
	}

	if sess.Values["logged"] == true {
		if u, ok := sess.Values["username"].(string); ok {
			w.Header().Set("X-Authenticated-User", u)
			logger = logger.With("username", u)
		}
		logger.Debug("User is logged in")
		return
	}

	redirectToURL := h.config.URL
	if redirectToParam := r.URL.Query().Get("rd"); redirectToParam != "" {
		redirectToURL += "?rd=" + redirectToParam
	} else {
		u, err := url.Parse(r.Header.Get("X-Forwarded-Uri"))
		if err != nil {
			logger.Warnf("invalid X-Forwarded-Uri: %q", r.Header.Get("X-Forwarded-Uri"))
		} else {
			redirectToURL += "?rd=" + u.String()
		}
	}
	logger.Debugw("redirecting unauthenticated user", "location", redirectToURL)
	http.Redirect(w, r, redirectToURL, http.StatusSeeOther)
}
