package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var sampleSymbols = []string{"AACC", "GOOG", "AAPL", "MSFT", "TSLA", "AMZN", "FB", "BABA"}

func worker(ctx context.Context, wg *sync.WaitGroup, coll *mongo.Collection) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		symbol := sampleSymbols[rand.Intn(len(sampleSymbols))]
		filter := bson.M{"stock_symbol": symbol}

		_ = coll.FindOne(context.Background(), filter)
	}
}

func main() {
	uri := flag.String("uri", "", "MongoDB URI (например mongodb://user:pass@host:27017)")
	concurrency := flag.Int("concurrency", 10, "Количество одновременных запросов (горутин)")
	dbName := flag.String("db", "market", "Имя базы данных")
	collName := flag.String("coll", "quotes", "Имя коллекции")
	showHelp := flag.Bool("help", false, "Показать помощь")
	flag.Parse()

	if *showHelp || *uri == "" || *concurrency <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s -uri <mongo-uri> -concurrency <n> [-db name] [-coll name]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	clientOpts := options.Client().ApplyURI(*uri)
	clientOpts.SetMaxPoolSize(uint64(*concurrency) * 2)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mongo.Connect error: %v\n", err)
		cancel()
		os.Exit(1)
	}
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Fprintf(os.Stderr, "mongo.Ping error: %v\n", err)
		_ = client.Disconnect(ctx)
		cancel()
		os.Exit(1)
	}
	cancel()

	coll := client.Database(*dbName).Collection(*collName)
	fmt.Printf("Connected to %s, collection %s. Starting %d reader workers. Send SIGINT to stop.\n",
		*uri, *collName, *concurrency)

	workCtx, workCancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(workCtx, &wg, coll)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	<-sigCh

	fmt.Println("Received SIGINT — stopping...")
	workCancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		fmt.Println("Timeout waiting for workers; exiting.")
	}

	disconnectCtx, dcancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dcancel()
	_ = client.Disconnect(disconnectCtx)
	fmt.Println("Disconnected. Exit complete.")
}
