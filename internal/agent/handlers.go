package agent

import (
	"net/http"
	"net/url"
)

func (a *Agent) solvingExpression(w http.ResponseWriter, r *http.Request) {
	val := url.Values{}
	val.Add("testVal", r.FormValue("expression"))
	resp, err := http.PostForm(a.AddrMainServer+"postResult", val)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func (a *Agent) redirectToMainServer(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.AddrMainServer, http.StatusSeeOther)
}