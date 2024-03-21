package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	// TODO: Instead of accepting the filename, maybe we should accept the file path.
	filenamePtr := flag.String("filename", "problems", "Name of the CSV which has the quiz")

	quizDurationPtr := flag.Int("duration", 30, "Duration of the quiz")

	flag.Parse()

	file, err := os.Open(fmt.Sprintf("./%v.csv", *filenamePtr))

	if err != nil {

		log.Fatal(err)

	}

	csvReader := csv.NewReader(file)

	records, err := csvReader.ReadAll()

	file.Close()

	numOfCorrectAns := 0

	quizDuration := time.NewTimer(time.Duration(*quizDurationPtr) * time.Second)

out:
	for _, v := range records {

		c := make(chan string)

		go func(v []string, c chan string) {
			fmt.Printf("%v => ", v[0])

			buf := bufio.NewReader(os.Stdin)

			userAnswer, err := buf.ReadString('\n')

			if err != nil {

				log.Fatal(err)

			}

			c <- strings.TrimRight(userAnswer, "\n")

		}(v, c)

		select {

		case res := <-c:

			if res == v[1] {
				numOfCorrectAns++
			}

		case <-quizDuration.C:

			fmt.Printf("\n\n")

			close(c)

			break out

		}

		fmt.Println()

	}

	fmt.Printf("Final Score: %v/%v\n", numOfCorrectAns, len(records))

}
