package cron

import (
	"math/rand"
	"time"

	"github.com/scalent.io/scalent-hrms/pkg/log"
)

// const (
// 	minCronInterval = 15 * time.Second
// 	maxCronInterval = 60 * time.Second
// )

func StartCronJobs() error {
	log.Info("********************************Starting cron jobs", "")

	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			log.Info("------------------Running leads verification cron----------------", "")

		}
	}()

	log.Info("Lead verification cron started with random interval", "")
	return nil
}

// func randomCronInterval() time.Duration {
// 	return minCronInterval + time.Duration(rand.Int63n(int64(maxCronInterval-minCronInterval)))
// }

func randomCronInterval(baseInterval time.Duration) time.Duration {
	if baseInterval <= 0 {
		baseInterval = 15 * time.Second
	}
	randomDelay := time.Duration(rand.Int63n(int64(baseInterval)))
	return baseInterval + randomDelay
}
