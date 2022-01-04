package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	problemsPath = flag.String("problems_path", "./problems.csv", "Set the file to read problems data, default to problems.csv in current direcotory")
	timeLimit    = flag.Int("time_limit", 30, "Specify a time limit within which user must be able to answer all questions")
)

type quiz struct {
	timer *time.Timer
}

func (q *quiz) run() {
	go func() {
		<-q.timer.C
		q.stop()
	}()
	loadProblems()
	askProblems()
	checkProblems()
	q.timer.Stop()
}

func (q *quiz) stop() {
	fmt.Println("Time's up!")
	checkProblems()
	os.Exit(0)
}

type problem struct {
	prompt     string
	answer     string
	userAnswer string
}

var problems []problem

func checkErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func main() {
	// flag will only be recoganized after parsing
	flag.Parse()
	fmt.Printf("Current time limit is %d seconds, start the quiz?(Y/n)\n", *timeLimit)
	r := bufio.NewReaderSize(os.Stdin, 1)
	input, _ := r.ReadString('\n')
	if strings.TrimSpace(input) == "Y" {
		q := quiz{
			timer: time.NewTimer(time.Duration(*timeLimit) * time.Second),
		}
		q.run()
	}
}

func loadProblems() {
	f, err := os.Open(*problemsPath)
	defer f.Close()
	checkErr(err)
	csvReader := csv.NewReader(f)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		checkErr(err)
		problems = append(problems, problem{
			prompt: record[0],
			answer: record[1],
		})
	}
}

func askProblems() {
	r := bufio.NewReader(os.Stdin)
	for i, p := range problems {
		fmt.Printf("Q: %s, what's your answer?\n", p.prompt)
		asn, err := r.ReadString('\n')
		checkErr(err)
		// NOTE:
		// 1. remember slice gives you the copy, not the pointer to the original value
		// 2. values read from bufio contains newline
		problems[i].userAnswer = strings.TrimSpace(asn)
	}
}

func checkProblems() {
	correctCount := 0
	for _, p := range problems {
		if p.answer == p.userAnswer {
			correctCount++
		}
	}
	fmt.Printf("You got %d of %d problems correct.\n", correctCount, len(problems))
}
