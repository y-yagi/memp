package main

import (
	"fmt"
	"sort"

	"code.cloudfoundry.org/bytefmt"
	"github.com/elastic/go-sysinfo"
	"github.com/y-yagi/color"
)

type process struct {
	name   string
	memory uint64
}

func main() {
	processes, err := sysinfo.Processes()
	if err != nil {
		panic(err)
	}

	memories := map[string]uint64{}
	for _, process := range processes {
		info, _ := process.Info()
		if len(info.Name) == 0 {
			continue
		}

		memory, _ := process.Memory()
		memories[info.Name] += memory.Resident
	}

	sortedProcesses := []process{}
	for k, v := range memories {
		if len(k) != 0 {
			sortedProcesses = append(sortedProcesses, process{name: k, memory: v})
		}
	}

	sort.Slice(sortedProcesses, func(i, j int) bool {
		return sortedProcesses[i].memory > sortedProcesses[j].memory
	})

	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	for _, process := range sortedProcesses {
		fmt.Printf("%v %v\n", green(process.name), bold(bytefmt.ByteSize(process.memory)))
	}
}
