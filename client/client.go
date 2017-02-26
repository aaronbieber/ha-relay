/*
Home Automation Relay Client

See README.md for details.
*/
package main

import (
	"bufio"
	"flag"
	"github.com/aaronbieber/ha-relay/config"
	"github.com/aaronbieber/ha-relay/crypto"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var host = flag.String("host", "localhost", "The hostname or IP to connect to; defaults to \"localhost\".")
var port = flag.Int("port", 8765, "The port to connect to; defaults to 8765.")
var cmd_dir = flag.String("cmd_dir", "commands", "Path to a directory containing valid command scripts.")

func main() {
	conf := config.Conf()

	flag.Parse()
	key := []byte(conf.Main.Key)

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
				switch {
				case text == "ping":
					log.Printf("Ping? Pong.")

				case text == "hello":
					log.Printf("The server said hello. Hello, server.")

				default:
					go command(text, key)
				}
			}

			if !ok {
				log.Printf("Reached EOF, dropping this connection.")
				conn.Close()
				break
			}

			if scanner.Err() != nil {
				log.Printf("!! Error reading from %s.\n", conn.RemoteAddr())
				log.Print(scanner.Err())
				break
			}
		}
	}
}

func scanCommands() map[string]string {
	var commands = make(map[string]string)

	files, err := ioutil.ReadDir(*cmd_dir)
	if err != nil {
		panic("Could not read the commands directory.")
	}

	for _, f := range files {
		ext := filepath.Ext(f.Name())
		name := f.Name()[0 : len(f.Name())-len(ext)]
		commands[name] = filepath.Join(*cmd_dir, f.Name())
	}

	return commands
}

func command(command string, key []byte) {
	command, err := crypto.Decrypt(key, command)
	if err != nil {
		log.Printf("crypto> !! Error decrypting %s", command)
		return
	}

	commands := scanCommands()

	if script, ok := commands[command]; ok {
		log.Printf("%s > Executing %s...", command, script)

		out, err := exec.Command(script).Output()

		for _, s := range strings.Split(string(out), "\n") {
			if len(s) > 0 {
				log.Printf("%s > %s", command, s)
			}
		}

		if err != nil {
			log.Printf("%s > !! Error running %s: %s", command, script, err)
		}

	} else {
		log.Printf("!! Unknown command: %s", command)
	}
}
