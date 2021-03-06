package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	stats "github.com/fukata/golang-stats-api-handler"
	"github.com/kaneshin/playground/go/concurrency-pattern/payload"
)

func init() {
	go func() {
		num := 0
		for {
			if n := runtime.NumGoroutine(); n != num {
				s := stats.GetStats()
				fmt.Printf("Stats\n")
				fmt.Printf(" |`-GoroutineNum: %v\n", s.GoroutineNum)
				fmt.Printf(" |`-Memory: Alloc %v, Mallocs %v, Frees %v\n",
					float64(s.MemoryAlloc)/1000.0, float64(s.MemoryMallocs)/1000.0, float64(s.MemoryFrees)/1000.0)
				fmt.Printf(" |`-Heap: Alloc %v, Sys %v, Idle %v\n",
					float64(s.HeapAlloc)/1000.0, float64(s.HeapSys)/1000.0, float64(s.HeapIdle)/1000.0)
				fmt.Printf("  `-GC: Num %v, PerSecond %v, Pause %v\n", s.GcNum, s.GcPerSecond, s.GcPause)
				num = n
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Read the body into a string for json decoding
		var collection payload.Collection
		if err := json.NewDecoder(r.Body).Decode(&collection); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		for _, instance := range collection.Instances {
			do(instance)
		}

		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":8080", nil)
}

func numMaxQueues() int {
	num, _ := strconv.Atoi(os.Getenv("MAX_QUEUES"))
	if num == 0 {
		num = 1
	}
	return num
}

func numMaxWorkers() int {
	num, _ := strconv.Atoi(os.Getenv("MAX_WORKERS"))
	if num == 0 {
		num = 1
	}
	return num
}
