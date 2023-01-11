package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

		files, err := readFileLines(b.Paths)
		if err != nil {
			return res, fmt.Errorf("readFileLines(); %w", err)
		}

		for _, content := range files {
			res.Results = append(res.Results, content...)
		}
	}

	// TODO: Need to implement pagination but is it per file?

	return res, nil
}

type logsPerFile map[string][]logs.Result

// readFileLines will take a list of file paths
// and then return each lines of those files.
func readFileLines(paths []string) (logsPerFile, error) {
	fileContents := make(logsPerFile, len(paths))
	for _, path := range paths {
		content, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("error reading file_path=%s; %w", path, err)
		}

		scanner := bufio.NewScanner(content)
		for scanner.Scan() {
			fileContents[path] = append(fileContents[path], logs.Result{
				// Id: , guess I can ignore the Id at this stage
				// Time: , not sure how to reliably get the time here. This varies based on the log type.
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
