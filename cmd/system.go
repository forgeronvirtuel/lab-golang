package cmd

import (
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
		fmt.Println("System Resource Monitor -- Not fully implemented yet")
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

		fmt.Printf("Found %d processes\n", len(pids))
		if showPids {
			fmt.Printf("PID\n")
			for _, pid := range sortedPIDs {
				fmt.Printf("%d\n", pid)
			}
		}
	},
}

func init() {
	systemCmd.Flags().BoolVar(&showPids, "show-pids", false, "Show all process IDs")

	rootCmd.AddCommand(systemCmd)
}

type processInfo struct {
	pid int
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
				pids[pid] = processInfo{pid: pid}
			}
		}
	}
	return pids, nil
}
