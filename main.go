package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/x/term"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/process"
)

var usage int
var e error
var w, h int
var blocks []rune = []rune{' ', '▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
var hist []float64

func main() {
	for {
		w, h, e = term.GetSize(os.Stdout.Fd())
		//		usage, e := lscpu()
		//		if e != nil {
		//			fmt.Printf("E: %v\r", e)
		//			continue
		//		}
		//		fmt.Printf("%v\r", usage)

		time.Sleep(500 * time.Millisecond)

		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		fmt.Printf("Cpu usage:\n\n")
		gdraw()

		for r, num := range proc(w / 2) {
			if r >= h {
				break
			}
			mv(w/2+1, r+1)
			fmt.Print(num)
		}
		use, e := cpu.Percent(0, false)

		if e == nil && len(use) > 0 {
			hist = append(hist, use[0])
			if len(hist) > w {
				hist = hist[1:]
			}
		}
		for y := 1; y <= h; y++ {
			mv(w/2, y)
			fmt.Print("│")
		}
	}

}

func check(err error) {
	if err != nil {
		return
	}

}

func ram() {
}
func lscpu() (float64, error) {
	cpu, e := cpu.Percent(0, true)

	if e != nil {
		return 0, e
	}
	return cpu[0], nil
}

func proc(w int) []string {
	var screen []string
	proc, e := process.Processes()
	check(e)
	//	fmt.Printf("Total: %v\n\n", len(proc))
	for i := range len(proc) {
		n, e := proc[i].Name()
		check(e)
		cm, e := proc[i].Cmdline()
		check(e)
		out := fmt.Sprintf("%-8d %-20s (%s)", proc[i].Pid,
			n, cm)
		if len(out) > w {
			out = out[:w-3] + "..."
		}
		screen = append(screen, out)
	}
	return screen

}

func mv(x, y int) {
	if y < 0 {
		cmd := exec.Command("stty", "size")
		cmd.Stdin = os.Stdin
		out, e := cmd.Output()
		check(e)
		var rows, cols int
		fmt.Sscanf(string(out), "%d %d", &rows, &cols)
		fmt.Printf("\033[%d;%dH", rows, x)

	} else {
		fmt.Printf("\033[%d;%dH", y, x)
	}
}

func uindex(usage float64) rune {
	inx := int(usage / 100 * float64(len(blocks)-1))
	return blocks[inx]
}
func gdraw() {
	gy := h / 4
	gh := gy - 3
	if h-gy < gh {
		gh = h - gy
	}

	gw := w/2 - 1
	if gh < 1 || gw < 1 {
		return
	}
	if len(hist) > gw {
		hist = hist[len(hist)-gw:]
	}

	upc := len(blocks) - 1

	mv(1, gy)
	for x := 1; x <= gw; x++ {
	}

	for col, samp := range hist {
		x := gw - len(hist) + col + 1
		height := samp / 100 * float64(gh*upc)

		color := "\x1b[32m"
		if samp >= 60 {
			color = "\x1b[33m"
		}
		if samp >= 85 {
			color = "\x1b[31m"
		}

		for row := 1; row <= gh; row++ {
			fill := int(math.Ceil(height - float64((row-1)*upc)))
			if fill < 0 {
				fill = 0
			}
			if fill > upc {
				fill = upc
			}
			if fill == 0 {
				continue
			}

			mv(x, gy-row)
			fmt.Printf("%s%c\x1b[0m", color, blocks[fill])

			mv(x, gy+row-1)
			if fill == upc {
				fmt.Printf("%s%c\x1b[0m", color, blocks[upc])
			} else {
				fmt.Printf("%s\x1b[7m%c\x1b[0m", color, blocks[upc-fill])
			}
		}
	}
}
