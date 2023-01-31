//小白数据库
//表信息
package xbdb

import (
	"bytes"
	"strconv"
	"strings"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var actmap map[string]func(k, v []byte) (r ReInfo)

//表的类
type Table struct {
	Name   string
	Db     *leveldb.DB
	Select *Select
	Ac     *Autoinc
	Ifo    TableInfo
}

func NewTable(db *leveldb.DB, name string) *Table {
	return &Table{
		Name:   name,
		Db:     db,
		Select: NewSelect(name, db),
		Ifo:    NewTableInfo(db).GetInfo(name),
	}
}

//遍历分词
func (t *Table) ForDisparte(nr string, ftlen int) (disparte []string) {
	var knr string //, fid
	var ml, cl int
	var r, idxstr []rune
	r = []rune(nr)
	cl = len([]rune(nr))
	for cl > 0 {
		if cl >= ftlen {
			ml = ftlen
		} else {
			ml = cl
		}
		idxstr = r[:ml]
		knr = string(idxstr)
		disparte = append(disparte, knr)
		r = r[1:]
		cl = len(r)
	}
	return
}

//添加或删除一条记录，以及相关索引等所有数据
func (t *Table) Act(vals [][]byte, Act string) (r ReInfo) {
	var updatefield []bool
	return t.Acts(vals, Act, updatefield)
}

//添加或删除一条记录，以及相关索引等所有数据等事务
//updatefield,修改时用。用于记录那个字段需要修改。与字段一一对应。
//修改某个某些字段时，不用把所有索引都删除再重新添加，导致性能不高和不灵活。
func (t *Table) Acts(vals [][]byte, Act string, updatefield []bool) (r ReInfo) {
	if actmap == nil {
		actmap = map[string]func(k, v []byte) (r ReInfo){
			"insert": t.put,
			"delete": t.del,
		}
	}
	r = t.ActPK(vals, Act)
	if !r.Succ {
		return
	}
	//添加表索引
	idx := -1
	var ivs []string
	idxfields := ""
	var idxval, idxvals []byte
	for _, iv := range t.Ifo.Idxs {
		if iv == "" {
			continue
		}
		ivs = strings.Split(iv, ",")
		if len(updatefield) > 0 { //len(updatefield) == 0 则是添加，否则是修改的情况
			if !isUpdateField(ivs, updatefield) {
				continue //不是修改字段则退出，不用添加或删除原有索引
			}
		}
		for i := 0; i < len(ivs); i++ { //组织单个或组合索引key
			idx, _ = strconv.Atoi(ivs[i])
			idxval = vals[idx]
			idxfields += t.Ifo.Fields[idx]       //累加
			idxvals = JoinBytes(idxvals, idxval) //累加
			if i != len(ivs)-1 {                 //不是末尾，则加分隔符
				idxfields += IdxSplit
				idxvals = JoinBytes(idxvals, []byte(IdxSplit))
			}
		}
		r = t.ActIDX([]byte(idxfields), idxvals, vals[0], []byte{}, Act)
		if !r.Succ {
			return
		}
		//重置
		idxfields = ""
		idxvals = idxvals[:0]
	}
	//添加表全文索引
	ftlen, _ := strconv.Atoi(t.Ifo.FTLen)
	var ftIdx []string
	for _, i := range t.Ifo.FullText {
		if i == "" {
			continue
		}
		idx, _ = strconv.Atoi(i)
		if len(updatefield) > 0 { //len(updatefield) == 0 则是添加，否则是修改的情况
			if !updatefield[idx] { //非修改字段不添加/删除索引
				continue
			}
		}
		ftIdx = t.ForDisparte(string(vals[idx]), ftlen)
		for p, f := range ftIdx {
			t.ActIDX([]byte(t.Ifo.Fields[idx]), []byte(f), vals[0], IntToBytes(p), Act)
			if !r.Succ {
				return
			}
		}
	}
	r.Succ = true
	r.Info = "ok"
	return
}

//是不是修改字段。
//由于支持组合索引，故而需要循环，看似复制些
func isUpdateField(ks []string, updatefield []bool) bool {
	isf := false
	idx := 0
	for i := 0; i < len(ks); i++ { //单个或组合索引。由于支持组合索引，故而需要循环，看似复制些
		idx, _ = strconv.Atoi(ks[i])
		isf = isf || updatefield[idx]
	}
	return isf
}

//添加/删除主键数据，即添加/删除一条记录。
func (t *Table) ActPK(vals [][]byte, Act string) (r ReInfo) {
	key := t.Select.GetPkKey(vals[0])
	r = actmap[Act](key, bytes.Join(vals[1:], []byte(Split)))
	if !r.Succ {
		return
	}
	return
}

//添加/删除一条索引数据。
func (t *Table) ActIDX(idxfield, idxvalue, pkvalue, val []byte, Act string) (r ReInfo) {
	//bySplit := []byte(Split)
	//k=ca,fid-3-7 v=
	//prefix := JoinBytes([]byte(t.Ifo.Name+","), idxFieldname, bySplit, idxFieldvalue, bySplit, PKvalue)
	key := t.Select.GetIdxPrefixKey(idxfield, idxvalue, pkvalue) //getIdxPrefixKey
	r = actmap[Act](key, val)
	if !r.Succ {
		return
	}
	return
}

//字符串转byte
//params字段对应的字符串map
func (t *Table) StrToByte(params map[string]string) (r [][]byte) {
	for i, v := range t.Ifo.Fields {
		r = append(r, t.Ifo.TypeChByte(t.Ifo.FieldType[i], params[v]))
	}
	return
}

//将记录转换为map
func (t *Table) RDtoMap(Rd []byte) (r map[string]string) {
	r = make(map[string]string, len(t.Ifo.Fields))
	vs := t.Split(Rd)
	for i, v := range vs {
		r[t.Ifo.Fields[i]] = t.Ifo.ByteChString(t.Ifo.FieldType[i], v) //将包括分隔符的转义数据恢复
	}
	return
}

//将记录分开并转义数据恢复
func (t *Table) Split(Rd []byte) (r [][]byte) {
	r = SplitRd(Rd)
	return
}

var Bufpool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func (t *Table) DataToJson(tbd *TbData) (r *bytes.Buffer) {
	r = t.DataToJsonforIfo(tbd, &t.Ifo)
	return
}

func (t *Table) DataToJsonforIfo(tbd *TbData, Ifo *TableInfo) (r *bytes.Buffer) {
	if tbd == nil {
		return
	}
	r = Bufpool.Get().(*bytes.Buffer)
	if r.Len() > 0 { //保证数据不混乱
		r.Reset()
	}
	var rdmap map[string]string
	jsonstr := ""
	valstr := ""
	r.WriteString("{\"result\":[")
	for j, v := range tbd.Rd {
		if v == nil {
			continue
		}
		r.WriteString("{")
		rdmap = t.RDtoMap(v)
		/*
			[{"id":2,"title":"金刚经","fid":1,"isleaf":"0"},
			{"id":3,"title":"六祖坛经","fid":1,"isleaf":"0"}]
		*/
		for i, fv := range Ifo.FieldType {
			switch fv {
			case "string", "time", "bool":
				valstr = strconv.Quote(rdmap[Ifo.Fields[i]])
			default:
				valstr = rdmap[Ifo.Fields[i]]
			}
			jsonstr = "\"" + Ifo.Fields[i] + "\":" + valstr
			if i != len(rdmap)-1 {
				jsonstr += ","
			}
			r.WriteString(jsonstr)
		}
		r.WriteString("}")
		if j != len(tbd.Rd)-1 {
			r.WriteString(",")
		}
	}
	r.WriteString("]}")
	tbd.Release()
	return
}

//获取字段在表中的索引id
func (t *Table) GetFieldIdx(field string) int {
	for i, fv := range t.Ifo.Fields {
		if field == fv { //得到pv在Fields中索引id
			return i
		}
	}
	return -1
}

//根据字段索引判断是否索引
func (t *Table) FieldIsIdx(idx int) bool {
	var i int
	if idx == -1 {
		return false
	}
	for _, v := range t.Ifo.Idxs {
		i, _ = strconv.Atoi(v)
		if idx == i {
			return true
		}
	}
	return false
}
