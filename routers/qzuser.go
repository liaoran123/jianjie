package routers

import (
	"encoding/json"
	"jianjie/pubgo"
	"jianjie/xbdb"
	"net/http"
)

var (
	qzusermethod map[string]func(w http.ResponseWriter, req *http.Request) //查询添加修改删除操作处理。

)

//群组
func Qzuser(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/qzuser/" + req.Method)

	//req.qzusermethod
	if qzusermethod == nil {
		qzusermethod = make(map[string]func(w http.ResponseWriter, req *http.Request), 4)
		qzusermethod["POST"] = qzuserpost     //添加
		qzusermethod["GET"] = qzuserget       //查询
		qzusermethod["DELETE"] = qzuserdelete //删除
		qzusermethod["PUT"] = qzuserput       //qzuserput       //修改
	}
	if f, ok := qzusermethod[req.Method]; ok {
		f(w, req)
	}
}
func qzuserpost(w http.ResponseWriter, req *http.Request) {
	var r xbdb.ReInfo
	params := postparas(req) //  postparas(req)
	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	tbname := params["tbname"]
	buserid := Table[tbname].Ifo.FieldChByte("userid", params["userid"])
	tdb := Table[tbname].Select.WhereIdx([]byte("userid"), buserid, true, 0, -1)
	if len(tdb.Rd) < 3 { //最多能建立3个群组
		r = PPOST(params)
	} else {
		r.Info = "最多能建立3个群组"
	}
	tdb.Release()
	json.NewEncoder(w).Encode(r)
}

//打开群组文章
func qzuserget(w http.ResponseWriter, req *http.Request) {
	params := getparas(req) //  postparas(req)
	tbname := params["tbname"]
	buserid := Table[tbname].Ifo.FieldChByte("userid", params["userid"])
	tdb := Table[tbname].Select.WhereIdx([]byte("type"), buserid, true, 0, -1)
	r := Table[tbname].DataToJson(tdb)
	w.Write(r.Bytes())
	r.Reset()
	xbdb.Bufpool.Put(r)
}
func qzuserdelete(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	params := postparas(req)
	Table["u"].Del(params["id"])
}
func qzuserput(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	params := postparas(req)
	Table["u"].Upd(params)
}
