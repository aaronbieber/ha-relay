/*
Home Automation Relay Client

See README.md for details.
*/
package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"strconv"
	"time"
)

var host = flag.String("host", "localhost", "The hostname or IP to connect to; defaults to \"localhost\".")
var port = flag.Int("port", 8765, "The port to connect to; defaults to 8765.")

func main() {
	flag.Parse()

	dest := *host + ":" + strconv.Itoa(*port)
	log.Printf("Connecting to %s...\n", dest)

	for {
		conn, err := net.Dial("tcp", dest)

		if err != nil {
			log.Println("Failed to connect; trying again in 10 seconds...")
			time.Sleep(10 * time.Second)
			continue
		}

		log.Printf("Connected to %s.\n", conn.RemoteAddr())

		for {
			scanner := bufio.NewScanner(conn)
			conn.SetDeadline(time.Now().Add(15 * time.Second))
			ok := scanner.Scan()
			text := scanner.Text()

			if len(text) > 0 {
				log.Printf("Received command: %s.\n", text)
			}

			if !ok {
				log.Printf("Reached EOF, dropping this connection.")
				conn.Close()
				break
			}

			if scanner.Err() != nil {
				log.Printf("Error reading from %s.\n", conn.RemoteAddr())
				log.Print(scanner.Err())
				break
			}
		}
	}
}
