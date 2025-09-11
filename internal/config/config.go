package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DBType string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	Port   int
}

func RunNewSetupWizard() (*Config, error) {
	cfg := &Config{}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Select DB type (1: PostgreSQL, 2: MySQL, 3: SQLite, 4:MSSql ): ")
	dbType, _ := reader.ReadString('\n')
	dbType = strings.TrimSpace(dbType)
	switch dbType {
	case "1":
		cfg.DBType = "postgres"
	case "2":
		cfg.DBType = "mysql"
	case "3":
		cfg.DBType = "sqlite"
	case "4":
		cfg.DBType = "mssql"
	}

	fmt.Print("Select DB host [localhost]: ")
	dbHost, _ := reader.ReadString('\n')
	dbHost = strings.TrimSpace(dbHost)
	switch dbHost {
	case "":
		cfg.DBHost = "localhost"
	default:
		cfg.DBHost = dbHost
	}

	fmt.Print("Select DB port: ")
	dbPort, _ := reader.ReadString('\n')
	dbPort = strings.TrimSpace(dbPort)
	switch {
	case dbPort == "":
		switch cfg.DBType {
		case "postgres":
			cfg.DBPort = "5432"
		case "mysql":
			cfg.DBPort = "3306"
		case "mssql":
			cfg.DBPort = "1433"
		case "sqlite":
			cfg.DBPort = "5432"
		}
	default:
		cfg.DBPort = dbPort
	}

	fmt.Print("Select DB user: ")
	dbUser, _ := reader.ReadString('\n')
	dbUser = strings.TrimSpace(dbUser)
	switch {
	case dbUser == "":
		dbUser = "root"
	default:
		cfg.DBUser = dbUser
	}

	fmt.Print("Select DB password: ")
	dbPassword, _ := reader.ReadString('\n')
	dbPassword = strings.TrimSpace(dbPassword)
	switch {
	case dbPassword == "":
		dbPassword = "root"
	default:
		cfg.DBPass = dbPassword
	}

	fmt.Print("Select DB name: ")
	dbName, _ := reader.ReadString('\n')
	dbName = strings.TrimSpace(dbName)
	cfg.DBName = dbName

	fmt.Print("Select service port: ")
	servicePort, _ := reader.ReadString('\n')
	servicePort = strings.TrimSpace(servicePort)
	if servicePort == "" {
		servicePort = "8080"
	}
	cfg.Port, _ = strconv.Atoi(servicePort)
	return cfg, nil
}
