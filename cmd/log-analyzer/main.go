package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pinjung/log-analyzer/internal/api"
	"github.com/pinjung/log-analyzer/internal/generator"
	"github.com/pinjung/log-analyzer/internal/processor"
	"github.com/pinjung/log-analyzer/internal/storage"
)

func main() {
	log.Println("Log Analyzer MVP started")

	stats := storage.NewStats()
	
	logGenerator := generator.NewGenerator(1000)
	logProcessor := processor.NewProcessor(stats, logGenerator.LogChannel())
	
	server := api.NewServer(stats)
	
	logGenerator.Start(3, 100)
	logProcessor.Start(2)
	
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := server.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				snapshot := stats.GetSnapshot()
				log.Printf("Stats - Total: %d, Errors: %d, Error Rate: %.2f%%", 
					snapshot.TotalCount, 
					snapshot.ErrorCount, 
					float64(snapshot.ErrorCount)/float64(snapshot.TotalCount)*100)
			}
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down...")
	logGenerator.Stop()
	logProcessor.Stop()
	log.Println("Log Analyzer stopped")
}