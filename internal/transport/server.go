package transport

import (
	"bufio"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/scriptdealer/tcp-pow/internal/config"
	"github.com/scriptdealer/tcp-pow/internal/observability"
	"github.com/scriptdealer/tcp-pow/internal/services"
)

type Server struct {
	address  string
	wisdom   []string
	listener net.Listener
	Log      *slog.Logger
	powAlg   *services.BasicHashcash
	limiter  *services.ConnectionPool
}

func NewServer() *Server {
	cfg := config.New()
	logger, _ := observability.NewLogger(slog.LevelDebug)

	return &Server{
		address: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		powAlg:  services.NewHC(5, time.Hour),
		limiter: services.NewConnectionPool(logger, 100, 5, time.Second*10),
		Log:     logger,
		wisdom: []string{
			"The only way to do great work is to love what you do.",
			"Life is what happens when you're busy making other plans.",
			"Spread love everywhere you go.",
			"The future belongs to those who believe in the beauty of their dreams.",
			"Whoever is happy will make others happy too.",
			"The journey of a thousand miles begins with a single step.",
			"In the middle of difficulty lies opportunity.",
			"The only limit to our realization of tomorrow will be our doubts of today.",
			"Happiness is not something ready-made. It comes from your own actions.",
			"Success is not final, failure is not fatal: it is the courage to continue that counts.",
		},
	}
}

func (s *Server) getRandomQuote() string {
	return s.wisdom[rand.Intn(len(s.wisdom))]
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}
	s.listener = listener

	fmt.Printf("Server listening on %s\n", s.address)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				continue
			}

			return fmt.Errorf("error accepting connection: %w", err)
		}
		go s.handle(conn)
	}
}

func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) ResetLimiter() {
	s.limiter.Reset()
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
	difficulty := s.limiter.CountRequests(ip) // for each recorded request, difficulty goes up by 1
	if difficulty == services.ConnLimitReached {
		conn.Write([]byte("you're wise enough!\n"))
	}
	s.Log.Info("New challenge", "ip", ip, "difficulty", difficulty)
	challenge := s.powAlg.NewChallengeWithDifficulty(ip, difficulty)
	conn.Write([]byte(challenge + "\n")) //nolint:errcheck

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)

		return
	}

	if s.powAlg.Verify(challenge, strings.TrimSpace(response)) {
		quote := s.getRandomQuote()
		conn.Write([]byte(quote + "\n")) //nolint:errcheck
	} else {
		conn.Write([]byte("Invalid proof of work\n")) //nolint:errcheck
		// ToDo: s.limiter.AddPenalty(ip)
	}
}
