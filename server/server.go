/************************************************************/
/********      Two-Factor Authentication Server      ********/
/********            By Carl Amko                    ********/
/************************************************************/
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

// Generic mutex for locking map write access.
var mutex = &sync.Mutex{}

// Maps clientID to their data.
var clientMap = make(map[string]*ClientData)

// Config reference
var config Config

type Config struct {
	TimeoutSeconds int `toml: timeoutSeconds`
	NumKeyDigits   int `toml: numKeyDigits`
}

type ClientData struct {
	TimeRemaining time.Duration
	Key           string
}

func generateKey() (key string) {

	// Create a new, randomly seeded rand.
	// See "Example (Rand)" at https://golang.org/pkg/math/rand/
	//   for generating a uniquely random seed based on time.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate random numbers for the auth key.
	for i := 0; i < config.NumKeyDigits; i++ {
		key += strconv.Itoa(r.Intn(10))
	}
	return
}

func countdown(clientID string) {

	// Create new ticker to count down for this IP every second.
	ticker := time.NewTicker(time.Second)

	// Create a goroutine to update time remaining.
	go func(clientID string) {

		// Called every second.
		for range ticker.C {
			// Maps aren't thread-safe. Need to lock to ensure no write corruption.
			mutex.Lock()
			clientMap[clientID].TimeRemaining -= time.Second
			mutex.Unlock()

			// Time has expired.
			if clientMap[clientID].TimeRemaining <= time.Duration(0) {
				// Stop ticker.
				ticker.Stop()
				// Remove remote IP entry from maps.
				mutex.Lock()
				delete(clientMap, clientID)
				mutex.Unlock()

				fmt.Printf("Client %s key has expired.\n", clientID)
			} else {
				// Log for testing.
				fmt.Printf("Client %s has %f seconds remaining.\n", clientID, clientMap[clientID].TimeRemaining.Seconds())
			}
		}
	}(clientID)

}

func handler(w http.ResponseWriter, r *http.Request) {

	// Check if client is already mapped.
	_, bClientExists := clientMap[r.RemoteAddr]
	fmt.Printf("Incoming request '%s' from %s. \n", r.Method, r.RemoteAddr)

	// Only start a new timer if client doesn't already exist.
	if !bClientExists {
		// Create new entry in client map.
		clientMap[r.RemoteAddr] = &ClientData{time.Duration(config.TimeoutSeconds) * time.Second, generateKey()}

		// Set timer for config.TimeoutSeconds seconds.
		countdown(r.RemoteAddr)
		fmt.Printf("Client %s has been assigned key %s for %d seconds!\n", r.RemoteAddr, clientMap[r.RemoteAddr].Key, config.TimeoutSeconds)
	}

	// Write the key back to the source.
	//w.Write(*clientMap[r.RemoteAddr])
}

func main() {

	// Parse TOML config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Configuration parsed successfully:\n -Keys will generate with %d digits.\n -Keys will expire in %d seconds.\n", config.NumKeyDigits, config.TimeoutSeconds)

	// Start server.
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
