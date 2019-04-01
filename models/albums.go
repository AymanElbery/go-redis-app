package models

import (
	"errors"
	// Import the Radix.v2 pool package ...
	"log"
	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

// Declare a global db variable to store the Redis connection pool ...
var db *pool.Pool

func init() {
	var err error

	// Establish a pool of 10 connections ...
	db, err = pool.New("tcp", "localhost:6379", 10)
	if err != nil {
		log.Panic(err)
	}
}

// ErrNoAlbum error message ...
var ErrNoAlbum = errors.New("models: no album found")

// Album struct ...
type Album struct {
	Title  string
	Artist string
	Price  float64
	Likes  int
}

func populateAlbum(allAlbum map[string]string) (*Album, error) {
	var err error
	ab := new(Album)
	ab.Title = allAlbum["title"]
	ab.Artist = allAlbum["artist"]
	ab.Price, err = strconv.ParseFloat(allAlbum["price"], 64)
	if err != nil {
		return nil, err
	}
	ab.Likes, err = strconv.Atoi(allAlbum["likes"])
	if err != nil {
		return nil, err
	}
	return ab, nil
}

// FindAlbum func ...
func FindAlbum(id string) (*Album, error) {

	// fetch a single Redis connection from the pool ...
	conn, err := db.Get()
	if err != nil {
		return nil, err
	}

	defer db.Put(conn)

	// Fetch the details of a specific album ...
	allAlbum, err := conn.Cmd("HGETALL", "album:"+id).Map()
	if err != nil {
		return nil, err
	} else if len(allAlbum) == 0 {
		return nil, ErrNoAlbum
	}

	return populateAlbum(allAlbum)
}

// IncrementLikes func ...
func IncrementLikes(id string) error {
	conn, err := db.Get()
	if err != nil {
		return err
	}
	defer db.Put(conn)

	// check that an album with the given id exists ...
	exists, err := conn.Cmd("EXISTS", "album:"+id).Int()
	if err != nil {
		return err
	} else if exists == 0 {
		return ErrNoAlbum
	}

	// Use the MULTI command to inform Redis that we are starting a new transaction ...
	err = conn.Cmd("MULTI").Err
	if err != nil {
		return err
	}

	// Increment the number of likes in the album hash by 1 ...
	err = conn.Cmd("HINCRBY", "album:"+id, "likes", 1).Err
	if err != nil {
		return err
	}

	// And we do the same with the increment on our sorted set.
	err = conn.Cmd("ZINCRBY", "likes", 1, id).Err
	if err != nil {
		return err
	}

	// Execute both commands in our transaction together as an atomic group ...
	err = conn.Cmd("EXEC").Err
	if err != nil {
		return err
	}
	return nil
}

// FindTopThree func ...
func FindTopThree() ([]*Album, error) {
	conn, err := db.Get()
	if err != nil {
		return nil, err
	}
	defer db.Put(conn)

	// Begin an infinite loop.
	for {
		// Instruct Redis to watch the likes sorted set for any changes ...
		err = conn.Cmd("WATCH", "likes").Err
		if err != nil {
			return nil, err
		}

		// Use the ZREVRANGE command to fetch the album ids with the highest liks ...
		reply, err := conn.Cmd("ZREVRANGE", "likes", 0, 2).List()
		if err != nil {
			return nil, err
		}

		// Use the MULTI command to inform Redis that we are starting a new transaction ...
		err = conn.Cmd("MULTI").Err
		if err != nil {
			return nil, err
		}

		// Loop through the ids returned by ZREVRANGE ...
		for _, id := range reply {
			err := conn.Cmd("HGETALL", "album:"+id).Err
			if err != nil {
				return nil, err
			}
		}

		// Execute the transaction ...
		ereply := conn.Cmd("EXEC")
		if ereply.Err != nil {
			return nil, err
		} else if ereply.IsType(redis.Nil) {
			continue
		}

		// convert the transaction reply to an array ...
		areply, err := ereply.Array()
		if err != nil {
			return nil, err
		}

		// Create a new slice to store the album details ...
		abs := make([]*Album, 3)

		// Iterate through the array of Resp objects ...
		for i, reply := range areply {
			mreply, err := reply.Map()
			if err != nil {
				return nil, err
			}
			ab, err := populateAlbum(mreply)
			if err != nil {
				return nil, err
			}
			abs[i] = ab
		}

		return abs, nil
	}
}
