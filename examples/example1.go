package main

import (
	"sync"

	log "github.com/Astera-org/easylog"
)

func main() {
	wg := new(sync.WaitGroup)

	if err := log.Init(
		log.SetLevel(log.DEBUG),
		log.SetFilePath("./"),
		log.SetFileName("test.log"),
		log.SetMaxSize(1),
	); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				log.Info("EXAMPLE: %d", j)
			}
		}()
	}
	wg.Wait()
}
