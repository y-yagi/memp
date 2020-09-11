package main

import (
	"fmt"
	"sort"

	"code.cloudfoundry.org/bytefmt"
	"github.com/elastic/go-sysinfo"
	"github.com/y-yagi/color"
)

type process struct {
	name string
	rss  uint64
	vsz  uint64
}

func main() {
	processes, err := sysinfo.Processes()
	if err != nil {
		panic(err)
	}

	rsses := map[string]uint64{}
	vszs := map[string]uint64{}
	for _, process := range processes {
		info, _ := process.Info()
		if len(info.Name) == 0 {
			continue
		}

		memory, _ := process.Memory()
		rsses[info.Exe] += memory.Resident
		vszs[info.Exe] += memory.Virtual
	}

	sortedProcesses := []process{}
	for k, v := range rsses {
		if len(k) != 0 {
			sortedProcesses = append(sortedProcesses, process{name: k, rss: v, vsz: vszs[k]})
		}
	}

	sort.Slice(sortedProcesses, func(i, j int) bool {
		return sortedProcesses[i].rss > sortedProcesses[j].rss
	})

	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	for _, process := range sortedProcesses {
		fmt.Printf("%v %v %v\n", green(process.name), bold(bytefmt.ByteSize(process.rss)), bold(bytefmt.ByteSize(process.vsz)))
	}
}
