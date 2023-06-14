package routers

import (
	"encoding/base64"
	"jianjie/pubgo"
	"net/http"
)

func Emali(w http.ResponseWriter, req *http.Request) {
	pubgo.Tj.Brows("/Emali/" + req.Method)
	params := getparas(req) //postparas(req)
	if len(params) == 0 {
		w.Write([]byte("0"))
		return
	}
	bemali := Table["u"].Ifo.TypeChByte("email", params["email"])
	ok := Table["u"].Select.WhereIdxExist([]byte("email"), bemali)
	if ok {
		w.Write([]byte("1"))
		return
	}
	w.Write([]byte("0"))
}

func UpPsw(w http.ResponseWriter, req *http.Request) {
	pubgo.Tj.Brows("/UpPsw/" + req.Method)
	params := getparas(req) //postparas(req)
	if len(params) == 0 {
		w.Write([]byte("0"))
		return
	}
	email, err := base64.URLEncoding.DecodeString(params["email"])
	if err != nil {
		return
	}
	bemali := Table["u"].Ifo.TypeChByte("email", string(email))
	tbd := Table["u"].Select.WhereIdx([]byte("email"), bemali, false, 0, 1, []int{}, false)
	if tbd == nil {
		w.Write([]byte("0"))
		return
	}
	id := Table["u"].RDtoMap(tbd.Rd[0])["id"]
	ps := make(map[string]string)
	ps["id"] = id
	ps["psw"] = params["psw"]

	r := Table["u"].Upd(ps)
	if r.Succ {
		w.Write([]byte("1"))
	} else {
		w.Write([]byte(r.Info))
	}
}
