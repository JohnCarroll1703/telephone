package main

import (
	"github.com/getsentry/sentry-go"
	"sync"
	"telephone/internal/app"
	"telephone/internal/config"
	"time"
)

func main() {
	defer func() {
		err := recover()
		hub := sentry.CurrentHub()
		if err != nil {
			hub.Recover(err)
			sentry.Flush(time.Second * 5)
		}

		if err != nil {
			panic(err)
		}
	}()

	once := sync.Once{}

	cfg, err := config.GetConfig(&once)
	if err != nil {
		panic(err)
	}

	app.Run(cfg)
}
