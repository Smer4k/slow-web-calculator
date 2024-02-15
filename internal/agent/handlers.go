package agent

import (
	"net/http"
)

func (a *Agent) solvingExpression(w http.ResponseWriter, r *http.Request) {
	
}

func (a *Agent) redirectToMainServer(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.AddrMainServer, http.StatusSeeOther)
}