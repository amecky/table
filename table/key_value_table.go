package table

import "time"

type KeyValueTable struct {
	Table *Table
}

func NewKeyValueTable(name string) *KeyValueTable {
	tbl := Table{
		Description:  name,
		TableHeaders: []string{"Name", "Value"},
		Count:        0,
		Limit:        -1,
		Created:      time.Now().Format("2006-01-02 15:04"),
	}
	return &KeyValueTable{
		Table: &tbl,
	}
}

func (kvt *KeyValueTable) AddRow(key, value string) {
	r := kvt.Table.CreateRow()
	r.AddDefaultText(key)
	r.AddDefaultText(value)
}

func (kvt *KeyValueTable) AddFloat(key string, value float64, marker int) {
	r := kvt.Table.CreateRow()
	r.AddDefaultText(key)
	r.AddFloat(value, marker)
}

func (kvt *KeyValueTable) String() string {
	return kvt.Table.String()
}
