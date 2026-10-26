package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	peerIP := "175.33.209.242:5643"

	message := "hey"
	var payload []byte

	for _, char := range message {
		payload = append(payload, byte(char))
	}

	conn, err := net.Dial("tcp", peerIP)

	if err != nil {
		log.Fatal("An error has occured: ", err)
	}

	wg.Add(2)
	go writeToPeer(conn, payload, &wg)
	go listen(peerIP, &wg)
	wg.Wait()

}

func writeToPeer(conn net.Conn, message []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	n, err := conn.Write(message)

	if err != nil { // handle error inside of func
		log.Fatal("Error encountered when writing to peer:\n", err)
	}

	fmt.Println(n)

}

func listen(address string, wg *sync.WaitGroup) {
	defer wg.Done()

	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("Error creating listener obj:\n", err)
	}

	conn, err := listener.Accept()

	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			log.Fatal("Error occured when reading buffer:\n", err)
		}

		fmt.Println(string(buffer[:n]))

		if string(buffer[:n]) == "end" {
			break
		}

	}

}
