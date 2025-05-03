package net

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
)

type Server struct {
}

func (s *Server) Listen(rcvr any, ctx context.Context) {
	cfg := net.ListenConfig{}
	listener, err := cfg.Listen(ctx, "tcp", ":1234")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Сервер запущен на порту 1234")

	rpc.Register(rcvr)

	// Принимаем соединения и обслуживаем RPC-запросы
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
