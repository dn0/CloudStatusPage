package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	dummyCommandTag = "DUMMY"
)

type dummyClient struct{}

type dummyBatchResults struct{}

type dummyRows struct{}

type dummyRow struct{}

type dummyTx struct{}

func newDummyClient() *dummyClient {
	return &dummyClient{}
}

func (d *dummyClient) Begin(_ context.Context) (pgx.Tx, error) {
	slog.Debug("db.dummyClient.Begin()")
	return &dummyTx{}, nil
}

func (d *dummyClient) Exec(_ context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	slog.Debug("db.dummyClient.Exec()", "sql", sql, "arguments", arguments)
	return pgconn.NewCommandTag(dummyCommandTag), nil
}

func (d *dummyClient) Query(_ context.Context, sql string, args ...any) (pgx.Rows, error) {
	slog.Debug("db.dummyClient.Query()", "sql", sql, "args", args)
	return &dummyRows{}, nil
}

func (d *dummyClient) QueryRow(_ context.Context, sql string, args ...any) pgx.Row {
	slog.Debug("db.dummyClient.QueryRow()", "sql", sql, "args", args)
	return &dummyRow{}
}

func (d *dummyClient) SendBatch(_ context.Context, b *pgx.Batch) pgx.BatchResults {
	slog.Debug("db.dummyClient.SendBatch()", "batch", b)
	return &dummyBatchResults{}
}

func (d *dummyClient) Ping(_ context.Context) error {
	slog.Debug("db.dummyClient.Ping()")
	return nil
}

func (d *dummyClient) Close() {
	slog.Debug("db.dummyClient.Close()")
}

func (d *dummyBatchResults) Exec() (pgconn.CommandTag, error) {
	slog.Debug("db.dummyBatchResults.Exec()")
	return pgconn.NewCommandTag(dummyCommandTag), nil
}

func (d *dummyBatchResults) Query() (pgx.Rows, error) {
	slog.Debug("db.dummyBatchResults.Query()")
	return &dummyRows{}, nil
}

func (d *dummyBatchResults) QueryRow() pgx.Row {
	slog.Debug("db.dummyBatchResults.QueryRow()")
	return nil
}

func (d *dummyBatchResults) Close() error {
	slog.Debug("db.dummyBatchResults.Close()")
	return nil
}

func (d *dummyRows) Close() {}

func (d *dummyRows) Err() error {
	return nil
}

func (d *dummyRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag(dummyCommandTag)
}

func (d *dummyRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (d *dummyRows) Next() bool {
	return false
}

func (d *dummyRows) Scan(_ ...any) error {
	return nil
}

func (d *dummyRows) Values() ([]any, error) {
	return nil, nil
}

func (d *dummyRows) RawValues() [][]byte {
	return nil
}

func (d *dummyRows) Conn() *pgx.Conn {
	return nil
}

func (d *dummyRow) Scan(_ ...any) error {
	return nil
}

func (d *dummyTx) Begin(_ context.Context) (pgx.Tx, error) {
	slog.Debug("db.dummyTx.Begin()")
	return &dummyTx{}, nil
}

func (d *dummyTx) Exec(_ context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	slog.Debug("db.dummyTx.Exec()", "sql", sql, "arguments", arguments)
	return pgconn.NewCommandTag(dummyCommandTag), nil
}

func (d *dummyTx) Query(_ context.Context, sql string, args ...any) (pgx.Rows, error) {
	slog.Debug("db.dummyTx.Query()", "sql", sql, "args", args)
	return &dummyRows{}, nil
}

func (d *dummyTx) QueryRow(_ context.Context, sql string, args ...any) pgx.Row {
	slog.Debug("db.dummyTx.QueryRow()", "sql", sql, "args", args)
	return &dummyRow{}
}

func (d *dummyTx) Commit(_ context.Context) error {
	slog.Debug("db.dummyTx.Commit()")
	return nil
}

func (d *dummyTx) Rollback(_ context.Context) error {
	slog.Debug("db.dummyTx.Rollback()")
	return nil
}

func (d *dummyTx) CopyFrom(
	_ context.Context,
	tableName pgx.Identifier,
	columnNames []string,
	rowSrc pgx.CopyFromSource,
) (int64, error) {
	slog.Debug("db.dummyTx.CopyFrom()", "tableName", tableName, "columnNames", columnNames, "rowSrc", rowSrc)
	return 0, nil
}

func (d *dummyTx) SendBatch(_ context.Context, b *Batch) pgx.BatchResults {
	slog.Debug("db.dummyTx.SendBatch()", "batch", b)
	return &dummyBatchResults{}
}

func (d *dummyTx) LargeObjects() pgx.LargeObjects {
	slog.Debug("db.dummyTx.LargeObjects()")
	return pgx.LargeObjects{}
}

func (d *dummyTx) Prepare(_ context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	slog.Debug("db.dummyTx.Prepare()", "name", name, "sql", sql)
	return &pgconn.StatementDescription{}, nil
}

func (d *dummyTx) Conn() *pgx.Conn {
	return nil
}
