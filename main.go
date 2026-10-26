package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	peerIP := "127.0.0.1:5643"

	/*
		message := "hey"
		var payload []byte
		for _, char := range message {
			payload = append(payload, byte(char))
		}
	*/

	// bind first, so the port is open before we dial
	listener, err := net.Listen("tcp", ":5643")
	if err != nil {
		log.Fatal("Error creating listener obj:\n", err)
	}
	fmt.Println("Listener initialised")

	conn, err := net.DialTimeout("tcp", peerIP, 10*time.Second)
	if err != nil {
		log.Fatal("An error has occured: ", err)
	}
	fmt.Println("Connection initialised")

	wg.Add(2)
	go listen(listener, &wg)
	go writeToPeer(conn, &wg)
	wg.Wait()
}

func writeToPeer(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Entered WTP func")

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Println("Enter the text to send...\n")
		input, _ := reader.ReadString('\n')

		_, err := conn.Write([]byte(input))

		fmt.Println("Written to socket")

		if err != nil { // handle error inside of func
			log.Fatal("Error encountered when writing to peer:\n", err)
		}
	}

}

func listen(listener net.Listener, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := listener.Accept()

	if err != nil {
		log.Fatal("Error listening with the listener")
	}

	fmt.Println("Listener is listening")

	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		fmt.Println("Reading...")

		if err != nil {
			log.Fatal("Error occured when reading buffer:\n", err)
		}

		fmt.Println(string(buffer[:n]))

		if string(buffer[:n]) == "end" {
			break
		}

	}

}
