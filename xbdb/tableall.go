package xbdb

import "github.com/syndtr/goleveldb/leveldb"

//对表结构以及对应的数据进行同步更新
type TableAll struct {
	Table     *Table
	TableInfo *TableInfo
}

func NewTableAll(db *leveldb.DB, tbname string) *TableAll {
	Table := NewTable(db, tbname)
	TableInfo := NewTableInfo(db)
	return &TableAll{
		Table:     Table,
		TableInfo: TableInfo,
	}
}

//重新添加表结构，实质就是先删除后添加。
//仅仅在没有数据时有效和使用。
func (t *TableAll) ReInsIfo(ifo TableInfo) (r ReInfo) {
	r = t.TableInfo.Del(t.Table.Name)
	if !r.Succ {
		return
	}
	r = t.TableInfo.Create(t.Table.Name, ifo.FTLen, ifo.Fields, ifo.Fields, ifo.FieldType, ifo.Idxs, ifo.FullText)
	return
}
