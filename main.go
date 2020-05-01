package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/greytabby/grawl/crawler"
	"github.com/greytabby/grawl/storage"
)

var (
	site           string
	parallelism    int
	allowedHosts   string
	depth          int
	headlessChrome bool
	outputDir      string
	logger         = log.New(os.Stdout, "Grawl ", log.LstdFlags)
)

func main() {
	flag.StringVar(&site, "site", "", "Site to crawl")
	flag.StringVar(&outputDir, "output_dir", "", "Directory name for saving crawl result")
	flag.IntVar(&parallelism, "parallelism", 5, "Number of parallel execution of crawler")
	flag.IntVar(&depth, "depth", 1, "Limit number of follow links on crawling")
	flag.BoolVar(&headlessChrome, "headless_chrome", false, "Use headless chrome on crawling")
	flag.StringVar(&allowedHosts, "allowed_hosts", "", "Accessibel hosts. Use comma to specify multiple hosts")

	// Load argument from environment variables.
	flag.VisitAll(func(f *flag.Flag) {
		if s := os.Getenv(strings.ToUpper(f.Name)); s != "" {
			f.Value.Set(s)
		}
	})

	flag.Parse()
	err := run()
	if err != nil {
		os.Exit(1)
	}
}

func run() error {
	// Initial settings
	storage := storage.NewFileStorage(outputDir)
	lr := crawler.NewLimitRule()
	if allowedHosts != "" {
		ah := strings.Split(allowedHosts, ",")
		lr.AddAllowedHosts(ah...)
	}
	c := crawler.NewCrawlerWithLimitRule(site, depth, lr)
	if headlessChrome {
		c.UseHeadlessChrome()
	}
	c.SetParallelism(parallelism)

	c.OnVisited(func(cr *crawler.CrawlResult) {
		logger.Printf("Visited: %s", cr.URL.String())
		storage.Save(cr)
	})
	c.OnError(func(err error) {
		logger.Println(err)
	})

	logger.Printf("Output base directory: %s", outputDir)
	logger.Printf("Crawling site: %v", site)
	logger.Printf("Crawling max depth: %v", depth)
	logger.Println("Start Crawling...")
	c.Crawl()
	return nil
}
