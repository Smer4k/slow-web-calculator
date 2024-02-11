package agent

import (
	"fmt"
	"net/http"
)

func (a *Agent) solvingExpression(w http.ResponseWriter, r *http.Request) {

	resp, err := http.PostForm(a.AddrMainServer+"postResult", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	val := r.FormValue("expression")
	fmt.Println(val)
}

func (a *Agent) redirectToMainServer(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.AddrMainServer, http.StatusSeeOther)
}