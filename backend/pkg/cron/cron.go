package cron

import (
	"context"
	"fmt"
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

	if leadService == nil {
		err := fmt.Errorf("LeadService not initialized")
		log.Error(err.Error(), "")
		return err
	}

	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			log.Info("------------------Running leads verification cron----------------", "")
			log.Info("Checking for pending email verification", "")

			hasPending, errResp := leadService.HasPendingVerification(context.Background())
			if errResp != nil {
				log.Error(errResp.Error(), "")
				time.Sleep(15 * time.Second)
				continue
			}

			if !hasPending {
				log.Info("No pending emails to verify", "")
				time.Sleep(15 * time.Second)
				continue
			}

			log.Info("Running lead verification cron", "")
			// Verify one lead and get user's configured interval
			verificationInterval, errResp := VerifyPendingLeads()
			log.Info(fmt.Sprintf("Verification interval from user or db : %v", verificationInterval), "")
			if errResp != nil {
				log.Error(errResp.Error(), "")
				continue
			}

			// Add random delay based on user's configured interval
			delay := randomCronInterval(verificationInterval)

			log.Info(fmt.Sprintf("Next lead verification cron will run after: %v", delay), "")
			time.Sleep(delay)
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
