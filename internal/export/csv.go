package export

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"time"
	"track/internal/store"
)

func ExportToCSV(tasks []store.Task) error {
	tasksSlice := make([][]string, 0, len(tasks)+1)
	tasksSlice = append(tasksSlice, []string{"Date", "Title"})
	for _, t := range tasks {
		createdAt, _ := time.Parse(time.RFC3339, t.CreatedAt)
		tasksSlice = append(tasksSlice, []string{createdAt.Format("2003-02-06"), t.Title})
	}
	filePath, _ := getFilePath()
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	return writer.WriteAll(tasksSlice)
}

func getFilePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, "track_export_"+time.Now().Format("2003-02-06")), nil
}
