package progress

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

const (
	WIDTH    = 70
	FILL     = '='
	HEAD     = '>'
	EMPTY    = '-'
	LEFTEND  = '['
	RIGHTEND = ']'
)

type Bar struct {
	total      int
	current    int
	firstPrint bool
	mtx        *sync.Mutex
}

func NewBar(total int) *Bar {
	return &Bar{
		total:      total,
		firstPrint: true,

		mtx: &sync.Mutex{},
	}
}

func (b *Bar) add(n int) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.current += n
	if b.current > b.total {
		fmt.Println("errors: current value is greater than total value")
		os.Exit(1)
	}
}

func (b *Bar) print() {
	bar := b.string()

	// \033[nA 光标上移n行
	// \033[K 清除从光标到行尾的内容
	if !b.firstPrint {
		fmt.Print("\033[1A\033[2K")
	}

	fmt.Println(bar, b.completedPercentString())
	b.firstPrint = false
}

func (b *Bar) bytes() []byte {
	completedWidth := int(float64(WIDTH) * (b.completedPercent() / 100.00))

	var buf bytes.Buffer
	for i := 0; i < completedWidth; i++ {
		buf.WriteByte(FILL)
	}
	for i := 0; i < (WIDTH - completedWidth); i++ {
		buf.WriteByte(EMPTY)
	}

	bar := buf.Bytes()
	bar[0], bar[len(bar)-1] = LEFTEND, RIGHTEND

	if completedWidth > 0 && completedWidth < WIDTH {
		bar[completedWidth-1] = HEAD
	}

	return bar
}

func (b *Bar) string() string {
	return string(b.bytes())
}

func (b *Bar) completedPercent() float64 {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	return (float64(b.current) / float64(b.total)) * 100.00
}

func (b *Bar) completedPercentString() string {
	return fmt.Sprintf("%3.f%%", b.completedPercent())
}
