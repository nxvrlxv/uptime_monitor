package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	//"time"

	"github.com/nxvrlxv/hosts_monitor/internal/checker"
	"github.com/nxvrlxv/hosts_monitor/internal/config"
	"github.com/nxvrlxv/hosts_monitor/internal/handler"
	"github.com/nxvrlxv/hosts_monitor/internal/scheduler"
	"github.com/nxvrlxv/hosts_monitor/internal/storage"
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
	ctx := context.Background()
	st := storage.CreateNewStorage()
	h := handler.NewStorageHandler(st)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /targets", h.GetTargets)
	mux.HandleFunc("GET /targets/{id}", h.GetTargetById)

	go func() {
		fmt.Println("listening 8080")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("Error occured while hosting server: %w", err)
		}
	}()

	round := func(ctx_x context.Context) {
		resultScan := ch.CheckerPool(MaxParallelWorkers, ParsedSites, ctx_x)
		st.FillData(resultScan)
	}

	scheduler.RunTicker(ctx, round, 0)

}
