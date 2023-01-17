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

//所有表公共操作，查询，添加，修改，删除
func Pubtb(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/pubtb/" + req.Method)

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
	tbd := Table[params["tbname"]].Select.Record(key)
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
	var r xbdb.ReInfo
	params := postparas(req)
	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	r = posts[req.Method](params)
	json.NewEncoder(w).Encode(r)
}

//如果filedname的filedvalue不存在，则添加
func PPOST(params map[string]string) (r xbdb.ReInfo) {
	filedname := params["existfiled"]
	tbname := params["tbname"]
	filedvalue := params[filedname]
	if !exist(tbname, filedname, filedvalue) {
		r = Table[tbname].Ins(params)
	}
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

//判断某字段的值是否存在
func exist(tbanme, filedname, filedvalue string) bool {
	bfiledvalue := Table[tbanme].Ifo.TypeChByte(filedname, filedvalue)
	key := Table[tbanme].Select.GetIdxPrefix([]byte(filedname), bfiledvalue)
	_, ok := Table[tbanme].Select.IterPrefixMove(key, true)
	return ok
}
