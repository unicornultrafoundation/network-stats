package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/MariusVanDerWijden/node-crawler-backend/api"
	"github.com/MariusVanDerWijden/node-crawler-backend/input"
)

var (
	crawlerDBPath = flag.String("crawler-db-path", "operadb.sqlite", "Crawler Database SQLite Path")
	apiDBPath     = flag.String("api-db-path", "apidb.sqlite", "API Database SQLite Path")
	dropNodesTime = flag.Duration("drop-time", 24*time.Hour, "Time to drop crawled nodes")
)

func main() {
	flag.Parse()

	crawlerDB, err := sql.Open("sqlite3", *crawlerDBPath)
	if err != nil {
		panic(err)
	}
	shouldInit := false
	if _, err := os.Stat(*apiDBPath); os.IsNotExist(err) {
		shouldInit = true
	}
	nodeDB, err := sql.Open("sqlite3", *apiDBPath)
	if err != nil {
		panic(err)
	}
	if shouldInit {
		fmt.Println("DB did not exist, init")
		if err := createDB(nodeDB); err != nil {
			panic(err)
		}
	}
	var wg sync.WaitGroup
	wg.Add(3)
	// Start reading deamon
	go newNodeDeamon(&wg, crawlerDB, nodeDB)
	go dropDeamon(&wg, nodeDB)
	// Start the API deamon
	apiDeamon := api.New(nodeDB)
	go apiDeamon.HandleRequests(&wg)
	wg.Wait()
}

// newNodeDeamon reads new nodes from the crawler and puts them in the db
// Might trigger the invalidation of caches for the api in the future
func newNodeDeamon(wg *sync.WaitGroup, crawlerDB, nodeDB *sql.DB) {
	defer wg.Done()
	lastCheck := time.Time{}
	for {
		nodes, err := input.ReadRecentNodes(crawlerDB, lastCheck)
		if err != nil {
			fmt.Printf("Error reading nodes: %v\n", err)
			return
		}
		lastCheck = time.Now()
		if len(nodes) > 0 {
			err := InsertCrawledNodes(nodeDB, nodes)
			if err != nil {
				fmt.Printf("Error inserting nodes: %v\n", err)
			}
			fmt.Printf("%d nodes inserted\n", len(nodes))
		}
		time.Sleep(time.Second)
	}
}

func dropDeamon(wg *sync.WaitGroup, db *sql.DB) {
	defer wg.Done()
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		<-ticker.C
		err := dropOldNodes(db, *dropNodesTime)
		if err != nil {
			panic(err)
		}
	}
}
