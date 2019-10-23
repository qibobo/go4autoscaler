package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func defaultHander(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request.")
	target := os.Getenv("TARGET")
	if target == "" {
		target = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", target)
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
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	}()
}
func sleepHandler(w http.ResponseWriter, r *http.Request) {
	seconds, err := parseIntParam(r, "seconds")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Fprint(w, sleep(seconds))
	}()
}
func findPrimeHandler(w http.ResponseWriter, r *http.Request) {
	max, err := parseIntParam(r, "max")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Fprint(w, prime(max))
	}()
}
func allocateMemoryHandler(w http.ResponseWriter, r *http.Request) {
	mb, err := parseIntParam(r, "mb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Fprint(w, bloat(mb))
	}()
}

func sleep(seconds int) string {
	start := time.Now().UnixNano()
	time.Sleep(time.Duration(seconds) * time.Second)
	end := time.Now().UnixNano()
	return fmt.Sprintf("Slept for %.2f secounds.\n", float64(end-start)/1000000000)
}
func allPrimes(N int) []int {

	var x, y, n int
	nsqrt := math.Sqrt(float64(N))

	is_prime := make([]bool, N)

	for x = 1; float64(x) <= nsqrt; x++ {
		for y = 1; float64(y) <= nsqrt; y++ {
			n = 4*(x*x) + y*y
			if n <= N && (n%12 == 1 || n%12 == 5) {
				is_prime[n] = !is_prime[n]
			}
			n = 3*(x*x) + y*y
			if n <= N && n%12 == 7 {
				is_prime[n] = !is_prime[n]
			}
			n = 3*(x*x) - y*y
			if x > y && n <= N && n%12 == 11 {
				is_prime[n] = !is_prime[n]
			}
		}
	}

	for n = 5; float64(n) <= nsqrt; n++ {
		if is_prime[n] {
			for y = n * n; y < N; y += n * n {
				is_prime[y] = false
			}
		}
	}

	is_prime[2] = true
	is_prime[3] = true

	primes := make([]int, 0, 1270606)
	for x = 0; x < len(is_prime)-1; x++ {
		if is_prime[x] {
			primes = append(primes, x)
		}
	}

	// primes is now a slice that contains all primes numbers up to N
	return primes
}

func bloat(mb int) string {
	b := make([]byte, mb*1024*1024)
	b[0] = 1
	b[len(b)-1] = 1
	return fmt.Sprintf("Allocated %v Mb of memory.\n", mb)
}

func prime(max int) string {
	p := allPrimes(max)
	if len(p) > 0 {
		return fmt.Sprintf("The largest prime less than %v is %v.\n", max, p[len(p)-1])
	} else {
		return fmt.Sprintf("There are no primes smaller than %v.\n", max)
	}
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
	http.HandleFunc("/prime", findPrimeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
