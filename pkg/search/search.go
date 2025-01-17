package search

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	chunkSize   = 1 * 1024 * 1024  // 1MB
	BufferLimit = 10 * 1024 * 1024 // 10MB
)

type Searcher struct {
	file   *os.File
	reader *bufio.Reader
}

func NewSearcher(path string) (*Searcher, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &Searcher{
		file:   file,
		reader: nil,
	}, nil
}

func (s *Searcher) seekToLineStart() error {
	searchPos, err := s.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("failed to get current position: %w", err)
	}

	if searchPos == 0 {
		return nil
	}

	buf := make([]byte, chunkSize)

	var lineStart int64 = 0

	for {
		readStart := searchPos - chunkSize
		if readStart < 0 {
			readStart = 0
		}

		if _, err := s.file.Seek(readStart, io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek buffer: %w", err)
		}

		// Read until searchPos
		n, err := io.ReadFull(s.file, buf[:searchPos-readStart])
		if err != nil && err != io.ErrUnexpectedEOF {
			return fmt.Errorf("failed to read: %w", err)
		}

		data := buf[:n]

		// searching for the last newline
		for i := len(data) - 1; i >= 0; i-- {
			if data[i] == '\n' {
				lineStart = readStart + int64(i) + 1 // +1 to move past the newline
				_, err = s.file.Seek(lineStart, io.SeekStart)
				if err != nil {
					return fmt.Errorf("failed to seek: %w", err)
				}
				return nil
			}
		}

		searchPos = readStart

		// If we've reached the start of the file
		if searchPos == 0 {
			lineStart = 0
			break
		}
	}

	if _, err := s.file.Seek(lineStart, io.SeekStart); err != nil {
		return fmt.Errorf("failed to seek: %w", err)
	}

	return nil
}

// FindLine reads a sorted text file and returns the first line >= searchTerm.
// If no line is found, it returns an empty string.
func (s *Searcher) FindLine(searchTerm string) (string, error) {
	var err error

	bufio.NewScanner(s.file)
	s.reader = bufio.NewReaderSize(s.file, int(BufferLimit))

	result, err := s.search([]byte(searchTerm))
	if err != nil {
		return "", fmt.Errorf("failed to search: %w", err)
	}

	return result, nil
}

func (s *Searcher) search(searchTerm []byte) (string, error) {
	fileInfo, err := s.file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	size := fileInfo.Size()
	if size == 0 {
		return "", nil
	}

	left, right := int64(0), size

	var result []byte
	for left < right {
		mid := (left + right) / 2
		_, err := s.file.Seek(mid, io.SeekStart)
		if err != nil {
			return "", fmt.Errorf("failed to seek, mid: %w", err)
		}

		if err := s.seekToLineStart(); err != nil {
			return "", fmt.Errorf("failed to seek to line start: %w", err)
		}
		s.reader.Reset(s.file)

		line, err := s.reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("failed to read line: %w", err)
		}

		if bytes.Compare(line, searchTerm) >= 0 {
			log.Debugf("Found a candidate: %s", line[:5])
			result = line
			right = mid
		} else {
			left = mid + 1
		}
	}
	return strings.TrimSpace(string(result)), nil
}
