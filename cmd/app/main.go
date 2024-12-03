package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"os"
	"sync"
	"ticket-tracker/config"
	"ticket-tracker/internal/http"
	scheduler2 "ticket-tracker/internal/infrastructure/scheduler"
	"ticket-tracker/pkg/db"
	"ticket-tracker/pkg/logger"
	"time"
)

func main() {
	loc, err := time.LoadLocation("Europe/Istanbul")
	if err != nil {
		logger.Logger.Fatalf("Timezone yüklenemedi: %v", err)
	}
	time.Local = loc
	var wg sync.WaitGroup
	wg.Add(2)

	go runScheduler(&wg)
	go runServer(&wg)

	wg.Wait()
}

func runServer(wg *sync.WaitGroup) {
	defer wg.Done()
	config.InitConfig()
	if err := db.InitDb(); err != nil {
		logger.Logger.Fatalf("Database initialization error: %v", err)
		os.Exit(1)
	}
	logger.Logger.Info("Server is starting...")
	if err := http.Init(); err != nil {
		fmt.Printf("Server initialization error: %v\n", err)
		os.Exit(1)
	}
}

func runScheduler(wg *sync.WaitGroup) {
	defer wg.Done()
	scheduler := scheduler2.GetTrainSchedulerInstance()
	err := gocron.Every(15).Seconds().Do(scheduler.Run)
	if err != nil {
		fmt.Printf("Scheduler error: %v\n", err)
		return
	}
	gocron.Start()
}
