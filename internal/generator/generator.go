package generator

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/pinjung/log-analyzer/pkg/analyzer"
)

type Generator struct {
	logChan chan *analyzer.Log
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewGenerator(bufferSize int) *Generator {
	ctx, cancel := context.WithCancel(context.Background())
	return &Generator{
		logChan: make(chan *analyzer.Log, bufferSize),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (g *Generator) Start(numWorkers int, logsPerSecond int) {
	for i := 0; i < numWorkers; i++ {
		go g.worker(i, logsPerSecond/numWorkers)
	}
}

func (g *Generator) worker(workerID int, logsPerSecond int) {
	ticker := time.NewTicker(time.Second / time.Duration(logsPerSecond))
	defer ticker.Stop()

	messages := []string{
		"User login successful",
		"Database connection established",
		"Request processed successfully",
		"Cache hit for key: user_123",
		"File uploaded successfully",
		"Payment processed",
		"Database connection failed",
		"Invalid user credentials",
		"Request timeout occurred",
		"Memory allocation failed",
		"Network connection lost",
		"System overload detected",
	}

	levels := []analyzer.LogLevel{
		analyzer.LogLevelInfo,
		analyzer.LogLevelInfo,
		analyzer.LogLevelInfo,
		analyzer.LogLevelWarn,
		analyzer.LogLevelError,
		analyzer.LogLevelDebug,
	}

	for {
		select {
		case <-g.ctx.Done():
			return
		case <-ticker.C:
			level := levels[rand.Intn(len(levels))]
			message := fmt.Sprintf("[Worker-%d] %s", workerID, messages[rand.Intn(len(messages))])
			
			log := analyzer.NewLog(level, message)
			
			select {
			case g.logChan <- log:
			case <-g.ctx.Done():
				return
			}
		}
	}
}

func (g *Generator) LogChannel() <-chan *analyzer.Log {
	return g.logChan
}

func (g *Generator) Stop() {
	g.cancel()
	close(g.logChan)
}