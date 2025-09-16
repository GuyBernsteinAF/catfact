package internal

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

func GetFact(url string) string {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic(fmt.Sprintf("bad status code: %d", res.StatusCode))
	}
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		panic(err)
	}

	return result["fact"].(string)
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
			ch <- strconv.Itoa(index+1) + ". " + fact
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
		fmt.Printf("%d. %s\n", i+1, fact)
	}
	fmt.Println()
}

func PhaseOne() {
	// phase one
	fact := GetFact("https://catfact.ninja/fact")
	fmt.Printf("%s\n", fact)
	fmt.Println()
}

func PhaseFour(amount int) []string {
	ch := make(chan string, amount)
	for i := 0; i < amount; i++ {
		go func() {
			fact := GetFact("https://catfact.ninja/fact")
			ch <- fact
		}()
	}
	var result []string
	for i := 0; i < amount; i++ {
		result = append(result, <-ch)
	}
	return result
}
