package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

func getRecords(fName string) [][]string {
	f, err := os.Open(fName)
	defer f.Close()
	if err != nil {
		fmt.Println("Error reading file:", err.Error())
		os.Exit(1)
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		fmt.Println("Error reading from file:", err.Error())
		os.Exit(1)
	}
	for _, record := range records {
		if len(record) != 2 {
			fmt.Println("csv-record has to many/few fields (expected 2):", record)
			os.Exit(1)
		}
	}
	return records
}

type GameData struct {
	totalQuestions, askedQuestions, correct, wrong int
}

func askUser(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Q: %s? A: ", question)
	answer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading answer", err.Error())
		os.Exit(1)
	}
	answer = strings.TrimSpace(answer)
	answer = strings.ToLower(answer)

	return answer
}

func evalUserAnswer(correctAnswer, userAnswer string, game *GameData) {
	if userAnswer == correctAnswer {
		game.correct++
	} else {
		game.wrong++
	}
	game.askedQuestions++

}

func main() {
	fmt.Println("Hello, World!")
	argsWithoutProg := os.Args[1:]

	problemFileName := "problems.csv"
	if len(argsWithoutProg) == 1 {
		problemFileName = argsWithoutProg[0]
	} else if len(argsWithoutProg) > 1 {
		fmt.Println("Only 1 problem file expected")
		os.Exit(1)
	}
	records := getRecords(problemFileName)

	c := make(chan int)
	defer close(c)

	var game GameData
    game.totalQuestions = len(records)

	logic := func(record []string) {
		question, answer := record[0], record[1]
		userAnswer := askUser(question)
		evalUserAnswer(answer, userAnswer, &game)
		c <- 0
	}

	timer := time.NewTimer(3 * time.Second)
Loop:
	for _, record := range records {
		go logic(record)
		select {
		case <-timer.C:
			fmt.Println("To slow :(")
			break Loop
		case <-c:
            // User answered question
		}
	}

    fmt.Printf("Total questions: %v, asked: %v, correct: %v, wrong: %v\n", game.totalQuestions, game.askedQuestions, game.correct, game.wrong)
}
