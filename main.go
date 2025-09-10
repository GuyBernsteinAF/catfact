package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

func getFact(url string) string {
	res, err := http.Get(url)
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

func main() {
	cadetsTask()
}

func cadetsTask() {
	fmt.Println("cadets task started")

	var phase int
	fmt.Print("Select a phase (1, 2, or 3): ")
	fmt.Scan(&phase)

	for {
		switch phase {
		case 1:
			fmt.Println()
			fmt.Println("Phase One:")
			phaseOne()
			break
		case 2:
			fmt.Println()
			fmt.Println("Phase Two:")
			phaseTwo()
			break
		case 3:
			fmt.Println()
			fmt.Println("Phase Three:")
			phaseThree()
			fmt.Println()
			break
		default:
			fmt.Println()
			fmt.Println("Invalid phase selected. Please choose 1, 2, or 3.")
			fmt.Print("Select a phase (1, 2, or 3): ")
			fmt.Scan(&phase)
			continue
		}
		break // Exit the loop when a valid phase is executed
	}

	fmt.Println("cadets task finished")
}

func phaseThree() {
	//phase three
	wg := sync.WaitGroup{}

	ch := make(chan string, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fact := getFact("https://catfact.ninja/fact")
			ch <- "" + strconv.Itoa(i+1) + ". " + fact
		}()
	}
	wg.Wait()
	close(ch)
	for i := range ch {
		fmt.Println(i)
	}

}

func phaseTwo() {
	// phase two
	for i := 0; i < 5; i++ {
		fact := getFact("https://catfact.ninja/fact")
		fmt.Printf("%d. %s\n", i+1, fact)
	}
	fmt.Println()
}

func phaseOne() {
	// phase one
	fact := getFact("https://catfact.ninja/fact")
	fmt.Printf("%s\n", fact)
	fmt.Println()
}
