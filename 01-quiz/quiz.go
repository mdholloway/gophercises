package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

func awaitAnswer(reader *bufio.Reader, input chan string) {
	line, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	input <- strings.TrimRight(line, "\n")
}

func main() {
	s := flag.Uint("time", 30, "Number of seconds the quiz should run")
	filename := flag.String("filename", "problems.csv",
		"CSV file containing quiz questions and answers")
	flag.Parse()

	data, err := ioutil.ReadFile(*filename)
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(strings.NewReader(string(data)))
	records, err := csvReader.ReadAll()
	rand.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})

	inputReader := bufio.NewReader(os.Stdin)
	input := make(chan string)

	fmt.Println("You have", *s, "seconds to answer", len(records), "questions.")
	fmt.Println("Press ENTER to start the quiz.")

	_, err = inputReader.ReadBytes('\n')
	if err != nil {
		panic(err)
	}

	var timeUp <-chan time.Time
	timeUp = time.After(time.Duration(*s) * time.Second)

	correct := 0

	Quiz:
		for i := 0; i < len(records); i++ {
			record := records[i]
			question := record[0]
			expected := record[1]

			fmt.Println(question)
			go awaitAnswer(inputReader, input)

			select {
			case answer := <-input:
				if answer == expected {
					correct++
				}
			case <-timeUp:
				break Quiz
			}
		}

	fmt.Println("You answered", correct, "out of", len(records), "correctly in", *s, "seconds.")
}
