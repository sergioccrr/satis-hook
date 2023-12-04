package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"
)

type HealthHandler struct{}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK!\n")
}

type GitLabHandler struct {
	ch chan string
}

type GitLabRequest struct {
	Project struct {
		Name string
		Path string `json:"path_with_namespace"`
	}
}

func (h *GitLabHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var gr GitLabRequest
	err := json.NewDecoder(r.Body).Decode(&gr)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "Request added to queue\n")

	go func(msg string) {
		fmt.Printf("Sending [%s] to queue\n", msg)
		h.ch <- msg
	}(gr.Project.Path)
}

func processQueue(ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		pkg := <-ch
		fmt.Printf("Read [%s] from queue\n", pkg)

		cmd := exec.Command("satis", "build", "satis.json", "/var/www/html", pkg)
		out, err := cmd.Output()
		if err != nil {
			fmt.Printf("Could not run command: %s\n", err)
		}
		fmt.Printf("Output: %s\n", out)
	}
}

func main() {
	var wg sync.WaitGroup

	queue := make(chan string)

	wg.Add(1)
	go processQueue(queue, &wg)

	health := HealthHandler{}
	http.Handle("/health", &health)

	gitlab := GitLabHandler{ch: queue}
	http.Handle("/gitlab", &gitlab)

	fs := http.FileServer(http.Dir("/var/www/html"))
	http.Handle("/", fs)

	defer close(queue)

	addr := ":80"

	log.Printf("Listening on %s...\n", addr)
	err := http.ListenAndServe(addr, nil)

	if errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server closed\n")
	} else if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}

	wg.Wait()
}
