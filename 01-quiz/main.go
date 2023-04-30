package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string // Question
	a string // Answer
}

func main() {
	// Flags
	csvFilename := flag.String("csv", "data/problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in secconds")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz order")
	flag.Parse()

	// Read in problems
	problems := readFile(*csvFilename)

	// Shuffle problems
	if *shuffle {
		rand.New(rand.NewSource(time.Now().UnixMicro()))
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	// Scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)

	// Start time - Ask user for input
	fmt.Print("Press Enter to Start Quiz.")
	scanner.Scan()
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// Loop though Problems
	correct := 0
	for i, p := range problems {

		// Ask Question
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)

		// Get user input
		go func() {
			scanner.Scan()
			answerCh <- strings.ToLower(strings.TrimSpace(scanner.Text()))
		}()

		if checkAnswer(*timer, p, &correct, answerCh) {
			break // Timer has expired
		}
	}

	// Final score output
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func readFile(fileName string) (problems []problem) {
	// Open File
	file, error := os.Open(fileName)
	if error != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", fileName))
	}

	// Read in line by line
	r := csv.NewReader(file)
	lineNum := 0
	for {
		line, err := r.Read()
		lineNum++

		if err == io.EOF {
			break
		}

		if err != nil {
			exit(fmt.Sprintf("Failed to parse the provided CSV file, Line %d: %s", lineNum, line))
		}

		// Parse Line
		parseLine(&problems, line)
	}

	return problems
}

func parseLine(problems *[]problem, line []string) {
	*problems = append(*problems, problem{
		q: line[0],
		a: strings.ToLower(strings.TrimSpace(line[1])), // We currently only expect single word/number answers. TrimSpaces will only trim left and right most, which is fine as we wont take " 1 1 " as an answer for 11.
	})
}

func checkAnswer(timer time.Timer, p problem, c *int, answerCh chan string) (isFinished bool) {
	select {
	case <-timer.C:
		fmt.Printf("\n")
		return true
	case answer := <-answerCh:
		if answer == p.a {
			*c++ // Answer was correct increment count
		}
		return false
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
