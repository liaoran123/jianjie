package routers

import (
	"encoding/json"
	"jianjie/pubgo"
	"jianjie/xbdb"
	"net/http"
)

var (
	pubtbmethod map[string]func(w http.ResponseWriter, req *http.Request) //查询添加修改删除操作处理。
)

func Pubtb(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/pubtb/" + req.Method)

	//req.pubtbmethod
	if pubtbmethod == nil {
		pubtbmethod = make(map[string]func(w http.ResponseWriter, req *http.Request), 4)
		pubtbmethod["POST"] = pubtbpost     //添加
		pubtbmethod["GET"] = pubtbget       //查询
		pubtbmethod["DELETE"] = pubtbdelete //删除
		pubtbmethod["PUT"] = pubtbput       //pubtbput       //修改
	}
	if f, ok := pubtbmethod[req.Method]; ok {
		f(w, req)
	}
}
func pubtbpost(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var r xbdb.ReInfo
	params := postparas(req)
	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	r = Table[params["tbname"]].Ins(params)
	json.NewEncoder(w).Encode(r)
}
func pubtbget(w http.ResponseWriter, req *http.Request) {
	params := getparas(req)
	key := Table[params["tbname"]].Ifo.FieldChByte("id", params["id"])
	tbd := Table[params["tbname"]].Select.OneRecord(key)
	json := Table[params["tbname"]].DataToJson(tbd)
	w.Write(json.Bytes())
	json.Reset()
	xbdb.Bufpool.Put(json)
}
func pubtbdelete(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var r xbdb.ReInfo
	params := postparas(req)
	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	r = Table[params["tbname"]].Del(params["id"])
	json.NewEncoder(w).Encode(r)
}
func pubtbput(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var r xbdb.ReInfo
	params := postparas(req)
	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	r = Table[params["tbname"]].Upd(params)
	json.NewEncoder(w).Encode(r)
}
