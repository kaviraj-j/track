package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"track/internal/export"
	"track/internal/store"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		log.Fatal("Use: `track help` to know how to use commands")
	}

	switch args[1] {
	case "ls", "list":
		filter, _ := getTagValue("filter", args)
		tasks, err := store.ListTasks(filter)
		if err != nil {
			log.Fatal(err)
		}
		printTasks(tasks)
	case "export":
		filter, _ := getTagValue("filter", args)
		tasks, err := store.ListTasks(filter)
		if err != nil {
			log.Fatal(err)
		}
		if export.ExportToCSV(tasks) != nil {
			log.Fatal("failed to export tasks")
		}

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

func getTagValue(tagName string, args []string) (value string, found bool) {
	for i, arg := range args {
		if arg == "--"+tagName {
			if len(args) > i+1 {
				value = args[i+1]
				found = true
				return
			}
			return
		}
	}
	return
}
