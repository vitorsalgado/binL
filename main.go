package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	conf := loadConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var db *sql.DB
	limit := 5
	backoff := 3 * time.Second
	factor := 2
	for i := 0; i < limit; i++ {
		if i > 0 {
			time.Sleep(backoff)
			backoff = time.Duration(factor) * backoff
		}
		var err error
		db, err = sql.Open("mysql", conf.ConnString)
		if err != nil {
			logger.Error("error opening database", slog.String("error", err.Error()))
			os.Exit(1)
			break
		}

		ctxPing, cancelPing := context.WithTimeout(ctx, 1*time.Second)
		defer cancelPing()
		err = db.PingContext(ctxPing)
		if err != nil {
			slog.Warn("error pinging the database", slog.String("error", err.Error()))
			continue
		}

		break
	}

	if db == nil {
		logger.Error("unable to connect to the database")
		os.Exit(1)
		return
	}

	canalConf := canal.NewDefaultConfig()
	canalConf.Dump = canal.DumpConfig{ExecutionPath: ""}
	canalConf.Addr = conf.CanalAddr
	canalConf.User = conf.CanalUser
	canalConf.Password = conf.CanalPassword
	canalConf.ServerID = 1
	canalConf.Flavor = "mysql"
	canalConf.IncludeTableRegex = conf.CanalTableRegex

	logger.Info(fmt.Sprintf("connecting to the database on: %s", conf.ConnString))

	can, err := canal.NewCanal(canalConf)
	if err != nil {
		logger.Error("error creating the canal", slog.String("error", err.Error()))
		os.Exit(1)
		return
	}

	can.SetEventHandler(&RowOnlyEventHandler{logger})

	go func() {
		<-ctx.Done()
		defer cancel()

		can.Close()

		err = db.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	err = can.Run()
	if err != nil {
		logger.Error("error after running the canal", slog.String("error", err.Error()))
		os.Exit(1)
		return
	}
}
