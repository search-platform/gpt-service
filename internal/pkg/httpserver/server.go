package httpserver

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
	"github.com/search-platform/gpt-service/internal/pkg/httpserver/middlewares"
	"github.com/search-platform/gpt-service/internal/pkg/httpserver/protofiber"
)

var logger = log.With().Str("name", "http-server").Logger()

type Server struct {
	app *fiber.App

	cfg Config
}

func New(cfg *Config) (*Server, error) {
	s := &Server{
		cfg: *cfg,
	}

	return s, s.Init()
}

func (s *Server) Init() error {
	cfg := fiber.Config{
		ErrorHandler:       middlewares.ErrorHandler,
		Prefork:            false,
		DisableDefaultDate: false,
	}
	if s.cfg.UseProtobuf {
		pf := protofiber.NewProtofiber(protofiber.Opts{
			UseProtoNames: s.cfg.UseProtoNames,
		})
		cfg.JSONEncoder = pf.Marhshal
		cfg.JSONDecoder = pf.Unmarshal
	}

	s.app = fiber.New(cfg)

	s.registerDefaultMiddlewares()

	return nil
}

func (s *Server) Run() error {
	return s.app.Listen(":" + strconv.Itoa(s.cfg.Port))
}

func (s *Server) Shutdown() error {
	if err := s.app.Shutdown(); err != nil {
		logger.Error().Err(err).Msg("failed to shutdown server")
		return err
	}
	return nil
}

func (s *Server) RegisterMiddleware(middlewares ...interface{}) {
	s.app.Use(middlewares...)
}

func (s *Server) registerDefaultMiddlewares() {
	app := s.app

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	app.Use(recover.New())
}

func (s *Server) App() *fiber.App {
	return s.app
}
