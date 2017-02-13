package htmlPDF

type Table struct {
	n        int64 //tr
	m        int64 //td
	actulaTd map[int64]string
	content  map[int64]map[int64]string
	col      map[int64]float64
}

func NewTable() *Table {
	return &Table{
		n:        0,
		m:        0,
		actulaTd: nil,
		content:  make(map[int64]map[int64]string),
		col:      make(map[int64]float64),
	}
}

func (t *Table) startTr() {
	t.actulaTd = make(map[int64]string)
	t.m = 0
	t.n++
}

func (t *Table) endTr() {
	t.content[t.n] = t.actulaTd
}

func (t *Table) addTd(text string) {
	size := pdf.GetStringWidth(text)
	if size > t.col[t.m] {
		t.col[t.m] = size
	}
	t.actulaTd[t.m] = text
	t.m++
}

func (t *Table) printSelf() {
	var i, j int64
	for i = 1; i < int64(len(t.content))+1; i++ {
		for j = 0; j < int64(len(t.content[i])); j++ {
			td := t.content[i][j]
			pdf.SetFont(pageFontFamily, pageFontStyle, pageFontSize)
			pdf.Text(pdf.GetX(), pdf.GetY(), td)
			x := pdf.GetX()
			pdf.SetX(x + t.col[j] + 2)
			newline = false
		}
		setNewLine()
	}
}
