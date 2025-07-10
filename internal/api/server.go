package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pinjung/log-analyzer/internal/storage"
)

type Server struct {
	stats  *storage.Stats
	engine *gin.Engine
}

func NewServer(stats *storage.Stats) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	
	engine.Use(gin.Recovery())
	
	server := &Server{
		stats:  stats,
		engine: engine,
	}
	
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	api := s.engine.Group("/stats")
	{
		api.GET("/errors", s.getErrors)
		api.GET("/all", s.getAllStats)
		api.POST("/reset", s.resetStats)
	}
	
	s.engine.GET("/health", s.health)
}

func (s *Server) getErrors(c *gin.Context) {
	errorCount := s.stats.GetErrorCount()
	c.JSON(http.StatusOK, gin.H{
		"error_count": errorCount,
	})
}

func (s *Server) getAllStats(c *gin.Context) {
	snapshot := s.stats.GetSnapshot()
	c.JSON(http.StatusOK, snapshot)
}

func (s *Server) resetStats(c *gin.Context) {
	s.stats.Reset()
	c.JSON(http.StatusOK, gin.H{
		"message": "Stats reset successfully",
	})
}

func (s *Server) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

func (s *Server) Run(addr string) error {
	return s.engine.Run(addr)
}