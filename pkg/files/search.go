package files

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/flanksource/commons/logger"
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

		files := readFilesLines(b.Paths, q.Labels)
		for _, content := range files {
			res.Results = append(res.Results, content...)
		}
	}

	return res, nil
}

type logsPerFile map[string][]logs.Result

// readFilesLines takes a list of file paths and returns each lines of those files.
// If labels are also passed, it'll attach those labels to each lines of those files.
func readFilesLines(paths []string, labelsToAttach map[string]string) logsPerFile {
	fileContents := make(logsPerFile, len(paths))
	for _, path := range paths {
		fInfo, err := os.Stat(path)
		if err != nil {
			logger.Warnf("error get file stat. path=%s; %w", path, err)
			continue
		}

		file, err := os.Open(path)
		if err != nil {
			logger.Warnf("error opening file. path=%s; %w", path, err)
			continue
		}

		// All lines of the same file will share these labels
		labels := mergeMap(map[string]string{"filepath": path}, labelsToAttach)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fileContents[path] = append(fileContents[path], logs.Result{
				Time:    fInfo.ModTime().Format(time.RFC3339),
				Labels:  labels,
				Message: strings.TrimSpace(scanner.Text()),
			})
		}
	}

	return fileContents
}

func matchQueryLabels(want, have map[string]string) bool {
	for label, val := range want {
		if val != have[label] {
			return false
		}
	}

	return true
}

// mergeMap will merge map b into a.
// On key collision, map b takes precedence.
func mergeMap(a, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}

	return a
}
