package hosts

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/dimchansky/utfbom"
)

func TestGetDomainsFromHost(t *testing.T) {
	tests := []struct {
		name         string
		initialData  string
		expectedData []string
	}{
		{
			name:         "get domains within markers",
			initialData:  "0.0.0.0 www.youtube.com\n#focusmode:start\n127.0.0.1 www.instagram.com\n127.0.0.1 www.facebook.com\n#focusmode:end",
			expectedData: []string{"www.instagram.com", "www.facebook.com"},
		},
		{
			name:         "no domains within markers",
			initialData:  "127.0.0.1 www.google.com\n#focusmode:start\n#focusmode:end",
			expectedData: []string{},
		},
		{
			name:         "no markers present",
			initialData:  "127.0.0.1 www.google.com",
			expectedData: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempHostsFile, cleanFile := createTempFile(t, tt.initialData)
			defer cleanFile()

			store := &HostsStore{hostsFile: tempHostsFile}
			defer store.Close()

			got, err := store.GetDomainsFromHost()
			if err != nil {
				t.Fatalf("GetDomainsFromHost() error: %v", err)
			}

			assertDomains(t, got, tt.expectedData)
		})
	}
}

func TestAddDomainsToHost(t *testing.T) {
	tests := []struct {
		name             string
		initialData      string
		expectedFileData string
		expectedDomains  []string
	}{
		{
			name:             "adding domains to empty hosts file",
			initialData:      ``,
			expectedDomains:  []string{"www.instagram.com", "www.youtube.com"},
			expectedFileData: "#focusmode:start\n127.0.0.1 www.instagram.com\n127.0.0.1 www.youtube.com\n#focusmode:end",
		},
		{
			name:             "overwriting domains to existing hosts file",
			initialData:      "#focusmode:start\n127.0.0.1 www.youtube.com\n127.0.0.1 www.instagram.com\n127.0.0.1 www.github.com\n#focusmode:end",
			expectedDomains:  []string{"www.youtube.com", "www.facebook.com", "www.amazon.com"},
			expectedFileData: "#focusmode:start\n127.0.0.1 www.youtube.com\n127.0.0.1 www.facebook.com\n127.0.0.1 www.amazon.com\n#focusmode:end",
		},
		{
			name:             "overwriting domains to hosts file without deleting contents",
			initialData:      "# comment\n#comment\n#focusmode:start\n127.0.0.1 www.youtube.com\n127.0.0.1 www.instagram.com\n127.0.0.1 www.github.com\n#focusmode:end",
			expectedDomains:  []string{"www.youtube.com", "www.facebook.com", "www.amazon.com"},
			expectedFileData: "# comment\n#comment\n#focusmode:start\n127.0.0.1 www.youtube.com\n127.0.0.1 www.facebook.com\n127.0.0.1 www.amazon.com\n#focusmode:end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempHostsFile, cleanFile := createTempFile(t, tt.initialData)
			defer cleanFile()

			store := &HostsStore{hostsFile: tempHostsFile}
			defer store.Close()

			// add domains to hosts file
			err := store.AddDomainsToHost(tt.expectedDomains)
			if err != nil {
				t.Fatalf("AddDomainsToHost() error: %v", err)
			}

			got, err := store.GetDomainsFromHost()
			if err != nil {
				t.Fatalf("GetDomainsFromHost() error: %v", err)
			}
			assertDomains(t, got, tt.expectedDomains)
			assertFileData(t, tempHostsFile, tt.expectedFileData)
		})
	}
}

func TestCleanDomainsFromHost(t *testing.T) {
	tests := []struct {
		name             string
		initialData      string
		expectedFileData string
		expectedDomains  []string
	}{
		{
			name:             "delete all domains in host",
			initialData:      "# comment\n#comment\n#focusmode:start\n127.0.0.1 www.youtube.com\n127.0.0.1 www.instagram.com\n127.0.0.1 www.github.com\n#focusmode:end",
			expectedFileData: "# comment\n#comment\n#focusmode:start\n#focusmode:end",
			expectedDomains:  []string{},
		},
		{
			name:             "empty hosts markers",
			initialData:      "127.0.0.1 www.instagram.com\n# 0.0.0.0 www.docker.com\n#focusmode:start\n#focusmode:end",
			expectedFileData: "127.0.0.1 www.instagram.com\n# 0.0.0.0 www.docker.com\n#focusmode:start\n#focusmode:end",
			expectedDomains:  []string{},
		},
		{
			name:             "empty hosts markers with stub domains",
			initialData:      "127.0.0.1 www.instagram.com\n# 0.0.0.0 www.docker.com\n#focusmode:start\n0.0.0.0 www.youtube.com\n#focusmode:end",
			expectedFileData: "127.0.0.1 www.instagram.com\n# 0.0.0.0 www.docker.com\n#focusmode:start\n#focusmode:end",
			expectedDomains:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempHostsFile, cleanFile := createTempFile(t, tt.initialData)
			defer cleanFile()

			store := &HostsStore{hostsFile: tempHostsFile}
			defer store.Close()

			err := store.CleanDomains()
			if err != nil {
				t.Fatalf("CleanDomains() error: %v", err)
			}

			got, err := store.GetDomainsFromHost()
			if err != nil {
				t.Fatalf("GetDomainsFromHost() error: %v", err)
			}
			assertDomains(t, got, tt.expectedDomains)
			assertFileData(t, tempHostsFile, tt.expectedFileData)
		})
	}
}

func TestDeleteDomainFromHost(t *testing.T) {
	tests := []struct {
		name             string
		domainToDelete   string
		initialData      string
		expectedFileData string
		expectedDomains  []string
		expectedErr      error
	}{
		{
			name:             "empty domain list",
			domainToDelete:   "www.instagram.com",
			initialData:      ``,
			expectedFileData: ``,
			expectedDomains:  []string{},
			expectedErr:      errNoDomains,
		},
		{
			name:             "remove instagram.com",
			domainToDelete:   "www.instagram.com",
			initialData:      "#focusmode:start\n127.0.0.1 www.instagram.com\n127.0.0.1 www.facebook.com\n127.0.0.1 www.snapchat.com\n#focusmode:end",
			expectedFileData: "#focusmode:start\n127.0.0.1 www.facebook.com\n127.0.0.1 www.snapchat.com\n#focusmode:end",
			expectedDomains:  []string{"www.facebook.com", "www.snapchat.com"},
			expectedErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempHostsFile, cleanFile := createTempFile(t, tt.initialData)
			defer cleanFile()

			store := &HostsStore{hostsFile: tempHostsFile}
			defer store.Close()

			err := store.DeleteDomainFromHost(tt.domainToDelete)
			assertError(t, err, tt.expectedErr)

			d, err := store.GetDomainsFromHost()
			if err != nil {
				t.Errorf("GetDomainsFromHost() error: %v", err)
			}

			assertDomains(t, d, tt.expectedDomains)
			assertFileData(t, tempHostsFile, tt.expectedFileData)
		})
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("Wrong error type got: %v, want: %v", got, want)
	}
}

func assertDomains(t testing.TB, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want = %v", got, want)
	}
}

func assertFileData(t testing.TB, f *os.File, want string) {
	t.Helper()
	got, err := extractFileData(f)
	if err != nil {
		t.Errorf("Error extracting file data %v", err)
	}

	if got != want {
		t.Errorf("\ngot = %q\nwant = %q", got, want)
	}
}

// returns file contents as a string
func extractFileData(f *os.File) (string, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	reader := utfbom.SkipOnly(bufio.NewReader(f))
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpFile.Write([]byte(initialData))
	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return tmpFile, removeFile
}
