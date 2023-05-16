package utils

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Bar struct {
	mu      sync.Mutex
	graph   string
	rate    string
	percent int
	current int
	total   int
	start   time.Time
}

func NewBar(current, total int) *Bar {
	bar := new(Bar)
	bar.current = current
	bar.total = total
	bar.start = time.Now()
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < bar.percent; i += 1 {
		bar.rate += bar.graph
	}
	return bar
}

func NewBarWithGraph(start, total int, graph string) *Bar {
	bar := NewBar(start, total)
	bar.graph = graph
	return bar
}

func (bar *Bar) getPercent() int {
	return int((float64(bar.current) / float64(bar.total)) * 100)
}

func (bar *Bar) getTime() (s string) {
	u := time.Now().Sub(bar.start).Seconds()
	h := int(u) / 3600
	m := int(u) % 3600 / 60
	if h > 0 {
		s += strconv.Itoa(h) + "h "
	}
	if h > 0 || m > 0 {
		s += strconv.Itoa(m) + "m "
	}
	s += strconv.Itoa(int(u)%60) + "s"
	return
}

func (bar *Bar) Load() {
	last := bar.percent
	bar.percent = bar.getPercent()
	for i := 0; i < bar.percent-last; i++ {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-100s]% 3d%%    %2s   %d/%d", bar.rate, bar.percent, bar.getTime(), bar.current, bar.total)
}

func (bar *Bar) Reset(current int) {
	bar.mu.Lock()
	defer bar.mu.Unlock()
	bar.current = current
	bar.Load()

}

func (bar *Bar) Add(i int) {
	bar.mu.Lock()
	defer bar.mu.Unlock()
	bar.current += i
	if bar.current > bar.total {
		bar.current = bar.total
	}
	bar.Load()
}
