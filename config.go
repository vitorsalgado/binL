package main

import (
	"os"
	"strings"
)

const (
	DBConnString      = "DB_CONN"
	DBCanalAddr       = "DB_CANAL_ADDR"
	DBCanalUser       = "DB_CANAL_USER"
	DBCanalPassword   = "DB_CANAL_PASSWORD"
	DBCanalTableRegex = "DB_CANAL_TABLE_REGEX"
)

type Config struct {
	ConnString      string
	CanalAddr       string
	CanalUser       string
	CanalPassword   string
	CanalTableRegex []string
}

func loadConfig() Config {
	return Config{
		ConnString:      envStr(DBConnString, "root:@tcp(127.0.0.1:3306)/binl"),
		CanalAddr:       envStr(DBCanalAddr, "localhost:3306"),
		CanalUser:       envStr(DBCanalUser, "root"),
		CanalPassword:   envStr(DBCanalPassword, ""),
		CanalTableRegex: strings.Split(envStr(DBCanalTableRegex, "binl*"), ","),
	}
}

func envStr(n, def string) string {
	str := os.Getenv(n)
	if len(str) == 0 {
		return def
	}

	return str
}
