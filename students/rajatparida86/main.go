package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	flag "github.com/ogier/pflag"
)

var (
	file    = flag.StringP("file", "f", "problems.csv", "Quiz file")
	timeOut = flag.IntP("timeout", "t", 20, "Timeout for quiz")
	random  = flag.BoolP("shuffle", "s", false, "Shuffle the questions in the quiz")
	correct int
)

type problem struct {
	question string
	answer   string
}

func newProblem(question string, answer string) *problem {
	return &problem{
		question,
		answer,
	}
}

func main() {
	flag.Parse()
	quiz, err := getQuiz()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	if *random {
		quiz = shuffleQuiz(quiz)
	}
	ch := make(chan bool, 1)
	defer close(ch)
	timer := time.NewTimer(time.Duration(*timeOut) * time.Second)
	defer timer.Stop()
	fmt.Print("The kid quiz. Hit the Enter key to start...")

	go runQuiz(quiz, ch)

	select {
	case <-ch:
		fmt.Printf("Quiz complete.")
	case <-timer.C:
		fmt.Printf("\nYou ran out of time.")
	}
	fmt.Printf("\nScore: %d/%d", correct, len(quiz))
}

func shuffleQuiz(quiz []*problem) []*problem {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(quiz); n > 0; n-- {
		randomIndex := r.Intn(n)
		quiz[n-1], quiz[randomIndex] = quiz[randomIndex], quiz[n-1]
	}
	return quiz
}

func runQuiz(quiz []*problem, ch chan bool) {
	reader := bufio.NewReader(os.Stdin)
	_, _, err := reader.ReadRune()
	if err != nil {
		fmt.Printf("Abort!!!Abort!!!Something bad happened....\n %s", err)
		os.Exit(1)
	}
	for _, p := range quiz {
		fmt.Printf("%v=", p.question)
		res, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Unable to read your answer, %s", err)
			os.Exit(1)
		}
		answer := cleanUP(res)

		if answer == strings.ToLower(p.answer) {
			correct++
		}
	}
	ch <- true
}

func cleanUP(s string) string {
	return strings.TrimSpace(
		strings.ToLower(
			strings.Replace(s, "\n", "", -1)))
}

func getQuiz() ([]*problem, error) {
	var quiz []*problem
	csvFile, err := os.Open(*file)
	if err != nil {
		fmt.Printf("Unable to get the quiz: %s", err)
		os.Exit(1)
	}
	defer csvFile.Close()

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get the quiz %s", err)
	}
	for _, line := range lines {
		if len(line) > 2 {
			return nil, errors.New("Quiz input file invalid")
		}
		quiz = append(quiz, newProblem(line[0], line[1]))
	}
	return quiz, nil
}
