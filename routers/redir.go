package routers

import (
	"net/http"
)

func Redir(w http.ResponseWriter, req *http.Request) {
	//http.Redirect(w, req, "/err/?url="+req.RequestURI, http.StatusFound)
	//http.Redirect(w, req, "/wenhuashuo/", http.StatusFound)
	http.Redirect(w, req, "/", http.StatusFound)
}
