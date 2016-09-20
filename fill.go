package main

import (
	"fmt"
	// Import the Radix.v2 redis package.
	"github.com/mediocregopher/radix.v2/redis"
	"log"
)

func main() {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	// Importantly, use defer to ensure the connection is always properly
	// closed before exiting the main() function.
	defer conn.Close()

	// key, followed by the various hash fields and values).
	resp := conn.Cmd("HMSET", "hash:test", "hash", "test", "url", "http://www.lalala.com", "sourceip", "127.0.0.1", "date", "05/10/2016", "clickstats", 0)
	// Check the Err field of the *Resp object for any errors.
	if resp.Err != nil {
		log.Fatal(resp.Err)
	}

	fmt.Println("Electric Ladyland added!")
}
