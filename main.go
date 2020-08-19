package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "net/http/pprof"
)

var b1k []byte
var b5k [1024 * 5]byte
var b10k [1024 * 10]byte

func init() {
	for i := 0; i < 1024; i++ {
		b1k = append(b1k, 65)
	}
	for i := 0; i < 1024*5; i++ {
		b5k[i] = byte(65)
	}
	for i := 0; i < 1024*10; i++ {
		b10k[i] = byte(65)
	}
}

func defaultHander(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request.")
	target := os.Getenv("TARGET")
	if target == "" {
		target = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", target)
}
func largeResponseHandler(w http.ResponseWriter, r *http.Request) {
	kb, err := parseIntParam(r, "kb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := []byte{}
	for i := 1; i <= kb; i++ {
		result = append(result, b1k...)
	}
	fmt.Fprintf(w, "Hello %s!\n", result)
}
func sleepAndLargeResponseHandler(w http.ResponseWriter, r *http.Request) {
	kb, err := parseIntParam(r, "kb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := []byte{}
	for i := 1; i <= kb; i++ {
		result = append(result, b1k...)
	}

	// fmt.Fprintf(w, "Hello %s!\n", result)
	seconds, err := parseIntParam(r, "seconds")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%s, Hello %s\n", sleep(seconds), result)
}

func consumeCpuHandler(w http.ResponseWriter, r *http.Request) {
	count, err := parseIntParam(r, "count")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	arr := make([]int, count)
	for i := 0; i < count; i++ {
		arr[i] = rand.Int()
	}
	for i := 0; i < count; i++ {
		for j := 0; j < count-i-1; j++ {
			if arr[j] > arr[j+1] {
				tmp := arr[j]
				arr[j] = arr[j+1]
				arr[j+1] = tmp
			}
		}
	}
	fmt.Fprintf(w, "finished bubble sort for %d numbers\n", count)
}
func sleepHandler(w http.ResponseWriter, r *http.Request) {
	seconds, err := parseIntParam(r, "seconds")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, sleep(seconds))
}
func sleepMsHandler(w http.ResponseWriter, r *http.Request) {
	ms, err := parseIntParam(r, "ms")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, sleepMs(ms))
}

func allocateMemoryHandler(w http.ResponseWriter, r *http.Request) {
	mb, err := parseIntParam(r, "mb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, bloat(mb))
}

func sleep(seconds int) string {
	start := time.Now().UnixNano()
	time.Sleep(time.Duration(seconds) * time.Second)
	end := time.Now().UnixNano()
	return fmt.Sprintf("Slept for %.2f seconds.\n", float64(end-start)/1000000000)
}

func sleepMs(ms int) string {
	start := time.Now().UnixNano()
	time.Sleep(time.Duration(ms) * time.Millisecond)
	end := time.Now().UnixNano()
	return fmt.Sprintf("Slept for %.2f ms.\n", float64(end-start)/1000000)
}

func bloat(mb int) string {
	b := make([]byte, mb*1024*1024)
	b[0] = 1
	b[len(b)-1] = 1
	return fmt.Sprintf("Allocated %v Mb of memory.\n", mb)
}

func parseIntParam(r *http.Request, param string) (int, error) {
	if value := r.URL.Query().Get(param); value != "" {
		i, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		if i == 0 {
			return i, nil
		}
		return i, nil
	}
	return 0, nil
}

func main() {
	log.Print("Hello world sample started.")

	http.HandleFunc("/", defaultHander)
	http.HandleFunc("/cpu", consumeCpuHandler)
	http.HandleFunc("/memory", allocateMemoryHandler)
	http.HandleFunc("/sleep", sleepHandler)
	http.HandleFunc("/sleepms", sleepMsHandler)
	http.HandleFunc("/largeresponse", largeResponseHandler)
	http.HandleFunc("/sleepandlarge", sleepAndLargeResponseHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
