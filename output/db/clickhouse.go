package db

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var conn driver.Conn

func InitClickhouse() {
	var (
		ctx = context.Background()
		err error
	)

	conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "", // TODO environment
		},
		Debugf: func(format string, v ...any) {
			fmt.Printf(format, v)
		},
	})
	if err != nil {
		log.Fatalf("fail to connect clickhouse, %v", err)
		return
	}
	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Fatalf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
	}
}

func WriteTraffic(host, method, url, ip, status, reqBody, resBody string) {
	err := conn.Exec(context.Background(),
		`INSERT INTO traffic (ts, host, method, url, ip, status, req_body, res_body)
         VALUES (NOW(), ?, ?, ?, ?, ?, ?, ?)`,
		host, method, url, ip, status, reqBody, resBody)
	if err != nil {
		slog.Error("fail to write table traffic", "error", err)
	}
}
