package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(linkChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range linkChan {
		time.Sleep(1 * time.Second)
		fmt.Printf("Done processing link #%s\n", url)
	}

}

func main() {
	t := time.Now()
	links := make([]string, 50)
	for i := 0; i < 50; i++ {
		links[i] = fmt.Sprintf("%d", i+1)
	}

	linkChan := make(chan string, 10) // HL123
	wg := new(sync.WaitGroup)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(linkChan, wg)
	}

	for _, link := range links {
		linkChan <- link
	}

	close(linkChan)
	fmt.Println("Waiting")
	wg.Wait()
	fmt.Println("Total time", time.Since(t))
}
