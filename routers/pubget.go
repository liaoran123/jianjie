package routers

import (
	"bytes"
	"jianjie/pubgo"
	"jianjie/xbdb"
	"net/http"
	"strconv"
	"strings"
)

// 根据索引查询，返回json
func PubIDXGet(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/PubIDXGet/" + req.Method)
	params := getparas(req)
	idxfield := params["idxfield"]
	idxvalue := params[idxfield]
	tbname := params["tbname"]
	key := Table[tbname].Ifo.FieldChByte(idxfield, idxvalue)
	b, _ := strconv.Atoi(params["b"])
	count, _ := strconv.Atoi(params["count"])
	var sfs []string //strings.Split(params["showFileds"], ",")
	showFileds := []int{}
	if params["showFileds"] != "" {
		sfs = strings.Split(params["showFileds"], ",")
		showFileds = Table[tbname].Ifo.GetFieldIds(sfs) //处理要显示的字段
	}
	_, distinct := params["distinct"]
	tbd := Table[params["tbname"]].Select.WhereIdx([]byte(params["idxfield"]), key, params["asc"] == "1", b, count, showFileds, distinct)
	if tbd != nil {
		var json *bytes.Buffer
		if len(showFileds) == 0 {
			json = Table[tbname].DataToJsonApp(tbd)
		} else {
			ifo := Table[tbname].Ifo.GetIfoForFields(*Table[tbname].Ifo, sfs)
			json = Table[tbname].DataToJsonforIfoApp(tbd, &ifo) //  DataToJson(tbd)
		}
		w.Write(json.Bytes())
		json.Reset()
		xbdb.Bufpool.Put(json)
	} else {
		w.Write([]byte(""))
	}
}

/*
//根据索引结果的集合打开对应的关联表记录，返回json
//对应SQL:select f1,f2 from a where a.id in (select wid from b where idx=xx limit xx)
func PubIdIn(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/PubIDXGetPKs/" + req.Method)
	params := getparas(req)
	idxfield := params["idxfield"]
	idxvalue := params[idxfield]
	tbname := params["tbname"]
	key := Table[tbname].Ifo.FieldChByte(idxfield, idxvalue)
	b, _ := strconv.Atoi(params["b"])
	count, _ := strconv.Atoi(params["count"])
	sfs := strings.Split(params["JsonFiled"], ",")
	showFileds := Table[tbname].Ifo.GetFieldIds(sfs)                                                   //处理要显示的字段
	tbd := Table[tbname].Select.WhereIdx([]byte(params["idxfield"]), key, false, b, count, showFileds) // (select wid from b where idx=xx limit xx)

	if len(tbd.Rd) > 0 {
		//sfs := strings.Split(params["showFileds"], ",")
		//showFileds := Table[tbname].Ifo.GetFieldIds(sfs) //处理要显示的字段
		//var sfs []string //strings.Split(params["showFileds"], ",")
		showFileds = []int{}
		if params["showFileds"] != "" {
			sfs = strings.Split(params["showFileds"], ",")
			showFileds = Table[tbname].Ifo.GetFieldIds(sfs) //处理要显示的字段
		}
		tbd := Table[params["opentb"]].Select.Records(tbd.Rd, showFileds) //select f1,f2 from a where a.id in
		if tbd == nil {
			w.Write([]byte(""))
			return
		}
		var json *bytes.Buffer
		opentb := params["opentb"]
		if len(showFileds) == 0 {
			json = Table[opentb].DataToJson(tbd)
		} else {
			ifo := Table[tbname].Ifo.GetIfoForFields(Table[opentb].Ifo, sfs)
			json = Table[opentb].DataToJsonforIfo(tbd, &ifo) //  DataToJson(tbd)
		}
		w.Write(json.Bytes())
		json.Reset()
		xbdb.Bufpool.Put(json)
	} else {
		w.Write([]byte(""))
	}
}
*/
//根据打开表的Top，顺序或倒序，返回json
func PubGetTB(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/PubGetTB/" + req.Method)
	params := getparas(req)

	b, _ := strconv.Atoi(params["b"])
	count, _ := strconv.Atoi(params["count"])
	tbname := params["tbname"]
	//sfs := strings.Split(params["showFileds"], ",")
	//showFileds := Table[tbname].Ifo.GetFieldIds(sfs) //处理要显示的字段
	var sfs []string //strings.Split(params["showFileds"], ",")
	showFileds := []int{}
	if params["showFileds"] != "" {
		sfs = strings.Split(params["showFileds"], ",")
		showFileds = Table[tbname].Ifo.GetFieldIds(sfs) //处理要显示的字段
	}
	_, distinct := params["distinct"]
	tbd := Table[tbname].Select.FindPrefix([]byte(tbname+xbdb.Split), params["asc"] == "1", b, count, showFileds, distinct)
	if tbd != nil {
		var json *bytes.Buffer
		if len(showFileds) == 0 {
			json = Table[tbname].DataToJsonApp(tbd)
		} else {
			ifo := Table[tbname].Ifo.GetIfoForFields(*Table[tbname].Ifo, sfs)
			json = Table[tbname].DataToJsonforIfoApp(tbd, &ifo) //  DataToJson(tbd)
		}
		w.Write(json.Bytes())
		json.Reset()
		xbdb.Bufpool.Put(json)
	} else {
		w.Write([]byte(""))
	}
}

// 根据打开表的Top，顺序或倒序，返回json
func PubGetTBOne(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/PubGetTBOne/" + req.Method)
	params := getparas(req)

	tbname := params["tbname"]
	id := params["id"]
	bid := Table[tbname].Ifo.FieldChByte(Table[tbname].Ifo.Fields[0], id) //默认第一个必须是主键。Table[tbname].Ifo.Fields[0]
	//sfs := strings.Split(params["showFileds"], ",")
	//showFileds := Table[tbname].Ifo.GetFieldIds(sfs) //处理要显示的字段
	var sfs []string //strings.Split(params["showFileds"], ",")
	showFileds := []int{}
	if params["showFileds"] != "" {
		sfs = strings.Split(params["showFileds"], ",")
		showFileds = Table[tbname].Ifo.GetFieldIds(sfs) //处理要显示的字段
	}
	tbd := Table[tbname].Select.Record(bid, showFileds)
	if tbd != nil {
		var json *bytes.Buffer
		if len(showFileds) == 0 {
			json = Table[tbname].DataToJsonApp(tbd)
		} else {
			ifo := Table[tbname].Ifo.GetIfoForFields(*Table[tbname].Ifo, sfs)
			json = Table[tbname].DataToJsonforIfoApp(tbd, &ifo) //  DataToJson(tbd)
		}
		w.Write(json.Bytes())
		json.Reset()
		xbdb.Bufpool.Put(json)
	} else {
		w.Write([]byte(""))
	}
}
