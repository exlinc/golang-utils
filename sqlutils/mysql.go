package sqlutils

import "fmt"

// MySQLConnectionURL generates a connection URL for the built-in MySQL driver that uses the utf8mb4 charset and parses timestamps
func MySQLConnectionURL(hostname string, port int64, username string, password string, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", username, password, hostname, port, dbname)
}
