package hosts

import (
	"bufio"
	"io"
	"log"
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

type HostsStore struct {
	hostsFile *os.File
}

func NewHostsStore() *HostsStore {
	hostsFile, err := os.OpenFile(hostsPath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &HostsStore{hostsFile: hostsFile}
}

func (hs *HostsStore) GetDomainsFromHost() ([]string, error) {
	domains := []string{}

	// reset file pointer
	if _, err := hs.hostsFile.Seek(0, io.SeekStart); err != nil {
		return domains, err
	}

	// skip BOM if present
	reader := utfbom.SkipOnly(bufio.NewReader(hs.hostsFile))
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

func (hs *HostsStore) CleanDomains() error {
	// reset file pointer
	if _, err := hs.hostsFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// skip BOM if present
	reader := utfbom.SkipOnly(bufio.NewReader(hs.hostsFile))
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	// overwrite hosts with empty domains
	newData, err := updateHostData(string(data), []string{})
	if err != nil {
		return err
	}

	if _, err := hs.hostsFile.Seek(0, 0); err != nil {
		return err
	}
	if _, err := hs.hostsFile.WriteString(newData); err != nil {
		return err
	}
	if err := hs.hostsFile.Truncate(int64(len(newData))); err != nil {
		return err
	}

	return nil
}

func (hs *HostsStore) AddDomainsToHost(domains []string) error {
	// reset file pointer
	if _, err := hs.hostsFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// skip BOM if present
	reader := utfbom.SkipOnly(bufio.NewReader(hs.hostsFile))
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	newData, err := updateHostData(string(data), domains)
	if err != nil {
		return err
	}

	if _, err := hs.hostsFile.Seek(0, 0); err != nil {
		return err
	}
	if _, err := hs.hostsFile.WriteString(newData); err != nil {
		return err
	}
	if err := hs.hostsFile.Truncate(int64(len(newData))); err != nil {
		return err
	}

	return nil
}

// close file
func (hs *HostsStore) Close() error {
	if hs.hostsFile != nil {
		err := hs.hostsFile.Close()
		if err != nil {
			return err
		}
		hs.hostsFile = nil
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

		if line == commentStart {
			inFocus = true
			continue
		}
		if line == commentEnd {
			inFocus = false
			continue
		}

		if inFocus {
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			re := regexp.MustCompile(`^\s*\d{1,3}(?:\.\d{1,3}){3}\s+`)
			line = re.ReplaceAllString(line, "")

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

// updates and overwrites domains within start-end markers
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
