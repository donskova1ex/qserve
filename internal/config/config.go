package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultHost     = "localhost"
	DefaultUser     = "root"
	DefaultPassword = ""
	DefaultPort     = 8080

	DefaultPostgresPort = "5432"
	DefaultMySQLPort    = "3306"
	DefaultMSSQLPort    = "1433"
	DefaultSQLitePort   = "5432"
)

const (
	DBTypePostgres = "postgres"
	DBTypeMySQL    = "mysql"
	DBTypeSQLite   = "sqlite"
	DBTypeMSSQL    = "mssql"
)

type ConfigReader interface {
	ReadConfig() (*Config, error)
}

type ConfigValidator interface {
	Validate(cfg *Config) error
}

type Config struct {
	DBType string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	Port   int
}

func (c *Config) Validate() error {
	if c.DBType == "" {
		return errors.New("database type is required")
	}

	switch c.DBType {
	case DBTypePostgres, DBTypeMySQL, DBTypeSQLite, DBTypeMSSQL:
		// OK
	default:
		return fmt.Errorf("unsupported database type: %s", c.DBType)
	}

	if c.DBHost == "" {
		return errors.New("database host is required")
	}

	if c.DBPort == "" {
		return errors.New("database port is required")
	}

	if port, err := strconv.Atoi(c.DBPort); err != nil || port <= 0 || port > 65535 {
		return fmt.Errorf("invalid database port: %s. Port must be between 1 and 65535", c.DBPort)
	}

	if c.DBUser == "" {
		return errors.New("database user is required")
	}

	if c.DBName == "" {
		return errors.New("database name is required")
	}

	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid service port: %d. Port must be between 1 and 65535", c.Port)
	}

	return nil
}

type InteractiveConfigReader struct {
	reader *bufio.Reader
}

func NewInteractiveConfigReader() *InteractiveConfigReader {
	return &InteractiveConfigReader{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (r *InteractiveConfigReader) readInput(prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := r.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(input), nil
}

func getDefaultPort(dbType string) string {
	switch dbType {
	case DBTypePostgres:
		return DefaultPostgresPort
	case DBTypeMySQL:
		return DefaultMySQLPort
	case DBTypeMSSQL:
		return DefaultMSSQLPort
	case DBTypeSQLite:
		return DefaultSQLitePort
	default:
		return DefaultPostgresPort
	}
}

func parseDBType(input string) (string, error) {
	switch input {
	case "1":
		return DBTypePostgres, nil
	case "2":
		return DBTypeMySQL, nil
	case "3":
		return DBTypeSQLite, nil
	case "4":
		return DBTypeMSSQL, nil
	default:
		return "", errors.New("invalid DB type selection. Please choose 1, 2, 3, or 4")
	}
}

func (r *InteractiveConfigReader) ReadConfig() (*Config, error) {
	cfg := &Config{}

	dbTypeInput, err := r.readInput("Select DB type (1: PostgreSQL, 2: MySQL, 3: SQLite, 4: MSSQL): ")
	if err != nil {
		return nil, fmt.Errorf("failed to read DB type: %w", err)
	}

	cfg.DBType, err = parseDBType(dbTypeInput)
	if err != nil {
		return nil, err
	}

	dbHost, err := r.readInput("Select DB host [localhost]: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read DB host: %w", err)
	}
	switch dbHost {
	case "":
		cfg.DBHost = DefaultHost
	default:
		cfg.DBHost = dbHost
	}

	dbPort, err := r.readInput("Select DB port: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read DB port: %w", err)
	}

	if dbPort == "" {
		cfg.DBPort = getDefaultPort(cfg.DBType)
	} else {
		if port, err := strconv.Atoi(dbPort); err != nil || port <= 0 || port > 65535 {
			return nil, fmt.Errorf("invalid port number: %s. Port must be between 1 and 65535", dbPort)
		}
		cfg.DBPort = dbPort
	}

	dbUser, err := r.readInput("Select DB user: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read DB user: %w", err)
	}
	switch dbUser {
	case "":
		cfg.DBUser = DefaultUser
	default:
		cfg.DBUser = dbUser
	}

	dbPassword, err := r.readInput("Select DB password: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read DB password: %w", err)
	}
	cfg.DBPass = dbPassword

	dbName, err := r.readInput("Select DB name: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read DB name: %w", err)
	}
	cfg.DBName = dbName

	servicePort, err := r.readInput("Select service port [8080]: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read service port: %w", err)
	}

	if servicePort == "" {
		cfg.Port = DefaultPort
	} else {
		port, err := strconv.Atoi(servicePort)
		if err != nil {
			return nil, fmt.Errorf("invalid service port: %s. Port must be a valid number", servicePort)
		}
		cfg.Port = port
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

func RunNewSetupWizard() (*Config, error) {
	reader := NewInteractiveConfigReader()
	return reader.ReadConfig()
}

func NewConfigFromDefaults(dbType, dbName string) *Config {
	return &Config{
		DBType: dbType,
		DBHost: DefaultHost,
		DBPort: getDefaultPort(dbType),
		DBUser: DefaultUser,
		DBPass: DefaultPassword,
		DBName: dbName,
		Port:   DefaultPort,
	}
}
