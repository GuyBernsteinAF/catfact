package internal

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// GetFact fetches a cat fact from the given URL
// Returns empty string on any error
func GetFact(url string) string {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	res, err := client.Get(url)
	if err != nil {
		return "" // Return empty string instead of panic
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "" // Return empty string for non-200 status codes
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "" // Return empty string for JSON decode errors
	}

	// Check if fact exists and is a string
	factInterface, ok := result["fact"]
	if !ok {
		return "" // Return empty string if "fact" field doesn't exist
	}

	fact, ok := factInterface.(string)
	if !ok {
		return "" // Return empty string if fact is not a string
	}

	return fact
}

func PhaseThree() {
	//phase three
	wg := sync.WaitGroup{}

	ch := make(chan string, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			fact := GetFact("https://catfact.ninja/fact")
			if fact != "" { // Only add non-empty facts
				ch <- strconv.Itoa(index+1) + ". " + fact
			}
		}(i)
	}
	wg.Wait()
	close(ch)
	for fact := range ch {
		fmt.Println(fact)
	}
}

func PhaseTwo() {
	// phase two
	for i := 0; i < 5; i++ {
		fact := GetFact("https://catfact.ninja/fact")
		if fact != "" { // Only print non-empty facts
			fmt.Printf("%d. %s\n", i+1, fact)
		}
	}
	fmt.Println()
}

func PhaseOne() {
	// phase one
	fact := GetFact("https://catfact.ninja/fact")
	if fact != "" {
		fmt.Printf("%s\n", fact)
	}
	fmt.Println()
}

func PhaseFour(amount int) []string {
	ch := make(chan string, amount)
	wg := sync.WaitGroup{}

	for i := 0; i < amount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fact := GetFact("https://catfact.ninja/fact")
			if fact != "" { // Only send non-empty facts
				ch <- fact
			}
		}()
	}

	// Close channel after all goroutines complete
	go func() {
		wg.Wait()
		close(ch)
	}()

	var result []string
	for fact := range ch {
		result = append(result, fact)
	}

	return result
}

// PhaseFourWithURL allows testing with a custom URL
func PhaseFourWithURL(amount int, url string) []string {
	ch := make(chan string, amount)
	wg := sync.WaitGroup{}

	for i := 0; i < amount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fact := GetFact(url)
			if fact != "" {
				ch <- fact
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var result []string
	for fact := range ch {
		result = append(result, fact)
	}

	return result
}
