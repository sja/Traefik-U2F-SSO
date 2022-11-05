package handler

import (
	"net/http"
)

// Index is the main page where the user logs in or registers
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	sess, err := h.sessionsStore.Get(r, "auth_session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// If the user is logged in, the page shown is static/loggedin.html
	if sess.Values["logged"] == true {
		h.logger.Debugw("User is not logged in",
			"Session", sess.ID,
		)
		newURL := r.URL.Query().Get("rd")
		if newURL != "" {
			http.Redirect(w, r, newURL, http.StatusSeeOther)
			return
		}

		h.render(w, "static/loggedin.html")
		return
	}
	if err := sess.Save(r, w); err != nil {
		h.logger.Errorf("error in peristing: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// If the registration is allowed in config.json, the page shown is static/index.html, that allows registration using a token.
	if h.config.Registration.Allowed {
		h.render(w, "static/index.html")
		return
	}
	// If this page is shown, only logging in using the authenticators is allowed
	h.render(w, "static/justlogin.html")
}
