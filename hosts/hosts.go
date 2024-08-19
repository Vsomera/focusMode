package hosts

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/dimchansky/utfbom"
)

// extracts data from hosts file
func GetDomainsFromHost() ([]string, error) {
	domains := []string{}
	hostsFile, err := os.OpenFile(hostsPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return domains, err
	}

	defer hostsFile.Close()

	data, err := io.ReadAll(hostsFile)
	if err != nil {
		return domains, err
	}
	// remove BOM (Byte Order Mark) if present at the beginning of the file
	data, err = io.ReadAll(utfbom.SkipOnly(bytes.NewReader(data)))
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

	re := regexp.MustCompile(`^\s*\d{1,3}(?:\.\d{1,3}){3}\s+`) // regex to match ip addresses
	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {

		line := strings.ToLower(strings.TrimSpace(scanner.Text()))

		// skip comments or empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// remove ip address and trailing spaces
		line = re.ReplaceAllString(line, "")

		// remove inline comments
		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		if line != "" {
			domains = append(domains, line)
		}

	}

	if err := scanner.Err(); err != nil {
		return domains, err
	}

	return domains, nil
}
