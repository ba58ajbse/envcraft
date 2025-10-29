package fs

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func ReadLines(filePath string) ([]string, error) {
	envFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer envFile.Close()

	lines := []string{}
	reader := bufio.NewReader(envFile)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				if len(line) > 0 || len(lines) == 0 {
					lines = append(lines, line)
				}
				break
			}
			return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func WriteLines(filePath string, lines []string) error {
	out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %s for writing: %w", filePath, err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	for _, line := range lines {
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("error writing to file %s: %w", filePath, err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}
