package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	for {
		processList := getProcesses(port)
		if len(processList) == 0 {
			fmt.Println("No processes found listening on port", port)
			return
		}

		fmt.Println("Processes listening on port", port, ":")
		for i, process := range processList {
			fmt.Printf("[%d] %s\n", i, process)
		}

		fmt.Print("Enter the index of the process to kill or press Enter to kill all: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			for _, process := range processList {
				killProcess(process)
			}
		} else if input == "q" {
			return
		} else {
			index, err := strconv.Atoi(input)
			if err != nil || index < 0 || index >= len(processList) {
				fmt.Println("Invalid input. Please enter a valid index.")
				continue
			}
			killProcess(processList[index])
		}
	}
}

func getProcesses(port string) []string {
	cmd := exec.Command("lsof", "-i", ":"+port)
	output, err := cmd.CombinedOutput()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok && exitError.ExitCode() == 1 {
			return []string{}
		}
		fmt.Println("Error executing lsof:", err)
		os.Exit(1)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) <= 1 {
		return []string{}
	}
	// strip header
	lines = lines[1:]

	nonEmptyLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}
	lines = nonEmptyLines

	return lines
}

func killProcess(process string) {
	fields := strings.Fields(process)
	if len(fields) < 2 {
		fmt.Println("Invalid process information:", process)
		return
	}

	pid := fields[1]
	cmd := exec.Command("kill", "-9", pid)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error killing process:", err)
	} else {
		fmt.Println("Killed process with PID:", pid)
	}
}
