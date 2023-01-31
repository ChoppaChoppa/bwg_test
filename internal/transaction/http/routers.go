package http

func InitRouters(s *Server) {
	transaction := s.Group("/transaction")

	transaction.POST("/out", s.handler.Output)
	transaction.POST("/in", s.handler.Input)
}
