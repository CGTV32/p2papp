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
)

func main() {

	var dataBytes []byte = fileToBytes("/home/cgtv32/Documents/p2ptest/test.txt")

	var wg sync.WaitGroup

	peerIP := "127.0.0.1:5643"

	listener, err := net.Listen("tcp", ":5643")
	if err != nil {
		log.Fatal("Error creating listener:", err)
	}
	defer listener.Close()

	fmt.Println("Listener initialised")

	conn, err := net.DialTimeout("tcp", peerIP, 10*time.Second)
	if err != nil {
		log.Fatal("Error connecting to peer:", err)
	}

	fmt.Println("Connection initialised")

	wg.Add(3)

	go listen(listener, &wg)
	go writeToPeer(conn, &wg, dataBytes)

	wg.Wait()

}

func writeToPeer(conn net.Conn, wg *sync.WaitGroup, message []byte) {
	defer wg.Done()
	defer conn.Close()

	_, err := conn.Write(message)
	if err != nil {
		log.Println("Error writing to peer:", err)
		return
	}

	fmt.Println("Written bytes:", len(message))
}

func listen(listener net.Listener, wg *sync.WaitGroup) {
	defer wg.Done()

	// discover a UPnP or NAT-PMP gateway on the local network
	natMapper := nat.Any()

	extPort, err := natMapper.AddMapping("TCP", 30303, 30303, "my-custom-p2p-app", 20*time.Minute)
	if err != nil {
		fmt.Printf("Failed to map port (is UPnP disabled on router?): %v\n", err)
		return
	}

	fmt.Printf("Successfully opened port! External internet traffic on port %d will forward here.\n", extPort)

	// clean up and close the port when done
	defer natMapper.DeleteMapping("TCP", int(extPort), 30303)

	conn, err := listener.Accept()
	if err != nil {
		log.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Listener is listening")

	receivedData, err := io.ReadAll(conn)
	if err != nil {
		log.Println("Error reading data:", err)
		return
	}

	fmt.Println("Received bytes:", receivedData)
}

func fileToBytes(filePath string) []byte {

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	return content

}
