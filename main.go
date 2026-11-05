package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {
	for {
		shit, err := cpu.Percent(time.Second, false)
		check(err)
		mem, err := mem.VirtualMemory()
		check(err)
		fmt.Printf("%v\n", shit[0])
		fmt.Printf("%v\n\n", mem.Total)
	}

}

func check(err error) {
	if err != nil {
		fmt.Printf("E: %v", err)
		return
	}

}

func ram() {
	fmt.Printf("%v\n", mem.ExVirtualMemory.ActiveAnon4)
}
