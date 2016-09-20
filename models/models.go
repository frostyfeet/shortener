package models

import (
	"errors"
	// Import the Radix.v2 pool package, NOT the redis package.
	"github.com/mediocregopher/radix.v2/pool"
	"log"
	"strconv"
)

// Declare a global db variable to store the Redis connection pool.
var db *pool.Pool

func init() {
	var err error
	// Establish a pool of 10 connections to the Redis server listening on
	// port 6379 of the local machine.
	db, err = pool.New("tcp", "localhost:6379", 10)
	if err != nil {
		log.Panic(err)
	}
}

// Create a new error message and store it as a constant. We'll use this
// error later if the FindAlbum() function fails to find an album with a
// specific id.
var ErrNoUrl = errors.New("models: no album found")

type Url struct {
	Hash   string
	Url    string
	Ip     string
	Date   string
	Clicks int
}

func populateUrl(reply map[string]string) (*Url, error) {
	var err error
	ab := new(Url)
	ab.Hash = reply["hash"]
	ab.Url = reply["url"]
	ab.Ip = reply["sourceip"]
	ab.Date = reply["date"]
	ab.Clicks, err = strconv.Atoi(reply["clickstats"])
	if err != nil {
		return nil, err
	}
	return ab, nil
}

func FindUrl(id string) (*Url, error) {
	// Use the connection pool's Get() method to fetch a single Redis
	// connection from the pool.
	conn, err := db.Get()
	if err != nil {
		return nil, err
	}
	// Importantly, use defer and the connection pool's Put() method to ensure
	// that the connection is always put back in the pool before FindAlbum()
	// exits.
	defer db.Put(conn)

	// Fetch the details of a specific album. If no album is found with the
	// given id, the map[string]string returned by the Map() helper method
	// will be empty. So we can simply check whether it's length is zero and
	// return an ErrNoAlbum message if necessary.
	reply, err := conn.Cmd("HGETALL", "album:"+id).Map()
	if err != nil {
		return nil, err
	} else if len(reply) == 0 {
		return nil, ErrNoUrl
	}

	return populateUrl(reply)
}

func GetUrl(id string) (string, error) {
	// Use the connection pool's Get() method to fetch a single Redis
	// connection from the pool.
	conn, err := db.Get()
	if err != nil {
		return "", err
	}
	// Importantly, use defer and the connection pool's Put() method to ensure
	// that the connection is always put back in the pool before FindAlbum()
	// exits.
	defer db.Put(conn)

	// Fetch the details of a specific album. If no album is found with the
	// given id, the map[string]string returned by the Map() helper method
	// will be empty. So we can simply check whether it's length is zero and
	// return an ErrNoAlbum message if necessary.
	reply, err := conn.Cmd("HGET", "hash:"+id, "url").Str()
	if err != nil {
		return "", err
	} else if len(reply) == 0 {
		return "", ErrNoUrl
	}

	return reply, err
}
