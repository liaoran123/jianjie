package routers

import (
	"encoding/json"
	"jianjie/pubgo"
	"jianjie/xbdb"
	"net/http"
	"strconv"
	"time"
)

var (
	FjJianjiemethod map[string]func(w http.ResponseWriter, req *http.Request) //查询添加修改删除操作处理。

)

func FjJianjie(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json") ////返回数据格式是json

	pubgo.Tj.Brows("/FjJianjie/" + req.Method)

	//req.FjJianjiemethod
	if FjJianjiemethod == nil {
		FjJianjiemethod = make(map[string]func(w http.ResponseWriter, req *http.Request), 2)
		FjJianjiemethod["POST"] = FjJianjiepost //添加
		FjJianjiemethod["GET"] = FjJianjieget   //查询
		//FjJianjiemethod["DELETE"] = FjJianjiedelete //删除
		//FjJianjiemethod["PUT"] = FjJianjieput       //FjJianjieput       //修改
	}
	if f, ok := FjJianjiemethod[req.Method]; ok {
		f(w, req)
	}
}
func FjJianjiepost(w http.ResponseWriter, req *http.Request) {
	//mu.Lock()
	//defer mu.Unlock()

	var r xbdb.ReInfo
	params := postparas(req)

	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	ptype := params["type"]
	if ptype != "" { //浏览器不支持delete,只能在这里转
		if ptype == "0" {
			FjJianjiedelete(w, req, params)
		} else {
			FjJianjieput(w, req, params)
		}
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
	if params["sj"] == "" {
		params["sj"] = "now()"
	}
	r = Table["j"].Ins(params)
	json.NewEncoder(w).Encode(r)

}
func FjJianjieget(w http.ResponseWriter, req *http.Request) {
	params := getparas(req)
	key := Table["j"].Ifo.FieldChByte("id", params["id"])
	tbd := Table["j"].Select.Record(key, []int{})
	if tbd == nil {
		return
	}
	json := Table["j"].DataToJsonApp(tbd)
	w.Write(json.Bytes())
	json.Reset()
	xbdb.Bufpool.Put(json)
}
func FjJianjiedelete(w http.ResponseWriter, req *http.Request, params map[string]string) {
	//mu.Lock()
	//defer mu.Unlock()
	//params := postparas(req)
	var r xbdb.ReInfo

	//打开要删除的记录，获取时间，超过3天不能删除
	key := Table["j"].Ifo.FieldChByte("id", params["id"])
	tbd := Table["j"].Select.Record(key, []int{})
	tbm := Table["j"].RDtoMap(tbd.Rd[0])
	tbd.Release()
	if tbm["userid"] != params["userid"] { //删除的id和用户id对应才能删除，以防数据错乱和攻击。
		return
	}
	sj, _ := time.ParseInLocation("2006-01-02 15:04:05", tbm["sj"], time.Local)

	if time.Since(sj).Hours() > 72 {
		r.Info = "超过3天不能删除。"
		json.NewEncoder(w).Encode(r)
		return
	}
	r = Table["j"].Del(params["id"])
	json.NewEncoder(w).Encode(r)
}
func FjJianjieput(w http.ResponseWriter, req *http.Request, params map[string]string) {
	//mu.Lock()
	//defer mu.Unlock()
	//var r xbdb.ReInfo
	//params := postparas(req)
	/*
		if !store.Verify(params["capid"], params["code"], true) {
			r.Info = "验证码不正确！"
			json.NewEncoder(w).Encode(r)
			return
		}*/
	r := Table["j"].Upd(params)
	json.NewEncoder(w).Encode(r)
}
