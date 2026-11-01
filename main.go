package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
	"unsafe"

	mic "github.com/MarkKremer/microphone"
	upnp "github.com/NebulousLabs/go-upnp"
	"github.com/gopxl/beep"
)

/*
func main() {

}
*/

func main() {
	err := mic.Init()
	if err != nil {
		log.Fatalf("Failed to initialize PortAudio: %v", err)
	}

	defer mic.Terminate()

	audioData := recordAudio()
	fmt.Println(audioData)
	/*
		var wg sync.WaitGroup

		peerIP := "127.0.0.1:5643"

		// bind first, so the port is open before we dial
		listener, err := net.Listen("tcp", ":5643") // chose a random unused port
		if err != nil {
			log.Fatal("Error creating listener obj:\n", err)
		}
		fmt.Println("Listener initialised")

		conn, err := net.DialTimeout("tcp", peerIP, 10*time.Second) // timeout takes 2 minutes if you dont define it
		if err != nil {
			log.Fatal("An error has occured: ", err)
		}
		fmt.Println("Connection initialised")

		stopChan := make(chan int) // used to close the uPnP port once the listener receives the message

		wg.Add(3)
		go handleUPnPPort(stopChan, &wg)
		go listen(listener, &wg, stopChan)
		go writeToPeer(conn, &wg)
		wg.Wait()
	*/
}

func writeToPeer(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter the text to send...\n")
	input, _ := reader.ReadString('\n')

	_, err := conn.Write([]byte(input)) // write the message to the IP that was passed

	fmt.Println("Written to socket")

	if err != nil { // handle error inside of func
		log.Fatal("Error encountered when writing to peer:\n", err)
	}
}

func listen(listener net.Listener, wg *sync.WaitGroup, c chan int) {
	defer wg.Done()

	conn, err := listener.Accept() // waits until there is data being sent into the port

	if err != nil {
		log.Fatal("Error listening with the listener")
	}

	fmt.Println("Listener is listening")

	buffer := make([]byte, 4096) // TODO: dynamic buffer

	for {
		n, err := conn.Read(buffer)
		fmt.Println("Reading...")

		if err != nil {
			log.Fatal("Error occured when reading buffer:\n", err)
		}

		fmt.Println(string(buffer[:n]))
		c <- 0 // send into stopChan to tell uPnP handler to close the port

	}

}

func handleUPnPPort(c chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	IGD, err := upnp.Discover() // initiate a IGD protocol to do the port handling

	if err != nil {
		log.Fatal("error initiating IGD:\n", err)
	}

	IGD.Forward(5643, "uPnP port for my p2p transfer app") // desc is irrelevant

	<-c // stalls func until the channel gets some data

	err = IGD.Clear(5643) // close the port

	if err != nil {
		log.Fatal("Error occured when closing port:\n", err) // god forbid this happens
	}

}

func recordAudio() []byte {
	var sampleRate beep.SampleRate = 44100

	stream, _, err := mic.OpenDefaultStream(sampleRate, 1)
	if err != nil {
		log.Fatal("Error opening microphone stream: ", err)
	}
	defer stream.Close()

	err = stream.Start()
	if err != nil {
		log.Fatal("Error starting microphone: ", err)
	}

	fmt.Println("Recording for 5 seconds...")

	var samples []float32
	done := make(chan struct{}) // struct{} is a zero-size type, so it doesn't allocate memory

	go func() {
		time.Sleep(5 * time.Second)
		stream.Stop()
		stream.Close()
		close(done) // signals that recording is done
	}()

	buffer := make([][2]float64, 1024) // buffer to hold audio samples

	for {
		select {
		case <-done:
			fmt.Println("Recording finished") // checks if the channel is closed
			sampleBytes := unsafe.Slice((*byte)(unsafe.Pointer(&samples[0])), len(samples)*4)
			return sampleBytes
		default:
			n, ok := stream.Stream(buffer) // fill the buffer with audio samples
			if !ok {
				return nil // if the stream is closed, return the samples
			}

			for i := 0; i < n; i++ {
				samples = append(samples, float32(buffer[i][0])) // convert to float32 and append to samples
			}

		}
	}
}
