package routers

import (
	"fmt"
	"net/http"
)

func Redir(w http.ResponseWriter, req *http.Request) {
	//http.Redirect(w, req, "/err/?url="+req.RequestURI, http.StatusFound)
	//http.Redirect(w, req, "/wenhuashuo/", http.StatusFound)
	http.Redirect(w, req, "/", http.StatusFound)
}

//更新服务器
func Updb(w http.ResponseWriter, req *http.Request) {
	Table["j"].Select.ForDbase(updbfun)
}
func updbfun(k, v []byte) bool {
	fmt.Println(string(k), string(v))
	/*
		nk := bytes.Replace(k, []byte("--"), []byte("-"), -1)

		if bytes.Contains(k, []byte("j")) {
			nk = bytes.Replace(nk, []byte("-0"), []byte("- 0"), -1)
		}
		if bytes.Contains(k, []byte("u")) {
			nk = bytes.Replace(nk, []byte(".c"), []byte(". c"), -1)
		}

		nk = bytes.Replace(nk, []byte(".."), []byte("."), -1)

		nv := bytes.Replace(v, []byte("--"), []byte("-"), -1)
		if bytes.Contains(k, []byte("j")) {
			nv = bytes.Replace(nv, []byte("-0"), []byte("- 0"), -1)
		}
		if bytes.Contains(k, []byte("u")) {
			nv = bytes.Replace(nv, []byte(".c"), []byte(". c"), -1)
		}
		nv = bytes.Replace(nv, []byte(".."), []byte("."), -1)
		err := Xb.Delete(k, nil)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		err = Xb.Put(nk, nv, nil)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}*/
	return true
}
