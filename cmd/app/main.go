package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"os"
	"sync"
	"ticket-tracker/config"
	"ticket-tracker/internal/scheduler/tcdd/v2"
	"ticket-tracker/pkg/db"
	"ticket-tracker/pkg/logger"
	"ticket-tracker/pkg/server"
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
	if err := server.Init(); err != nil {
		fmt.Printf("Server initialization error: %v\n", err)
		os.Exit(1)
	}
}

func runScheduler(wg *sync.WaitGroup) {
	defer wg.Done()
	scheduler := v2.GetTrainSchedulerInstance()
	_, loadStationErr := scheduler.GetStationsOnce()
	if loadStationErr != nil {
		logger.Logger.Errorf("Station load error: %v", loadStationErr)
	}
	err := gocron.Every(15).Seconds().Do(scheduler.Run)
	if err != nil {
		fmt.Printf("Scheduler error: %v\n", err)
		return
	}
	gocron.Start()
}
