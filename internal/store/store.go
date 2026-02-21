package store

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"
)

type Task struct {
	Title     string `json:"title"`
	CreatedAt string `json:"createdAt"`
}

const fileName = "tasks.json"

func AddToList(task Task) error {
	tasks, err := getFileContent()
	if err != nil {
		return err
	}

	tasks = append(tasks, task)
	return updateFileContent(tasks)
}

func ListTasks(filter string) ([]Task, error) {
	var validFilters = []string{
		"today", "month", "year",
	}

	if filter != "" {
		if slices.Index(validFilters, filter) == -1 {
			return nil, errors.New("invalid filter")
		}
	}

	tasks, err := getFileContent()
	if err != nil || filter == "" {
		return tasks, err
	}

	var filteredTasks []Task
	now := time.Now()

	for _, t := range tasks {
		createdAt, err := time.Parse(time.RFC3339, t.CreatedAt)
		if err != nil {
			continue // skip invalid timestamps
		}

		createdAt = createdAt.Local()
		canAdd := false

		switch filter {

		case "today":
			canAdd =
				createdAt.Year() == now.Year() &&
					createdAt.YearDay() == now.YearDay()

		case "month":
			canAdd =
				createdAt.Year() == now.Year() &&
					createdAt.Month() == now.Month()

		case "year":
			canAdd =
				createdAt.Year() == now.Year()
		}

		if canAdd {
			filteredTasks = append(filteredTasks, t)
		}
	}

	return filteredTasks, nil
}
func getFileContent() ([]Task, error) {
	filePath, err := getFilePath()
	if err != nil {
		return nil, err
	}

	// If file does not exist, return empty slice
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []Task{}, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tasks []Task
	err = json.NewDecoder(file).Decode(&tasks)

	// If file is empty → treat as empty list
	if errors.Is(err, io.EOF) {
		return []Task{}, nil
	}

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func updateFileContent(tasks []Task) error {
	filePath, err := getFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

func getFilePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, fileName), nil
}
