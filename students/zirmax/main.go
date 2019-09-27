package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

var filename string
var timeout int
var shuffleEnabled bool

func main() {

	flag.StringVar(&filename, "filename", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.IntVar(&timeout, "timeout", 30, "the time limit for the quiz in seconds")
	flag.BoolVar(&shuffleEnabled, "shuffle", false, "if the questions shall be shuffled")
	flag.Parse()
	csvContent, err := readCSV(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	problems := []problem{}

	for _, line := range csvContent {
		problems = append(problems, problem{question: line[0], answer: line[1]})
	}

	if shuffleEnabled {
		problems = shuffle(problems)
	}

	var correctAnswers, falseAnswers int

	fmt.Println("Press enter to start the quiz ...")
	bufio.NewScanner(os.Stdin).Scan()
	defer os.Stdin.Close()

	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	defer timer.Stop()

	go func() {
		<-timer.C
		fmt.Printf("\nCorrect Answers %v, False Answers %v\n", correctAnswers, falseAnswers)
	}()

	var answer string
	done := make(chan bool)

	go func() {
		for _, problem := range problems {
			fmt.Printf("%v ", problem.question)
			fmt.Scan(&answer)
			if answer == problem.answer {
				correctAnswers++
			} else {
				falseAnswers++
			}
		}
		done <- true
	}()
	select {
	case <-done:
	case <-timer.C:
		fmt.Println("time's up!")
	}
	fmt.Printf("\nCorrect Answers %v, False Answers %v\n", correctAnswers, falseAnswers)
}

func readCSV(filename string) ([][]string, error) {
	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func shuffle(problems []problem) []problem {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	return problems
}
