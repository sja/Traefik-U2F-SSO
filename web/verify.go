package web

import (
	"net/http"
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
		logger.Debugw("User is logged in")
		return

	}
	logger.Debugw("User is not logged in")
	newURL := r.URL.Query().Get("rd")
	if newURL != "" {
		http.Redirect(w, r, h.config.URL+"?rd="+newURL, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, h.config.URL, http.StatusSeeOther)
}
