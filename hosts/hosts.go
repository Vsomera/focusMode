package hosts

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/dimchansky/utfbom"
)

const (
	CommentStart = "#focusmode:start"
	CommentEnd   = "#focusmode:end"
)

// extracts data from hosts file
func GetDomainsFromHost() ([]string, error) {
	domains := []string{}
	hostsFile, err := os.OpenFile(hostsPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return domains, err
	}
	defer hostsFile.Close()

	// skip BOM if present
	reader := utfbom.SkipOnly(bufio.NewReader(hostsFile))
	data, err := io.ReadAll(reader)
	if err != nil {
		return domains, err
	}

	domains, err = extractDomainsFromData(string(data))
	if err != nil {
		return domains, err
	}

	return domains, nil
}

// given data from bytes (which is converted from string) extract the domains
func extractDomainsFromData(data string) ([]string, error) {
	domains := []string{}
	inFocus := false

	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {

		line := strings.ToLower(strings.TrimSpace(scanner.Text()))

		// check if domains are in focus start and end markers
		if strings.Contains(line, CommentStart) {
			inFocus = true
			continue
		}
		if strings.Contains(line, CommentEnd) {
			inFocus = false
			continue
		}

		// only append domains if domains are within focus markers
		if inFocus {
			// skip comments or empty lines
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			// remove ip address and trailing spaces
			re := regexp.MustCompile(`^\s*\d{1,3}(?:\.\d{1,3}){3}\s+`)
			line = re.ReplaceAllString(line, "")

			// remove inline comments
			if idx := strings.Index(line, "#"); idx != -1 {
				line = strings.TrimSpace(line[:idx])
			}

			if line != "" {
				domains = append(domains, line)
			}
		}

	}

	if err := scanner.Err(); err != nil {
		return domains, err
	}

	return domains, nil
}
