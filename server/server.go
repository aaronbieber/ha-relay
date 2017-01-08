/*
Home Automation Relay Server

See README.md for details.
*/
package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Function to test for the existence of a file, safely.
func fileExists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// The main loop.
func main() {
	// Consume command line switches.
	var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
	var port = flag.Int("port", 8765, "The port to listen on; default is 8765.")
	var pipeflag = flag.String("pipe", "~/conduit", "A named pipe to listen on for external commands; default is ~/conduit.")
	flag.Parse()

	// Argument sanity check.
	pipef, _ := filepath.Abs(*pipeflag)

	pipeExists, _ := fileExists(pipef)
	if !pipeExists {
		panic("Given pipe file does not exist.")
	}

	// This will hold connected clients.
	var clients []net.Conn

	rm_queue := make(chan net.Conn)
	add_queue := make(chan net.Conn)
	cmd_queue := make(chan string)

	// Add a client to the clients slice.
	add_client := func(client net.Conn) {
		clients = append(clients, client)
	}

	// Remove a client from the clients slice.
	rm_client := func(client net.Conn) {
		index := -1
		for i, v := range clients {
			if client == v {
				index = i
			}
		}

		if index == -1 {
			return
		}

		clients = append(clients[:index], clients[index+1:]...)
	}

	// Handle messages through the client management channels.
	go func() {
		for {
			select {
			case c := <-add_queue:
				add_client(c)
			case c := <-rm_queue:
				rm_client(c)
			}
		}
	}()

	log.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	log.Printf("Listening on %s.\n", src)

	defer listener.Close()

	// Listen for new connections, handle connection failures.
	go func() {
		for {
			conn, err := listener.Accept()
			if conn == nil || err != nil {
				log.Printf("Some connection error: %s\n", err)
				continue
			}

			log.Printf("Accepted connection from %s.\n", conn.LocalAddr())
			conn.Write([]byte("hello\n"))

			// Send a heartbeat byte to the client every ten seconds.
			go func() {
				for {
					log.Printf("Sending heartbeat to %s.\n", conn.LocalAddr())
					_, err := conn.Write([]byte("ping\n"))
					if err != nil {
						log.Printf("Closing connection %s.\n", conn.LocalAddr())
						conn.Close()
						return
					}
					time.Sleep(10 * time.Second)
				}
			}()

			add_queue <- conn
		}
	}()

	// Open the pipe for reading.
	log.Println("Opening pipe listener...")
	conduit, err := os.Open(pipef)
	if err != nil {
		panic(err)
	}
	conduit_reader := bufio.NewReader(conduit)

	// Read from the pipe and post to the cmd_queue
	go func() {
		for {
			text, err := conduit_reader.ReadString('\n')
			if len(text) > 0 {
				if err == nil {
					text = text[:len(text)-1]
				}
				log.Printf("Conduit received: %s.", text)
				cmd_queue <- text
			}
		}
	}()

	// Monitor the cmd_queue and relay those commands to all connected clients
	for {
		select {
		case cmd := <-cmd_queue:
			log.Printf("Command: %s", cmd)

			if (len(clients)) > 0 {
				log.Printf("Writing to %d client(s).\n", len(clients))
			}

			bad_clients := []net.Conn{}
			for _, client := range clients {
				_, err := client.Write([]byte(cmd + "\n"))
				if err != nil {
					bad_clients = append(bad_clients, client)
					log.Printf("Bad connection: %s.\n", client.LocalAddr())
				}
			}

			// Remove bad connections.
			for _, client := range bad_clients {
				client.Close()
				rm_queue <- client
			}
		}
	}
}
