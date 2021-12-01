package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	last := -1
	numInc := 0
	for scanner.Scan() {
		line := scanner.Text()
		this, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}
		if this > last && last != -1 {
			numInc += 1
		}
		last = this
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Num increased: %d\n", numInc)
}
