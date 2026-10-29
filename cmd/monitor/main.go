package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nxvrlxv/hosts_monitor/internal/checker"
	"github.com/nxvrlxv/hosts_monitor/internal/config"
	"github.com/nxvrlxv/hosts_monitor/internal/scheduler"
)

var (
	FilePath           string = "/Users/enjoyer/go_edu/hosts_monitor/targets/targets.json"
	MaxParallelWorkers        = 5
)

func main() {
	ParsedSites, err := config.ParseTargets(FilePath)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	ch := checker.NewChecker()
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	round := func(ctx_x context.Context) {
		resultScan := ch.CheckerPool(MaxParallelWorkers, ParsedSites, ctx_x)
		for _, res := range resultScan {
			fmt.Printf("Host: %v\nIs Available: %v\nResponseTime: %v\nStatus: %v\n", res.Host, res.IsUp, res.ResponseTime, res.Status)
			fmt.Println("----------------")
		}
	}

	scheduler.RunTicker(ctx, round, 0)

}
