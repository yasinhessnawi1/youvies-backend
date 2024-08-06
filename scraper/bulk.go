package scraper

import (
	"log"
	"sync"
)

type BulkScraper struct {
	scrapers []Scraper
}

func NewBulkScraper(scrapers []Scraper) *BulkScraper {
	return &BulkScraper{scrapers: scrapers}
}

func (b *BulkScraper) ScrapeAll() error {
	log.Println("Starting bulk scraping...")
	var wg sync.WaitGroup
	var errors []error
	mu := sync.Mutex{}

	for _, scraper := range b.scrapers {

		wg.Add(1)
		go func(scr Scraper) {
			defer wg.Done()
			err := scr.Scrape()
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				log.Printf("Error scraping: %v", err)
			}
		}(scraper)
	}
	wg.Wait()

	if len(errors) > 0 {
		return errors[0]
	}

	log.Println("Bulk scraping completed.")
	return nil
}
