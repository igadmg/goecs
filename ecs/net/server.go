package net

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type Server struct {
	wg sync.WaitGroup
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
		select {
		case <-ctx.Done():
			fmt.Println("Закрываем сервер. Ожидаем задачи")
			s.wg.Wait()
			fmt.Println("сервер закрыт")
		default:
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Ошибка при принятии соединения:", err)
				continue
			}

			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				rpc.ServeConn(conn)
			}()
		}
	}
}
