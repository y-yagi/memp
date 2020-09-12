package main

import (
	"fmt"
	"os"
	"sort"

	"code.cloudfoundry.org/bytefmt"
	"github.com/elastic/go-sysinfo"
	"github.com/olekukonko/tablewriter"
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

	data := [][]string{}
	for _, process := range sortedProcesses {
		data = append(data, []string{process.name, fmt.Sprintf("%v", bytefmt.ByteSize(process.rss))})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Exe", "Rss"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data)
	table.Render()
}
