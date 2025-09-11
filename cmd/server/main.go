package main

import (
	"bytes"
	"catfacts/internal"
	"fmt"
	"io"
	"net/http"
	"os"
)

func captureStdout(f func()) string {
	// Save the original stdout
	originalStdout := os.Stdout

	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function in a goroutine to prevent blocking
	done := make(chan bool)
	var buf bytes.Buffer

	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	// Call the function that writes to stdout
	f()

	// Restore original stdout and close the writer
	w.Close()
	os.Stdout = originalStdout

	// Wait for the capture to complete
	<-done

	return buf.String()
}

func phaseOneAPI(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, captureStdout(internal.PhaseOne))
}

func phaseTwoAPI(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, captureStdout(internal.PhaseTwo))
}

func phaseThreeAPI(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, captureStdout(internal.PhaseThree))
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	http.HandleFunc("/phase-one", phaseOneAPI)
	http.HandleFunc("/phase-two", phaseTwoAPI)
	http.HandleFunc("/phase-three", phaseThreeAPI)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8090", nil)
}
