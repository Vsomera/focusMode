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
	defaultIpAddress = "127.0.0.1"
	commentStart     = "#focusmode:start"
	commentEnd       = "#focusmode:end"
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

// given a slice of hosts, appends host to hosts file
func AddDomainsToHost(domains []string) error {
	hostsFile, err := os.OpenFile(hostsPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer hostsFile.Close()

	// skip BOM if present
	reader := utfbom.SkipOnly(bufio.NewReader(hostsFile))
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	newData, err := updateHostData(string(data), domains)
	if err != nil {
		return err
	}

	if _, err := hostsFile.Seek(0, 0); err != nil {
		return err
	}
	if _, err := hostsFile.WriteString(newData); err != nil {
		return err
	}
	if err := hostsFile.Truncate(int64(len(newData))); err != nil {
		return err
	}

	return nil
}

// given data from bytes (which is converted from string) extract the domains
func extractDomainsFromData(data string) ([]string, error) {
	domains := []string{}
	inFocus := false

	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {

		line := strings.ToLower(strings.TrimSpace(scanner.Text()))

		// check if domains are in focus start and end markers
		if line == commentStart {
			inFocus = true
			continue
		}
		if line == commentEnd {
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

func updateHostData(initData string, domains []string) (string, error) {
	inFocus := false
	newData := ""

	scanner := bufio.NewScanner(strings.NewReader(initData))
	for scanner.Scan() {

		line := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if line == commentStart {
			inFocus = true
			continue
		} else if line == commentEnd {
			inFocus = false
			continue
		}

		if !inFocus {
			newData += line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		return initData, err
	}

	newData += commentStart + "\n"
	for _, d := range domains {
		newData += defaultIpAddress + " " + d + "\n"
	}
	newData += commentEnd + "\n"

	return newData, nil
}
