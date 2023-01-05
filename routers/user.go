package routers

import (
	"encoding/json"
	"jianjie/pubgo"
	"jianjie/xbdb"
	"net/http"
	"strconv"
	"strings"
)

var (
	usermethod map[string]func(w http.ResponseWriter, req *http.Request) //查询添加修改删除操作处理。

)

func User(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/user/" + req.Method)

	//req.usermethod
	if usermethod == nil {
		usermethod = make(map[string]func(w http.ResponseWriter, req *http.Request), 4)
		usermethod["POST"] = userpost     //添加
		usermethod["GET"] = userget       //查询
		usermethod["DELETE"] = userdelete //删除
		usermethod["PUT"] = userput       //userput       //修改
	}
	if f, ok := usermethod[req.Method]; ok {
		f(w, req)
	}
}
func userpost(w http.ResponseWriter, req *http.Request) {
	mu.Lock() //leveldb仅支持单进程数据操作。
	defer mu.Unlock()
	var r xbdb.ReInfo
	params := postparas(req)
	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	key := Table["u"].Select.GetIdxPrefix([]byte("email"), []byte(params["email"]))
	_, ok := Table["u"].Select.IterPrefixMove(key, true)
	if ok {
		r.Info = "邮箱已存在。"
		json.NewEncoder(w).Encode(r)
		return
	}
	params["psw"] = Md5(params["psw"])
	r = Table["u"].Ins(params)
	json.NewEncoder(w).Encode(r)
}
func userget(w http.ResponseWriter, req *http.Request) {
	var r xbdb.ReInfo
	params := getparas(req) //  postparas(req)
	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	params["psw"] = Md5(params["psw"])
	tbd := Table["u"].Select.WhereIdx([]byte("email"), []byte(params["email"]), true, 0, -1)
	if tbd == nil {
		r.Info = "密码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	rd := strings.Split(string(tbd.Rd[0]), xbdb.Split)
	psw := rd[2]
	id := xbdb.BytesToInt([]byte(rd[0]))
	sid := strconv.Itoa(id)
	if psw == params["psw"] {
		r.Info = "id:" + sid + ",法号:" + rd[3] + "\n登陆成功!"
		r.Succ = true
	} else {
		r.Info = "密码不对!"
		r.Succ = false
	}
	json.NewEncoder(w).Encode(r)

}
func userdelete(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	params := postparas(req)
	Table["u"].Del(params["id"])
}
func userput(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	params := postparas(req)
	Table["u"].Upd(params)
}
