package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"os/exec"
	"log"
)

func CrawlProc(parent int) (err error, processes []int) {
	children := map[int][]int{}
	positions, _ := os.ReadDir("/proc")

	for _, position := range positions {
		process, err := strconv.Atoi(position.Name())
		if err != nil {continue}
		data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", process))
		if err != nil {continue}

		infoString := string(data)
		fields := strings.Fields(infoString[strings.LastIndex(infoString, ")")+1:])
		if len(fields) < 2 {continue}

		pPosition, _ := strconv.Atoi(fields[1])
		children[pPosition] = append(children[pPosition], process)
	}

	final := []int{parent}
	for i := 0; i < len(final); i++ {
		final = append(final, children[final[i]]...)
	}
	return nil, final
}

type WindowQueryData struct {
	Data_PosX	    int
	Data_PosY	    int
	Data_Height	    int
	Data_Width          int
}

func SetWindow(parent int, data WindowQueryData) error {
	var preference = GetUserPreference()
	if !preference.RoseStrap.AllowWinTrans {return fmt.Errorf("SetWindow disabled by client")}
	log.Println("Modifying Window Options")
	var err, processes = CrawlProc(parent)
	if err != nil {return fmt.Errorf("Could not set window size! Reason: %w", err)}
	for _, process := range processes {
		var out, err = exec.Command("xdotool", "search", "--onlyvisible", "--pid", strconv.Itoa(process)).Output()
		if err != nil {continue}
		for _, program := range strings.Fields(string(out)) {
			if data.Data_Height != 0 && data.Data_Width != 0 {
				if err := exec.Command("xdotool", "windowsize", program,
					strconv.Itoa(data.Data_Width), strconv.Itoa(data.Data_Height)).Run(); err != nil {
						log.Printf("windowsize failed: %v\n", err)
					}
			}
			if data.Data_PosX != 0 && data.Data_PosY != 0 {
				if err := exec.Command("xdotool", "windowmove", program,
					strconv.Itoa(data.Data_PosX), strconv.Itoa(data.Data_PosY)).Run(); err != nil {
						log.Printf("windowmove failed: %v\n", err)
					}
			}
		}
	}
	return nil
}
