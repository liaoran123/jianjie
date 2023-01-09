package routers

import (
	"bytes"
	"fmt"
	"net/http"
)

func Redir(w http.ResponseWriter, req *http.Request) {
	//http.Redirect(w, req, "/err/?url="+req.RequestURI, http.StatusFound)
	//http.Redirect(w, req, "/wenhuashuo/", http.StatusFound)
	http.Redirect(w, req, "/", http.StatusFound)
}
func Updb(w http.ResponseWriter, req *http.Request) {
	Table["j"].Select.ForDbase(updbfun)
}
func updbfun(k, v []byte) bool {
	//fmt.Println(string(k), string(v))
	nk := bytes.Replace(k, []byte("--"), []byte("-"), -1)
	nk = bytes.Replace(nk, []byte("-0"), []byte("- 0"), -1)
	nk = bytes.Replace(nk, []byte(".."), []byte("."), -1)

	nv := bytes.Replace(v, []byte("--"), []byte("-"), -1)
	nv = bytes.Replace(nv, []byte("-0"), []byte("- 0"), -1)
	nv = bytes.Replace(nv, []byte(".."), []byte("."), -1)
	err := Xb.Delete(k, nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	err = Xb.Put(nk, nv, nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	return true
}
