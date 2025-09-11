package main

import (
	"catfacts/internal"
	"fmt"
)

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
			internal.PhaseOne()
			break
		case 2:
			fmt.Println()
			fmt.Println("Phase Two:")
			internal.PhaseTwo()
			break
		case 3:
			fmt.Println()
			fmt.Println("Phase Three:")
			internal.PhaseThree()
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
