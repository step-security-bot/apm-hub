package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/flanksource/flanksource-ui/apm-hub/api/logs"
)

type FileSearch struct {
	FilesBackend []logs.FileSearchBackend
}

func (t *FileSearch) Search(q *logs.SearchParams) (r logs.SearchResults, err error) {
	var res logs.SearchResults

	for _, b := range t.FilesBackend {
		if !matchQueryLabels(q.Labels, b.Labels) {
			continue
		}

		files, err := readFilesLines(b.Paths)
		if err != nil {
			return res, fmt.Errorf("readFilesLines(); %w", err)
		}

		for _, content := range files {
			res.Results = append(res.Results, content...)
		}
	}

	return res, nil
}

type logsPerFile map[string][]logs.Result

// readFilesLines will take a list of file paths
// and then return each lines of those files.
func readFilesLines(paths []string) (logsPerFile, error) {
	fileContents := make(logsPerFile, len(paths))
	for _, path := range paths {
		fInfo, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("error get file stat. path=%s; %w", path, err)
		}

		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("error opening file. path=%s; %w", path, err)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fileContents[path] = append(fileContents[path], logs.Result{
				Time: fInfo.ModTime().Format(time.RFC3339),
				// Labels: , all the records will have the same labels. Is it necessary to add it here?
				Message: strings.TrimSpace(scanner.Text()),
			})
		}
	}

	return fileContents, nil
}

func matchQueryLabels(want, have map[string]string) bool {
	for label, val := range want {
		if val != have[label] {
			return false
		}
	}

	return true
}
