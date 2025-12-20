package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Tool to monitor and manipulate system resources",
	Run: func(cmd *cobra.Command, args []string) {
		pids, err := listPids()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing PIDs: %v\n", err)
			return
		}

		// Sort PIDs for consistent output
		sortedPIDs := make([]int, 0, len(pids))
		for pid := range pids {
			sortedPIDs = append(sortedPIDs, pid)
		}
		sort.Ints(sortedPIDs)

		// Get stats for each PID
		for _, pid := range sortedPIDs {
			stat, err := getProcessStat(pid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting stat for PID %d: %v\n", pid, err)
				continue
			}
			p := pids[pid]
			p.Stat = stat
			pids[pid] = p
		}

		// Reorder list of processes
		var reorderedPIDs []processInfo
		for _, pid := range sortedPIDs {
			reorderedPIDs = append(reorderedPIDs, pids[pid])
		}

		// Jsonify and print the process list
		output, err := json.MarshalIndent(reorderedPIDs, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling process info: %v\n", err)
			return
		}
		fmt.Println(string(output))
	},
}

func init() {
	rootCmd.AddCommand(systemCmd)
}

type processInfo struct {
	Pid  int           `json:"pid"`
	Stat *ProccessStat `json:"stat,omitempty"`
}

func listPids() (map[int]processInfo, error) {
	// List all directories in /proc that are numeric (PIDs)
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc directory: %w", err)
	}

	pids := make(map[int]processInfo)
	for _, entry := range entries {
		if entry.IsDir() {
			if pid, err := strconv.Atoi(entry.Name()); err == nil {
				pids[pid] = processInfo{Pid: pid}
			}
		}
	}
	return pids, nil
}

type ProccessStat struct {
	Comm  string `json:"comm"`
	State string `json:"state,omitempty"`
	Ppid  int    `json:"ppid,omitempty"`
}

var stateMap = map[rune]string{
	'R': "Running",
	'S': "Sleeping (interruptible)",
	'D': "Sleeping (uninterruptible)",
	'Z': "Zombie",
	'T': "Stopped",
	'W': "Paging",
	'X': "Dead",
	'I': "Idle kernel thread",
}

func getProcessStat(pid int) (*ProccessStat, error) {
	content, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return nil, err
	}
	entries := bytes.Split(content, []byte(" "))
	ppidString := entries[3]
	ppid, err := strconv.Atoi(string(ppidString))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ppid: %w", err)
	}
	return &ProccessStat{
		Comm:  string(entries[1][1 : len(entries[1])-1]),
		State: stateMap[rune(entries[2][0])],
		Ppid:  ppid,
	}, nil
}
