package scheduler

import (
	"context"
	"esp-organizer/internal/domain/infoin"
	"log"
	"time"
)

// CurationScheduler manages automatic sample curation
type CurationScheduler struct {
	interval time.Duration
	stopChan chan bool
	service  *infoin.SemanticLinkService
}

// NewCurationScheduler creates a new curation scheduler
func NewCurationScheduler(intervalHours int) (*CurationScheduler, error) {
	service, err := infoin.NewSemanticLinkService()
	if err != nil {
		return nil, err
	}

	return &CurationScheduler{
		interval: time.Duration(intervalHours) * time.Hour,
		stopChan: make(chan bool),
		service:  service,
	}, nil
}

// Start begins the curation scheduler
func (cs *CurationScheduler) Start() {
	log.Printf("🔄 Starting curation scheduler (interval: %v)", cs.interval)

	// Run immediately on start
	cs.runCuration()

	// Then run on schedule
	ticker := time.NewTicker(cs.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				cs.runCuration()
			case <-cs.stopChan:
				ticker.Stop()
				log.Println("🛑 Curation scheduler stopped")
				return
			}
		}
	}()
}

// Stop stops the curation scheduler
func (cs *CurationScheduler) Stop() {
	cs.stopChan <- true
}

// runCuration executes the curation process
func (cs *CurationScheduler) runCuration() {
	ctx := context.Background()
	log.Println("⏳ Running automatic sample curation...")

	startTime := time.Now()

	// Call the method on the service instance, not as a package function
	if err := cs.service.CurateHighQualityExtractions(ctx); err != nil {
		log.Printf("❌ Curation failed: %v", err)
		return
	}

	duration := time.Since(startTime)
	log.Printf("✅ Curation completed in %v", duration)
}

// RunOnce runs curation once without scheduling
func RunOnce() error {
	scheduler, err := NewCurationScheduler(24)
	if err != nil {
		return err
	}

	scheduler.runCuration()
	return nil
}
