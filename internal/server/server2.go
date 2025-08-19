package server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sebzz2k2/vaultic/pkg/utils"
)

type Server struct {
	config *Config
	// engine   storage.StorageEngine
	listener net.Listener

	connections sync.Map
	connCount   int64

	shutdown bool
	mu       sync.RWMutex
	wg       sync.WaitGroup
}

type Config struct {
	Address        string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxConnections int

	MaxMessageSize int
}

func defaultConfig() *Config {
	return &Config{
		Address:        "localhost",
		Port:           5381,
		MaxConnections: 100,
		MaxMessageSize: 1024 * 1024, // 1 MB
	}
}

func New(cfg *Config) (*Server, error) {
	return &Server{
		config: cfg,
	}, nil
}

func (s *Server) Start() error {
	address := net.JoinHostPort(s.config.Address, fmt.Sprintf("%d", s.config.Port))
	listener, err := net.Listen("tcp", address)
	s.listener = listener
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.mu.RLock()
			if s.shutdown {
				s.mu.RUnlock()
				break
			}
			s.mu.RUnlock()

			log.Error().Err(err).Msg("Failed to accept connection")
			continue
		}

		if s.getConnectionCount() >= int64(s.config.MaxConnections) {
			log.Warn().
				Int64("current", s.getConnectionCount()).
				Int("max", s.config.MaxConnections).
				Msg("Connection limit reached, rejecting new connection")

			utils.WriteToClient(conn, "Server busy, please try again later\n")
			conn.Close()
			continue
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer s.removeConnection(conn)

	s.addConnection(conn)
	if s.config.ReadTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))
	}
	if s.config.WriteTimeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
	}

	// client := NewClient(conn, s.config)

	log.Info().
		Str("remote_addr", conn.RemoteAddr().String()).
		Int64("connection_count", s.getConnectionCount()).
		Msg("New client connected")

	// if err := client.Handle(); err != nil {
	// 	log.Error().
	// 		Err(err).
	// 		Str("remote_addr", conn.RemoteAddr().String()).
	// 		Msg("Client handling error")
	// }

	log.Info().
		Str("remote_addr", conn.RemoteAddr().String()).
		Msg("Client disconnected")
}
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.shutdown = true
	s.mu.Unlock()

	log.Info().Msg("Shutting down server")

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing listener")
		}
	}

	s.connections.Range(func(key, value interface{}) bool {
		if conn, ok := value.(net.Conn); ok {
			conn.Close()
		}
		return true
	})

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("All connections closed gracefully")
	case <-ctx.Done():
		log.Warn().Msg("Shutdown timeout reached, forcing close")
	}

	return nil
}

func (s *Server) removeConnection(conn net.Conn) {
	s.connections.Delete(conn)
	atomic.AddInt64(&s.connCount, -1)
	conn.Close()
}
func (s *Server) getConnectionCount() int64 {
	return atomic.LoadInt64(&s.connCount)
}
func (s *Server) addConnection(conn net.Conn) {
	s.connections.Store(conn, conn)
	atomic.AddInt64(&s.connCount, 1)
}
