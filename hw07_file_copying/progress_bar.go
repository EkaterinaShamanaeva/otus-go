package main

import "fmt"

type Bar struct {
	total       int64  // max task progress
	current     int64  // current progress
	percent     int64  // progress percentage
	prevPercent int64  // prev percentage
	rate        string // the actual progress bar to be printed
	graph       string // symbol for printing
}

func (bar *Bar) New(total int64) {
	bar.total = total
	bar.current = 0
	bar.percent = bar.getPercent()
	bar.prevPercent = 0
	bar.rate = ""
	bar.graph = "█"
}

func (bar *Bar) getPercent() int64 {
	return int64((float32(bar.current) / float32(bar.total)) * 100)
}

func (bar *Bar) Add(step int64) {
	bar.current += step
	bar.prevPercent = bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != bar.prevPercent {
		var i int64
		for ; i < bar.percent-bar.prevPercent; i++ {
			bar.rate += bar.graph
		}
		fmt.Printf("\r[%-100s]%3d%% %8d/%d", bar.rate, bar.percent, bar.current, bar.total)
	}
}

func (bar *Bar) Finish() {
	fmt.Println()
}
