package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"track/internal/store"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		log.Fatal("Use: `track help` to know how to use commands")
	}

	switch args[1] {
	case "ls", "list":
		tasks, err := store.ListTasks()
		if err != nil {
			log.Fatal(err)
		}
		printTasks(tasks)

	case "add":
		if len(args) < 3 {
			log.Fatal("Usage: `track add <message>`")
		}
		msgToAdd := args[2]
		err := store.AddToList(store.Task{
			Title:     msgToAdd,
			CreatedAt: time.Now().UTC().Format(time.RFC3339), // now -> UTC -> formatted string
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func printTasks(tasks []store.Task) {
	for _, t := range tasks {
		createdAtTime, _ := time.Parse(time.RFC3339, t.CreatedAt)
		fmt.Printf("%s  -  %s\n", createdAtTime.Local().Format("Mon Jan 2"), t.Title)
	}
}
