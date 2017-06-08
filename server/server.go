/************************************************************/
/********      Two-Factor Authentication Server      ********/
/********            By Carl Amko                    ********/
/************************************************************/
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strconv"
	"sync"
	"time"
	mgo "gopkg.in/mgo.v2"

	"github.com/rs/cors"
	"github.com/gorilla/mux"
	"github.com/BurntSushi/toml"
	"encoding/json"
	"io"
	"bytes"
	"gopkg.in/mgo.v2/bson"
)

// Generic mutex for locking map write access.
var mutex = &sync.Mutex{}

// Maps clientID to their data.
var clientMap = make(map[string]*ClientData)

// Database session for data manipulation.
var database * mgo.Database

// Config reference
var config Config

// DBconfig reference
var DBconfig DBConfig

// Configuration struct for TOML parsing.
type Config struct {
	TimeoutSeconds int
	NumKeyDigits   int
}

// DB configuration struct.
type DBConfig struct {
	DBConnectionURI string
}

// Struct for containing client data for ease of access.
type ClientData struct {
	TimeRemaining time.Duration
	Key           string
}

type User struct {
	Email string
	Serial string
	DeviceSerial string
}

const usableCharacters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRand() (r * rand.Rand) {
	// Create a new, randomly seeded rand.
	// See "Example (Rand)" at https://golang.org/pkg/math/rand/
	//   for generating a uniquely random seed based on time.
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return
}

func generateKey() (key string) {
	// Create a random object.
	r:= generateRand()

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

func checkSerialHandler(writer http.ResponseWriter, request *http.Request) {
	// Log incoming request.
	dump, _ := httputil.DumpRequest(request, true)
	fmt.Println(string(dump))

}

func newUserHandler(writer http.ResponseWriter, request *http.Request) {
	// Log incoming request.
	dump, _ := httputil.DumpRequest(request, true)
	fmt.Println(string(dump))

	var newUser User
	// Attempt to decode request body into User interface.
	err := json.NewDecoder(request.Body).Decode(&newUser)
	switch {
	case err == io.EOF:
		http.Error(writer, "No body provided to /user/new route!", 400)
		return
	case err!= nil:
		http.Error(writer, err.Error(), 400)
		panic(err)
	}

	userCollection := database.C("users")

	// User already exists, ignore insert.
	if count, err := userCollection.Find(bson.M{"email" : newUser.Email}).Count(); count > 0 || err != nil {
		errorMsg := fmt.Sprintf("User with email %s already exists!", newUser.Email)
		http.Error(writer, errorMsg, 400)
		fmt.Println(errorMsg)
		return
	}

	// Generate random serial for new user.
	r := generateRand()
	var buffer bytes.Buffer
	for i := 0; i < 12; i++ {
		buffer.WriteByte(usableCharacters[r.Intn(len(usableCharacters))])
	}

	newUser.Serial = buffer.String()
	fmt.Println(newUser)

	//err = userCollection.Insert(newUser)
	//if err != nil {
	//	println(err)
	//}

	buffer.Reset()
	json.NewEncoder(&buffer).Encode(newUser)
	writer.Write(buffer.Bytes())
}

func handler(writer http.ResponseWriter, request *http.Request) {
	// Check if client is already mapped.
	_, bClientExists := clientMap[request.RemoteAddr]
	fmt.Printf("Incoming request '%s' from %s. \n", request.Method, request.RemoteAddr)

	// Only start a new timer if client doesn't already exist.
	if !bClientExists {
		// Create new entry in client map.
		clientMap[request.RemoteAddr] = &ClientData{time.Duration(config.TimeoutSeconds) * time.Second, generateKey()}

		// Set timer for config.TimeoutSeconds seconds.
		countdown(request.RemoteAddr)
		fmt.Printf("Client %s has been assigned key %s for %d seconds!\n", request.RemoteAddr, clientMap[request.RemoteAddr].Key, config.TimeoutSeconds)
	}

	// Write the key back to the source.
	//writer.Write(*clientMap[request.RemoteAddr])
}

func main() {
	// Parse config TOML.
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Configuration parsed successfully:\n -Keys will generate with %d digits.\n -Keys will expire in %d seconds.\n", config.NumKeyDigits, config.TimeoutSeconds)

	// Parse DB config TOML.
	if _, err := toml.DecodeFile("DBconfig.toml", &DBconfig); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Database configuration parsed successfully:\n -Connection URI is %s.\n", DBconfig.DBConnectionURI)

	// Connect to Mongo database.
	dbSession, err := mgo.Dial(DBconfig.DBConnectionURI)
	defer dbSession.Close()

	database = dbSession.DB("tfa")
	if err != nil {
		panic(err)
	}

	// Create routes.
	r := mux.NewRouter()
	r.HandleFunc("/user/new", newUserHandler).Methods("POST")
	r.HandleFunc("/validate/serial", checkSerialHandler).Methods("POST")
	r.HandleFunc("/", handler).Methods("GET")

	// Create access control handler.
	handler := cors.Default().Handler(r)
	http.Handle("/", r)

	// Start server.
	fmt.Println(http.ListenAndServe(":8080", handler))
}
