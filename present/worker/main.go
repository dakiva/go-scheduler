package main

import (
	"fmt"
	"sync"
	"time"
)

func funcExecutor(funcChan chan func(), wg *sync.WaitGroup) {
	defer wg.Done()

	for f := range funcChan {
		f()
	}
}

func main() {
	t := time.Now()
	funcs := make([]func(), 50)
	for i := 0; i < 50; i++ {
		cp := i
		funcs[i] = func() {
			time.Sleep(1 * time.Second)
			fmt.Println("Done processing function", cp+1)
		}
	}

	funcChan := make(chan func(), 10)
	wg := new(sync.WaitGroup)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go funcExecutor(funcChan, wg)
	}

	for _, f := range funcs {
		funcChan <- f
	}

	close(funcChan)
	fmt.Println("Waiting")
	wg.Wait()
	fmt.Println("Total time", time.Since(t))
}
