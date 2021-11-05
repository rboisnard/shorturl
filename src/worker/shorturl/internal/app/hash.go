package app

import (
	"crypto/sha256"
	"fmt"
	"log"
)

// globals
const shortSize = 6
const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Shorten(abort <-chan struct{}, lurl string) <-chan string {
	// create channel to iterate over results
	results := make(chan string)

	// prepare a slice with sha256 and a copy of the first bytes at the end
	urlSha := sha256.Sum256([]byte(lurl))
	hash := make([]byte, sha256.Size+shortSize-1)
	copy(hash, urlSha[:])
	copy(hash[sha256.Size:], urlSha[:shortSize-1])
	log.Printf("+++ full hash '%x'", hash[:])

	symbolsLen := len(symbols)

	// start convertion with the symbol list
	for i := 0; i < shortSize-1; i++ {
		hash[i] = symbols[int(hash[i])%symbolsLen]
	}

	// start iterating in closure
	go func() {
		// close results channel when exiting to avoid memory leaks
		defer close(results)

		currentPos := 0
		for currentPos < sha256.Size {
			select {
			case <-abort:
				log.Printf("+++ aborting generator channel for Shorten")
				return
			default:
				// convert new element from the hash
				hash[currentPos+shortSize-1] = symbols[int(hash[currentPos+shortSize-1])%symbolsLen]
				log.Printf("+++ sending new short URL to generator channel in Shorten '%s'", hash[currentPos:currentPos+shortSize])
				results <- fmt.Sprintf("%s", hash[currentPos:currentPos+shortSize])
			}
			currentPos++
		}
	}()
	return results
}
