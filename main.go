package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/p2p/nat"
	//"main.go/frontend"
)

func main() {
	var wg sync.WaitGroup

	//var dataBytes []byte = fileToBytes("/home/cgtv32/Documents/p2ptest/test.txt") // hard coded for tests, points to self for testing
	peerSocket := "127.0.0.1:30303"

	wg.Add(1)
	var writeMode string = os.Args[1]
	var message []string = os.Args[2:]
	var messageBytes []byte

	for _, msg := range message {
		messageBytes = append(messageBytes, []byte(msg)...) // ... lets append() ingest item by item since it expects singular elements and not a list
		messageBytes = append(messageBytes, byte(' '))
	}

	if writeMode == "false" {
		go listen(30303, &wg)
		wg.Wait()
	} else if writeMode == "true" {
		go writeToPeerFromFile(peerSocket, &wg, messageBytes)
		wg.Wait()
	}

}

func writeToPeerFromFile(peerSocket string, wg *sync.WaitGroup, message []byte) {
	defer wg.Done()

	conn, err := net.DialTimeout("tcp", peerSocket, 10*time.Second) // init a connection
	if err != nil {
		log.Fatal("Error connecting to peer:\n", err)
	}
	defer conn.Close() // make sure it closes gracefully

	fmt.Println("Connection initialised")

	_, err = conn.Write(message) // write a message to the conn
	if err != nil {
		log.Println("Error writing to peer:", err)
		return
	}

	fmt.Println("Written bytes:", len(message)) // print how many bytes it sent
}

func writeToPeerFromCLI(peerSocket string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.DialTimeout("tcp", peerSocket, 10*time.Second) // init a connection
	if err != nil {
		log.Fatal("Error connecting to peer:\n", err)
	}
	defer conn.Close() // make sure it closes gracefully

	fmt.Println("Connection initialised")

	for {

	}
}

func listen(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	var listenerPort string = ":" + fmt.Sprintf("%d", port) // cant use conventional typecasting cuz it turns into a chinese char
	fmt.Println(listenerPort)

	// discover a UPnP or NAT-PMP gateway on the local network
	natMapper := nat.Any()

	extPort, err := natMapper.AddMapping("TCP", port, port, "my-custom-p2p-app", 20*time.Minute) // open the UPnP port to make sure it can listen
	if err != nil {
		fmt.Printf("Failed to map port (is UPnP disabled on router or firewall rule disallowing?): %v\n", err)
		return
	}

	listener, err := net.Listen("tcp", listenerPort) // init listener
	if err != nil {
		log.Fatal("Error creating listener:", err)
	}
	defer listener.Close() // ensure it exits gracefully

	fmt.Println("Listener initialised")

	// clean up and close the port when done
	defer natMapper.DeleteMapping("TCP", int(extPort), port)

	for {
		func() { // wrap it in a func to stop memory leaks with repeated defer statements
			conn, err := listener.Accept() // blocks the script until it receives a message
			if err != nil {
				log.Println("Error accepting connection:", err)
				return
			}
			defer conn.Close()

			// use readall so it reads until the entire message is received
			receivedData, err := io.ReadAll(conn)
			if err != nil {
				log.Println("Error reading data:", err)
				return
			}

			fmt.Println("Received text:", string(receivedData)) // use conventional typecasting
		}()
	}
}

func fileToBytes(filePath string) []byte {

	content, err := os.ReadFile(filePath) // returns raw bytes :)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	return content

}
