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
	pubtbmethod map[string]func(w http.ResponseWriter, req *http.Request) //查询添加修改删除操作处理。
	posts       map[string]func(params map[string]string) (r xbdb.ReInfo) //查询除外的执行都是类同的。
)

//所有表公共操作，查询，添加，修改，删除
func Pubtb(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/pubtb/" + req.Method)

	if pubtbmethod == nil {
		pubtbmethod = make(map[string]func(w http.ResponseWriter, req *http.Request), 2)
		pubtbmethod["POST"] = pubtbposts //添加
		pubtbmethod["GET"] = pubtbget    //查询
		//pubtbmethod["DELETE"] = pubtbposts //删除
		//pubtbmethod["PUT"] = pubtbposts    //pubtbput       //修改
	}
	if f, ok := pubtbmethod[req.Method]; ok {
		f(w, req)
	}
}

func pubtbget(w http.ResponseWriter, req *http.Request) {
	params := getparas(req)
	tbname := params["tbname"]
	key := Table[tbname].Ifo.FieldChByte("id", params["id"])
	tbd := Table[tbname].Select.Record(key)
	json := Table[tbname].DataToJson(tbd)
	w.Write(json.Bytes())
	json.Reset()
	xbdb.Bufpool.Put(json)
}
func pubtbposts(w http.ResponseWriter, req *http.Request) {
	var r xbdb.ReInfo
	params := postparas(req)
	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	r = POST(params)
	json.NewEncoder(w).Encode(r)
}

func POST(params map[string]string) (r xbdb.ReInfo) {
	//类似params解析器。根据不同的参数进行不同的判断和执行相应的函数
	if rdexist(params) {
		r.Info = "记录已存在！"
		return
	}
	if rdCount(params) {
		r.Info = "已超过记录数！"
		return
	}
	if posts == nil {
		posts = make(map[string]func(params map[string]string) (r xbdb.ReInfo), 3)
		posts["POST"] = PPOST     //添加
		posts["DELETE"] = PDELETE //删除
		posts["PUT"] = PPUT       //pubtbput       //修改
	}
	r = posts[params["method"]](params)
	return
}

//判断记录是否存在
func rdexist(params map[string]string) bool {
	filedname := params["existfiled"]
	if filedname == "" {
		return false
	}
	tbname := params["tbname"]
	filedvalue := params[filedname]
	bfiledvalue := Table[tbname].Ifo.TypeChByte(filedname, filedvalue)
	return Table[tbname].Select.WhereIdxExist([]byte(filedname), bfiledvalue)
}

//统计某索引存在的条数
func rdCount(params map[string]string) bool {
	filedname := params["countfiled"]
	if filedname == "" {
		return false
	}
	tblimit := make(map[string]string, 1) //直接在内存设置参数
	tblimit["qz-maxcount"] = "3"          //qz表maxcount=3
	tblimit["wz-maxcount"] = "3"
	tblimit["j-maxcount"] = "3"

	tbname := params["tbname"]
	filedvalue := params[filedname]
	bfiledvalue := Table[tbname].Ifo.TypeChByte(filedname, filedvalue)
	count := Table[tbname].Select.WhereIdxCount([]byte(filedname), bfiledvalue)
	maxcountstr := tblimit[tbname+"-maxcount"] //params["count"]
	maxcount, _ := strconv.Atoi(maxcountstr)
	return count > maxcount
}
func PPOST(params map[string]string) (r xbdb.ReInfo) {
	if params["sj"] == "" {
		params["sj"] = "now()"
	}
	r = Table[params["tbname"]].Ins(params)
	return
}
func PDELETE(params map[string]string) (r xbdb.ReInfo) {
	if rdTimeout(params) {
		r.Info = "已经超时，不能删除。"
		return
	}
	r = Table[params["tbname"]].Del(params["id"])
	return
}

//超时不能删除文章
func rdTimeout(params map[string]string) bool {
	timeoutfield := params["timeoutfield"] //timeout的字段
	if timeoutfield == "" {
		return false
	}
	tblimit := make(map[string]float64, 1) //直接在内存设置参数
	tblimit["j-timeout"] = 3 * 24          //* 60 * 60 * 1000  //j表timeout=3天不能删除
	tblimit["wz-timeout"] = 3 * 24         //* 60 * 60 * 1000 //j表timeout=3天不能删除

	tbname := params["tbname"]
	bid := Table[tbname].Ifo.FieldChByte(Table[tbname].Ifo.Fields[0], params["id"])
	tbd := Table[tbname].Select.Record(bid)
	rdmap := Table[tbname].RDtoMap(tbd.Rd[0])
	rdtime := rdmap[timeoutfield]
	sj, _ := time.ParseInLocation("2006-01-02 15:04:05", rdtime, time.Local)
	return time.Since(sj).Hours() > tblimit[tbname+"-timeout"]
	//	return time.Since(sj).Milliseconds() > tblimit[tbname+"-timeout"]
}
func PPUT(params map[string]string) (r xbdb.ReInfo) {
	r = Table[params["tbname"]].Upd(params)
	return
}
