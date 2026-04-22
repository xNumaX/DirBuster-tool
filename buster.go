package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

// testDir executes an HTTP request to a target URL to check if the directory exists
func testDir(baseURL string, dir string, wg *sync.WaitGroup) {
	// defer wg.Done() signals to the WaitGroup that this function has finished
	defer wg.Done()

	// Build the final URL
	target := baseURL + "/" + dir

	// Perform a GET request
	resp, err := http.Get(target)
	if err != nil {
		// Stampiamo esplicitamente l'errore per diagnosticare se la connessione fallisce
		fmt.Printf("[-] Error requesting %s: %v\n", target, err)
		return
	}
	// Always close the response Body to free network resources!
	defer resp.Body.Close()

	// 200 OK means the resource exists
	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Printf("[+] Found (200): %s\n", target)
	case http.StatusForbidden:
		fmt.Printf("[!] Access denied (403): %s\n", target)
	default:
		// Stampa tutte le altre risposte (come 404, 301, 502) per non ignorarle in silenzio
		fmt.Printf("[-] Ignored (%d): %s\n", resp.StatusCode, target)
	}
}

func main() {
	// Configure the "-w" flag for the wordlist
	wordlistPath := flag.String("w", "wordlist.txt", "Path to the wordlist file to use")

	// Parse the command line flags
	flag.Parse()

	// Check positional arguments (the URL)
	if len(flag.Args()) < 1 {
		fmt.Println("Error: Incorrect syntax.")
		fmt.Println("Usage: go run buster.go [-w wordlist.txt] <URL>")
		fmt.Println("Example 1: go run buster.go http://localhost")
		fmt.Println("Example 2: go run buster.go -w custom_list.txt http://localhost")
		return
	}
	baseURL := flag.Args()[0]

	// Validate the provided URL
	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		fmt.Println("Error: The provided argument does not appear to be a valid URL.")
		fmt.Println("Make sure to include http:// or https:// (e.g. http://example.com)")
		return
	}

	// Open the wordlist file
	file, err := os.Open(*wordlistPath)
	if err != nil {
		fmt.Printf("Error opening wordlist '%s': %v\n", *wordlistPath, err)
		return
	}
	defer file.Close()

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)

	fmt.Printf("Starting DirBuster attack on: %s\n", baseURL)
	fmt.Printf("Using wordlist: %s\n", *wordlistPath)
	fmt.Println("------------------------------------------------")

	wordCount := 0
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		wordCount++

		wg.Add(1)
		go testDir(baseURL, word, &wg)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading wordlist: %v\n", err)
	}
	fmt.Printf("Loaded %d words from wordlist.\n", wordCount)

	wg.Wait()

	fmt.Println("------------------------------------------------")
	fmt.Println("Scan completed successfully!")
}
