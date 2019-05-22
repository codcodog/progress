package progress

import (
	"time"
)

const DEFAULT_REFRESH_INTERVAL = time.Millisecond

type Progress struct {
	bar  *Bar
	done chan bool
}

func New(total int) *Progress {
	p := &Progress{
		NewBar(total),
		make(chan bool),
	}
	p.start()

	return p
}

func (p *Progress) Incr() {
	p.bar.add(1)
}

func (p *Progress) start() {
	go p.listen()
}

func (p *Progress) Stop() {
	p.done <- true
	<-p.done
}

func (p *Progress) listen() {
	for {
		select {
		case <-time.After(DEFAULT_REFRESH_INTERVAL):
			p.bar.print()
		case <-p.done:
			p.bar.print()
			close(p.done)
			return
		}
	}
}
