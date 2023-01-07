package routers

import (
	"encoding/json"
	"jianjie/pubgo"
	"jianjie/xbdb"
	"net/http"
)

var (
	pubtbmethod map[string]func(w http.ResponseWriter, req *http.Request) //查询添加修改删除操作处理。
	posts       map[string]func(params map[string]string) (r xbdb.ReInfo) //查询除外的执行都是类同的。
)

func Pubtb(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/pubtb/" + req.Method)

	//req.pubtbmethod
	if pubtbmethod == nil {
		pubtbmethod = make(map[string]func(w http.ResponseWriter, req *http.Request), 4)
		pubtbmethod["POST"] = pubtbposts   //添加
		pubtbmethod["GET"] = pubtbget      //查询
		pubtbmethod["DELETE"] = pubtbposts //删除
		pubtbmethod["PUT"] = pubtbposts    //pubtbput       //修改
	}
	if f, ok := pubtbmethod[req.Method]; ok {
		f(w, req)
	}
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
func pubtbposts(w http.ResponseWriter, req *http.Request) {
	if posts == nil {
		posts = make(map[string]func(params map[string]string) (r xbdb.ReInfo), 3)
		posts["POST"] = PPOST     //添加
		posts["DELETE"] = PDELETE //删除
		posts["PUT"] = PPUT       //pubtbput       //修改
	}
	mu.Lock()
	defer mu.Unlock()
	var r xbdb.ReInfo
	params := postparas(req)

	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		if params["tbname"] != "d" && req.Method == "POST" { //例外。点赞表添加不需要验证码。
			return
		}
	}
	r = posts[req.Method](params)
	json.NewEncoder(w).Encode(r)
}
func PPOST(params map[string]string) (r xbdb.ReInfo) {
	r = Table[params["tbname"]].Ins(params)
	return
}
func PDELETE(params map[string]string) (r xbdb.ReInfo) {
	r = Table[params["tbname"]].Del(params["id"])
	return
}
func PPUT(params map[string]string) (r xbdb.ReInfo) {
	r = Table[params["tbname"]].Upd(params)
	return
}

/*
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
	if !VerifyTime(params["md5code"], params["capid"]) {
		r.Info = "错误，请重试"
		json.NewEncoder(w).Encode(r)
		return
	}
	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	r = Table[params["tbname"]].Upd(params)
	json.NewEncoder(w).Encode(r)
}
*/
