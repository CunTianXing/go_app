package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//stopChan := make(chan struct{})
	wg := sync.WaitGroup{}
	//tickStoppedChan := make(chan struct{})
	wg.Add(1)
	go tick(ctx, &wg)
	//tockStoppedChan := make(chan struct{})
	wg.Add(1)
	go tock(ctx, &wg)

	wg.Add(1)
	go server(ctx, &wg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("main: received C-c - shutting down")
	fmt.Println("main: telling goroutines to stop")
	cancel()
	//<-tickStoppedChan
	//<-tockStoppedChan
	wg.Wait()
	fmt.Println("main: all goroutines have told us they've finished")
}

func tick(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tick: tick %s\n", now.UTC().Format("20060102-150405.000000000"))
		case <-ctx.Done():
			fmt.Println("tick: caller has told us to stop")
			return
		}
	}
}

func tock(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("tock: tick %s\n", now.UTC().Format("20060102-150405.000000000"))
		case <-ctx.Done():
			fmt.Println("tock: caller has told us to stop")
			return
		}
	}
}

func server(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("server: received request")
		time.Sleep(3 * time.Second)
		io.WriteString(w, "Finished!\n")
		fmt.Println("server: request finished")
	}))

	srv := &http.Server{Addr: ":8080", Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("Listen : %s\n", err)
		}
	}()
	<-ctx.Done()
	fmt.Println("server: caller has told us to stop")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
	fmt.Println("server gracefully stopped")
}
