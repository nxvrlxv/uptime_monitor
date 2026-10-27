package main

import (
	"context"
	"fmt"

	"github.com/nxvrlxv/hosts_monitor/internal/checker"
	"github.com/nxvrlxv/hosts_monitor/internal/config"
)

var (
	FilePath           string = "/Users/enjoyer/go_edu/hosts_monitor/targets.json"
	MaxParallelWorkers        = 5
)

func main() {
	ParsedSites, err := config.ParseTargets(FilePath)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	ch := checker.NewChecker()
	ctx := context.Background()

	resultScan := checker.CheckerPool(MaxParallelWorkers, ParsedSites, ctx, *ch)

	for _, res := range resultScan {
		fmt.Printf("Host: %v\nIs Available: %v\nResponseTime: %v\nStatus: %v\n", res.Host, res.IsUp, res.ResponseTime, res.Status)
		fmt.Println("----------------")
	}

}
