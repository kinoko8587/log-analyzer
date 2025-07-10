package processor

import (
	"context"
	"log"

	"github.com/pinjung/log-analyzer/internal/storage"
	"github.com/pinjung/log-analyzer/pkg/analyzer"
)

type Processor struct {
	stats   *storage.Stats
	ctx     context.Context
	cancel  context.CancelFunc
	logChan <-chan *analyzer.Log
}

func NewProcessor(stats *storage.Stats, logChan <-chan *analyzer.Log) *Processor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Processor{
		stats:   stats,
		ctx:     ctx,
		cancel:  cancel,
		logChan: logChan,
	}
}

func (p *Processor) Start(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go p.worker(i)
	}
}

func (p *Processor) worker(workerID int) {
	for {
		select {
		case <-p.ctx.Done():
			return
		case logEntry, ok := <-p.logChan:
			if !ok {
				return
			}
			
			p.processLog(logEntry)
		}
	}
}

func (p *Processor) processLog(logEntry *analyzer.Log) {
	p.stats.RecordLog(logEntry)
	
	if logEntry.IsError() {
		log.Printf("ERROR detected: %s", logEntry.Message)
	}
}

func (p *Processor) Stop() {
	p.cancel()
}

func (p *Processor) GetStats() *storage.Stats {
	return p.stats
}