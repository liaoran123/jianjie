//小白数据库
//表信息
package xbdb

//修改整条记录。等于先删除后添加。
func (t *Table) Updata(vals [][]byte) (r ReInfo) {
	r = t.Delete(vals[0])
	if !r.Succ {
		return
	}
	r = t.Insert(vals)
	return
}

//修改某个某些字段数据。不需要修改的字段索引不会删除。
func (t *Table) Upd(params map[string]string) (r ReInfo) {
	vals := t.StrToByte(params)
	key := JoinBytes(t.Select.GetTbKey(), vals[0])
	data, err := t.Db.Get(key, nil) //获取旧数据
	if err != nil {
		r.Info = err.Error()
		return
	}
	r = t.Act(vals, "delete") //删除旧数据
	if !r.Succ {
		return
	}
	//组织新数据
	var newvals [][]byte
	newvals = append(newvals, vals[0])
	newvals = append(newvals, SplitRd(data)...)
	for i, v := range vals {
		if v != nil {
			newvals[i] = v
		}
	}
	r = t.Insert(newvals) //添加新数据
	return
}
