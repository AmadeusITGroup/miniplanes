package www

import (
	"html/template"
	"net/http"
)

// AdminHandler renders the admin view template
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	push(w, "/static/style.css")
	push(w, "/static/navigation_bar.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fullData := map[string]interface{}{
		"NavigationBar": template.HTML(navigationBarHTML),
	}
	// pass from global variables
	if !BasicAuth(w, r, username, password) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin protected."`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
		return
	}

	render(w, r, adminViewTpl, "admin_view", fullData)
}
