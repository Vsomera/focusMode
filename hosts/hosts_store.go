package hosts

import (
	"bufio"
	"errors"
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

var (
	errNoDomains      = errors.New("blacklist is empty, no domains present")
	errDomainNotFound = errors.New("domain not found in blacklist")
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

	return extractDomainsFromData(string(data))
}

func (hs *HostsStore) CleanDomains() error {
	if _, err := hs.hostsFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	reader := utfbom.SkipOnly(bufio.NewReader(hs.hostsFile))
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	cleanData, err := updateHostData(string(data), []string{})
	if err != nil {
		return err
	}

	return hs.writeDataToHostFile(cleanData)
}

func (hs *HostsStore) DeleteDomainFromHost(domain string) error {
	currDomains, err := hs.GetDomainsFromHost()
	if err != nil {
		return err
	}
	if len(currDomains) < 1 {
		return errNoDomains
	}

	domainDeleted := false
	newDomains := []string{}
	for _, d := range currDomains {
		if d == domain {
			domainDeleted = true
			continue
		}
		newDomains = append(newDomains, d)
	}
	if !domainDeleted {
		return errDomainNotFound
	}

	return hs.AddDomainsToHost(newDomains)
}

func (hs *HostsStore) AddDomainsToHost(domains []string) error {
	if _, err := hs.hostsFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	reader := utfbom.SkipOnly(bufio.NewReader(hs.hostsFile))
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	newData, err := updateHostData(string(data), domains)
	if err != nil {
		return err
	}

	return hs.writeDataToHostFile(newData)
}

func (hs *HostsStore) writeDataToHostFile(data string) error {

	if _, err := hs.hostsFile.Seek(0, 0); err != nil {
		return err
	}
	if _, err := hs.hostsFile.WriteString(data); err != nil {
		return err
	}
	if err := hs.hostsFile.Truncate(int64(len(data))); err != nil {
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

	return domains, scanner.Err()
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
