package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type serverResult struct {
	Name   string
	Result string
}

func pingServer(server string) (string, error) {
	out, err := exec.Command("ping", server).Output()
	if err != nil {
		return "", err
	}

	log.Println(string(out))

	if strings.Contains(string(out), "Lost = 4,") {
		return "Failed", nil
	}
	return "Success", nil

}

// readServersFromFile reads the server list from a txt file
func readServersFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to read file: %s", err)
	}
	defer file.Close()

	var servers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		servers = append(servers, scanner.Text())
	}
	return servers, scanner.Err()
}

// writeResultsToCSV writes the ping results to a CSV file
func writeResultsToCSV(results []serverResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"Server", "Result"})

	for _, result := range results {
		err := writer.Write([]string{result.Name, result.Result})
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	fmt.Println("Starting ICMP ping...")

	servers, err := readServersFromFile("servers.txt")
	if err != nil {
		log.Fatal("Unable to find servers.txt file in the current directory. Error message: ", err)
	}

	var results []serverResult
	for _, server := range servers {
		result, err := pingServer(server)
		if err != nil {
			log.Printf("Error executing ping command %s: %s", server, err)
		}

		results = append(results, serverResult{Name: server, Result: result})
	}

	fileName := fmt.Sprintf("serverpulse_results_%s.csv", time.Now().Local().Format("01-02-2006"))
	err = writeResultsToCSV(results, fileName)
	if err != nil {
		log.Fatal("Error writing results file: ", err)
	}

	fmt.Println("Finished ICMP ping...")
}

func init() {
	f, err := os.OpenFile("serverpulse.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		os.Exit(1)
		return
	}

	log.SetOutput(f)
}
