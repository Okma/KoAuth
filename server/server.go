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

// Generic keyLock for locking map write access.
var serverLock = &sync.Mutex{}

// Maps clientID to their data.
var connectionFrequencyMap = make(map[string]int)

// Database session for data manipulation.
var database * mgo.Database

var userCollection * mgo.Collection

// Config reference
var config Config

// DBconfig reference
var DBconfig DBConfig

var serverSeed string

var serialTimeRemaining int

// Configuration struct for TOML parsing.
type Config struct {
	ServerSerialRefreshSeconds int
	NumKeyDigits   int
	NumSerialDigits int
}

// DB configuration struct.
type DBConfig struct {
	DBConnectionURI string
}

// KeyData packages relevant key data for parsing and transmittal to mobile.
type KeyData struct {
	ServerTimeRemaining int
	Key           string
}

type User struct {
	Email string
	Serial string
	DeviceSerial string
}

const MAX_CONNECTIONS_PER_MIN = 5
const usableCharacters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRand() (r * rand.Rand) {
	// Create a new, randomly seeded rand.
	// See "Example (Rand)" at https://golang.org/pkg/math/rand/
	//   for generating a uniquely random seed based on time.
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return
}

func generateKey(userSerial string) (key string) {
	var serialSum = 0
	for i := 0; i < config.NumSerialDigits; i++ {
		// Modulus by the length of the serial is a safety precaution in case the server serial's length is
		//   changed, as this would cause a discrepancy between user and server serial lengths.
		serialSum += int(userSerial[i % len(userSerial)]) + int(serverSeed[i % len(serverSeed)])
	}

	var buffer bytes.Buffer
	for i := 0; i < config.NumKeyDigits; i++ {
		buffer.WriteString(string(serialSum / (i + 1) % 10))
	}

	return buffer.String()
}

func generateSerial() (serial string) {
	r := generateRand()
	var buffer bytes.Buffer
	for i := 0; i < config.NumSerialDigits; i++ {
		buffer.WriteByte(usableCharacters[r.Intn(len(usableCharacters))])
	}

	return buffer.String()
}

func startCountdown() {
	// Create new ticker to count down for this IP every second.
	ticker := time.NewTicker(time.Second)

	// Set current time remaining to 0 so that a seed is generated at startup.
	serialTimeRemaining = 0

	// Create a goroutine to update time remaining.
	go func() {
		// Update time remaining every second.
		for range ticker.C {
			serverLock.Lock()
			serialTimeRemaining -= 1
			serverLock.Unlock()

			// Refresh time has expired.
			if serialTimeRemaining <= 0 {
				serverLock.Lock()

				// Generate new sever serial.
				serverSeed = generateSerial()

				// Reset connection limiter map.
				connectionFrequencyMap = make(map[string]int)

				// Reset refresh time to max.
				serialTimeRemaining = config.ServerSerialRefreshSeconds
				serverLock.Unlock()

				fmt.Printf("New server serial generated: %s.\n", serverSeed)
			}
		}
	}()

}

func parseRequestBodyToUser(writer http.ResponseWriter, request *http.Request) (newUser *User){
	// Attempt to decode request body into User interface.
	err := json.NewDecoder(request.Body).Decode(&newUser)
	switch {
	case err == io.EOF:
		http.Error(writer, "No body provided!", 400)
		return
	case err!= nil:
		http.Error(writer, err.Error(), 400)
		panic(err)
	}
	return
}

func checkSerialHandler(writer http.ResponseWriter, request *http.Request) {
	// Log incoming request.
	dump, _ := httputil.DumpRequest(request, true)
	fmt.Println(string(dump))

	// Parse request body to user interface.
	user := parseRequestBodyToUser(writer, request)

	// Create query to find the user(s) with matching serial.
	query := userCollection.Find(bson.M{"serial" : user.Serial})

	// Check that there is at least one user with the given serial.
	if count, _ := query.Count(); count == 0 {
		errorMsg := fmt.Sprintf("User with serial %s does not exist!", user.Serial)
		http.Error(writer, errorMsg, 400)
		fmt.Println(errorMsg)
		return
	}

	// Update database entry with received device serial.
	userCollection.Update(bson.M{"serial" : user.Serial}, bson.M{"$set" : bson.M{"deviceserial": user.DeviceSerial}})

	// Read updated user.
	var currentUserEntry User = User{}
	err := query.One(&currentUserEntry)
	if err != nil {
		panic(err)
	}
	fmt.Println(currentUserEntry)

	// Write updated user back to client.
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(currentUserEntry)
	writer.Write(buffer.Bytes())
}

func newUserHandler(writer http.ResponseWriter, request *http.Request) {
	// Log incoming request.
	dump, _ := httputil.DumpRequest(request, true)
	fmt.Println(string(dump))

	// Parse request body to user interface.
	user := parseRequestBodyToUser(writer, request)

	// User already exists, ignore insert.
	if count, err := userCollection.Find(bson.M{"email" : user.Email}).Count(); count > 0 || err != nil {
		errorMsg := fmt.Sprintf("User with email %s already exists!", user.Email)
		http.Error(writer, errorMsg, 400)
		fmt.Println(errorMsg)
		return
	}

	// Generate random serial for new user.
	user.Serial = generateSerial()
	fmt.Println(*user)

	err := userCollection.Insert(user)
	if err != nil {
		println(err)
	}

	buffer := new(bytes.Buffer)
	json.NewEncoder(*buffer).Encode(user)
	writer.Write(buffer.Bytes())
}

func handler(writer http.ResponseWriter, request *http.Request) {
	// Check if client is already mapped.
	connections, bClientExists := connectionFrequencyMap[request.RemoteAddr]
	fmt.Printf("Incoming request '%s' from %s. Request %d out of %d.\n", request.Method, request.RemoteAddr, connections, MAX_CONNECTIONS_PER_MIN)

	// Only start a new timer if client doesn't already exist.
	if !bClientExists {
		// Create new entry in client map.
		connectionFrequencyMap[request.RemoteAddr] = 1
	} else {
		// Increment rate limit for IP.
		connectionFrequencyMap[request.RemoteAddr] += 1
	}

	// If not over limit, generate and send key.
	if connectionFrequencyMap[request.RemoteAddr] <= MAX_CONNECTIONS_PER_MIN {
		key := generateKey(request.Header.Get("serial"))
		fmt.Printf("Client %s has been assigned key %s.\n", request.RemoteAddr, key)

		// Encode and write KeyData back to requester.
		keyData := KeyData{serialTimeRemaining, key}
		buffer := bytes.Buffer{}
		json.NewEncoder(&buffer).Encode(keyData)
		writer.Write(buffer.Bytes())
	} else {
		http.Error(writer, "Too many requests received within allotted time limit.", 400)
	}
}

func main() {
	// Parse config TOML.
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Configuration parsed successfully:\n-Serial will generate with %d digits.\n-Keys will generate with %d digits.\n-Server serial will generate every %d seconds.\n",
		config.NumSerialDigits, config.NumKeyDigits, config.ServerSerialRefreshSeconds)

	// Parse DB config TOML.
	if _, err := toml.DecodeFile("DBconfig.toml", &DBconfig); err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Database configuration parsed successfully:\n -Connection URI is %s.\n", DBconfig.DBConnectionURI)

	// Connect to Mongo database.
	dbSession, err := mgo.Dial(DBconfig.DBConnectionURI)
	if err != nil {
		panic(err)
	}
	defer dbSession.Close()

	// Begin server time updater.
	startCountdown()

	// Assign database references.
	database = dbSession.DB("tfa")
	userCollection = database.C("users")

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
