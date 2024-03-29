package routers

import (
	"encoding/json"
	"jianjie/pubgo"
	"jianjie/xbdb"
	"net/http"
)

var (
	usermethod map[string]func(w http.ResponseWriter, req *http.Request) //查询添加修改删除操作处理。

)

func User(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

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
	pubgo.Tj.Brows("/user/" + req.Method)
	//mu.Lock() //leveldb仅支持单进程数据操作。
	//defer mu.Unlock()
	var r xbdb.ReInfo
	params := postparas(req)

	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	bemali := Table["u"].Ifo.TypeChByte("email", params["email"])
	ok := Table["u"].Select.WhereIdxExist([]byte("email"), bemali)
	if ok {
		r.Info = "邮箱已存在。"
		json.NewEncoder(w).Encode(r)
		return
	}
	params["psw"] = Md5(params["psw"])
	if params["sj"] == "" {
		params["sj"] = "now()"
	}
	r = Table["u"].Ins(params)
	json.NewEncoder(w).Encode(r)
}
func userget(w http.ResponseWriter, req *http.Request) {
	pubgo.Tj.Brows("/user/" + req.Method)
	var r xbdb.ReInfo
	params := getparas(req) //  postparas(req)

	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	params["psw"] = Md5(params["psw"])
	bemali := Table["u"].Ifo.TypeChByte("email", params["email"])
	tbd := Table["u"].Select.WhereIdx([]byte("email"), bemali, true, 0, -1, []int{}, false)
	if tbd == nil {
		r.Info = "密码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	rdmap := Table["u"].RDtoMap(tbd.Rd[0])
	tbd.Release()
	if rdmap["pass"] == "0" { //被封用户
		r.Info = "该用户不存在。!"
		r.Succ = false
		json.NewEncoder(w).Encode(r)
	}
	psw := rdmap["psw"] //string(rd[2])
	id := rdmap["id"]   //xbdb.BytesToInt(rd[0])
	fahao := rdmap["fahao"]
	if psw == params["psw"] {
		r.Info = "id:" + id + ",法号:" + fahao + "\n登陆成功!"
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
	//mu.Lock()
	//defer mu.Unlock()
	params := postparas(req)
	Table["u"].Upd(params)
}
func UserDel(w http.ResponseWriter, req *http.Request) {
	params := postparas(req)
	if params["psw"] == "!nmantf123456!" {
		delete(params, "psw")
		r := Table["u"].Upd(params)
		if r.Succ {
			w.Write([]byte("1"))
		} else {
			w.Write([]byte("0"))
		}
	}
}
