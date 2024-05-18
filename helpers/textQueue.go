package helpers

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2/widget"
)

type Queue struct {
	Text []string
	grid *widget.TextGrid
	max  int
}

func NewQueue(max int, grid *widget.TextGrid) *Queue {
	q := &Queue{}
	q.SetMax(max)
	q.SetGrid(grid)

	return q
}

func (q *Queue) SetGrid(grid *widget.TextGrid) {
	q.grid = grid
}

func (q *Queue) Max() int { return q.max }
func (q *Queue) SetMax(max int) {
	if max > 0 {
		q.max = max
	}
}

func (q *Queue) Trim() {
	f := len(q.Text) - q.max
	if f > 0 {
		q.Text = q.Text[f:]
	}
}

func (q *Queue) Append(txt string) {
	txt = fmt.Sprintf("%s %s", time.Now().Format("2006/01/02 15:04:05"), txt)
	q.Text = append(q.Text, txt)
	q.Trim()

	// Match the log package's default format.
	q.grid.SetText(q.String())
}

func (q *Queue) String() string {
	return strings.Join(q.Text, "\n")
}
