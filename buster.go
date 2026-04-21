package main

import (
	"fmt"
	"net/http"
	"sync"
	"os"
	"bufio"
)

// testDir esegue una richiesta HTTP verso un URL target per vedere se la directory esiste
func testDir(baseURL string, dir string, wg *sync.WaitGroup) {
	// defer wg.Done() segnala al WaitGroup che questa funzione ha finito. Serve per la concorrenza.
	defer wg.Done()

	// Costruiamo l'URL finale
	target := baseURL + "/" + dir

	// Effettuiamo una richiesta di tipo GET
	resp, err := http.Get(target)
	if err != nil {
		// Ignoriamo eventuali errori di rete (es. timeout) per non spamare la console
		return
	}
	// Molto importante in Go: chiudiamo sempre il Body della risposta!
	defer resp.Body.Close()

	// 200 OK significa che la risorsa esiste
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("[+] Trovato (200): %s\n", target)
	} else if resp.StatusCode == http.StatusForbidden {
		fmt.Printf("[!] Accesso negato (403): %s\n", target)
	}
}