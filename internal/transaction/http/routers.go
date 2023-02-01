package http

func InitRouters(s *Server) {
	initTransactionRouters(s)
	initBalanceRoters(s)
}

func initTransactionRouters(s *Server) {
	transaction := s.Group("/transaction")

	transaction.POST("/out", s.handler.Output)
	transaction.POST("/in", s.handler.Input)
	transaction.GET("/get/:id", s.handler.GetTransactions)
}

func initBalanceRoters(s *Server) {
	balance := s.Group("/balance")

	balance.GET("/get/:id", s.handler.GetBalance)
}