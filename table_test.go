package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {

	tbl := NewTable("Test", []string{"1", "2"})
	if len(tbl.Headers) != 2 {
		t.Errorf("got %d, wanted %d", len(tbl.Headers), 2)
	}
}

func TestCreateRow(t *testing.T) {

	tbl := NewTable("Test", []string{"1", "2"})
	r := tbl.CreateRow()
	r.AddDefaultText("Test1")
	r.AddDefaultText("Test2")
	assert.Equal(t, 1, len(tbl.Rows))
	tr := tbl.Rows[0]
	assert.Equal(t, 2, len(tr.Cells))
	assert.Equal(t, "Test1", tr.Cells[0].Text)
	assert.Equal(t, "Test2", tr.Cells[1].Text)
}
