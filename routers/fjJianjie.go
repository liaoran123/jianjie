package routers

import (
	"encoding/json"
	"jianjie/pubgo"
	"jianjie/xbdb"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	FjJianjiemethod map[string]func(w http.ResponseWriter, req *http.Request) //查询添加修改删除操作处理。

)

func FjJianjie(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/FjJianjie/" + req.Method)

	//req.FjJianjiemethod
	if FjJianjiemethod == nil {
		FjJianjiemethod = make(map[string]func(w http.ResponseWriter, req *http.Request), 4)
		FjJianjiemethod["POST"] = FjJianjiepost     //添加
		FjJianjiemethod["GET"] = FjJianjieget       //查询
		FjJianjiemethod["DELETE"] = FjJianjiedelete //删除
		FjJianjiemethod["PUT"] = FjJianjieput       //FjJianjieput       //修改
	}
	if f, ok := FjJianjiemethod[req.Method]; ok {
		f(w, req)
	}
}
func FjJianjiepost(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	//params := postparas(req)
	//Table["jianjie"].Ins(params)

	var r xbdb.ReInfo
	params := postparas(req)
	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	//检查用户是否存在
	userid := params["userid"]
	if userid == "" {
		r.Info = "该用户不存在。"
		json.NewEncoder(w).Encode(r)
		return
	}
	_, err := strconv.Atoi(userid) //userid为字符串的话，FieldChByte会得到1
	if err != nil {
		r.Info = "该用户不存在。"
		json.NewEncoder(w).Encode(r)
		return
	}
	buserid := Table["u"].Ifo.FieldChByte("id", userid)
	if Table["u"].Select.GetPKValue(buserid) == nil {
		r.Info = "该用户不存在。"
		json.NewEncoder(w).Encode(r)
		return
	}
	params["sj"] = strings.Split(time.Now().String(), ".")[0]
	r = Table["j"].Ins(params)
	json.NewEncoder(w).Encode(r)

}
func FjJianjieget(w http.ResponseWriter, req *http.Request) {
	params := getparas(req)
	key := Table["j"].Ifo.FieldChByte("id", params["id"])
	tbd := Table["j"].Select.OneRecord(key)
	if tbd == nil {
		return
	}
	json := Table["j"].DataToJson(tbd)
	w.Write(json.Bytes())
	json.Reset()
	xbdb.Bufpool.Put(json)
}
func FjJianjiedelete(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	params := postparas(req)
	Table["j"].Del(params["id"])
}
func FjJianjieput(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	params := postparas(req)
	Table["j"].Upd(params)
}
