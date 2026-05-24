package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Handler функция которая выполняется в транзакции
type Handler func(ctx context.Context) error

// Client для работы с бд
type Client interface {
	DB() DB
	Close() error
}

type TxManager interface {
	ReedCommitted(ctx context.Context, f Handler) error
}

// Query обертка над запросом хранящая сам запрос
// Имя запроса используется для логирования и может использоваться еще гдето
type Query struct {
	Name     string
	QueryRaw string
}

// Transactor интерфейс для работы с транзакциями
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type SQLExecer interface {
	NameExecer
	QueryExecer
}

// NameExecer интерфейс для работы с именованными запросами с помощью тегов в структурах
type NameExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecerинтерфейс для работы с обычными запросами
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger обкртка для проверки соединения с БД
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB интерфейс для работы с БД
type DB interface {
	SQLExecer
	Pinger
	Close()
}
