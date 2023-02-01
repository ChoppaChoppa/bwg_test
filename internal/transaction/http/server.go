package http

import "github.com/labstack/echo/v4"

type IHandler interface {
	Input(c echo.Context) error
	Output(c echo.Context) error
	GetTransactions(c echo.Context) error
	GetBalance(c echo.Context) error
}

type Server struct {
	*echo.Echo
	host    string
	handler IHandler
}

func New(host string, handler IHandler) *Server {
	return &Server{
		Echo:    echo.New(),
		host:    host,
		handler: handler,
	}
}

func (s *Server) Run() error {
	InitRouters(s)

	s.HideBanner = true
	if err := s.Echo.Start(s.host); err != nil {
		return err
	}

	return nil
}
