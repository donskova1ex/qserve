package database

import (
	"context"
	"database/sql"
	"fmt"
	"qserve/internal/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
)

type DBConnector interface {
	Connect(ctx context.Context) error
	Close() error
	Ping(ctx context.Context) error
	GetDB() *sql.DB
	GetDriver() string
}

type QueryExecutor interface {
	ExecuteQuery(ctx context.Context, query string) ([]map[string]interface{}, error)
	ExecuteTransaction(ctx context.Context, query string) (int64, error)
}

type ConnectionManager struct {
	config *config.Config
	db     *sql.DB
	driver string
}

func NewConnectionManager(config *config.Config) *ConnectionManager {
	return &ConnectionManager{
		config: config,
	}
}

func (c *ConnectionManager) Connect(ctx context.Context) error {
	var dsn string
	var driver string

	switch c.config.DBType {
	case config.DBTypePostgres:
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.config.DBHost, c.config.DBPort, c.config.DBUser, c.config.DBPass, c.config.DBName)
		driver = "postgres"
	case config.DBTypeMySQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			c.config.DBUser, c.config.DBPass, c.config.DBHost, c.config.DBPort, c.config.DBName)
		driver = "mysql"
	case config.DBTypeSQLite:
		dsn = fmt.Sprintf("file:%s?cache=shared&mode=rwc", c.config.DBName)
		driver = "sqlite3"
	case config.DBTypeMSSQL:
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			c.config.DBUser, c.config.DBPass, c.config.DBHost, c.config.DBPort, c.config.DBName)
		driver = "sqlserver"
	default:
		return fmt.Errorf("unsupported database type: %s", c.config.DBType)
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	c.db = db
	c.driver = driver
	return nil
}

func (c *ConnectionManager) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *ConnectionManager) Ping(ctx context.Context) error {
	if c.db == nil {
		return fmt.Errorf("database not connected")
	}
	return c.db.PingContext(ctx)
}
