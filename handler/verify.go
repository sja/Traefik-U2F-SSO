package handler

import (
	"net/http"
)

func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	sess, err := h.sessionsStore.Get(r, h.config.Session.CookieName)
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

	redirectToURL := h.config.URL + "?rd="
	if redirectToParam := r.URL.Query().Get("rd"); redirectToParam != "" {
		redirectToURL += redirectToParam
	} else {
		u := h.getUrlFromForwardedHeaders(r)
		redirectToURL += u.String()
	}
	/*for name, values := range r.Header {
		logger = logger.With(name, strings.Join(values, ", "))
	}*/
	logger.Debugw("redirecting unauthenticated user", "location", redirectToURL)
	http.Redirect(w, r, redirectToURL, http.StatusSeeOther)
}
