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
			name: "get domains within markers",
			initialData: `
			0.0.0.0 www.youtube.com
			#focusmode:start
			127.0.0.1 www.instagram.com
			127.0.0.1 www.facebook.com
			#focusmode:end
			`,
			expectedData: []string{"www.instagram.com", "www.facebook.com"},
		},
		{
			name: "no domains within markers",
			initialData: `
			127.0.0.1 www.google.com
			#focusmode:start
			#focusmode:end
			`,
			expectedData: []string{},
		},
		{
			name: "no markers present",
			initialData: `
        	127.0.0.1 www.google.com
			`,
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
			name:            "adding domains to empty hosts file",
			initialData:     ``,
			expectedDomains: []string{"www.instagram.com", "www.youtube.com"},
			expectedFileData: `#focusmode:start
127.0.0.1 www.instagram.com
127.0.0.1 www.youtube.com
#focusmode:end`,
		},
		{
			name: "overwriting domains to existing hosts file",
			initialData: `#focusmode:start
127.0.0.1 www.youtube.com
127.0.0.1 www.instagram.com
127.0.0.1 www.github.com
#focusmode:end`,
			expectedDomains: []string{"www.youtube.com", "www.facebook.com", "www.amazon.com"},
			expectedFileData: `#focusmode:start
127.0.0.1 www.youtube.com
127.0.0.1 www.facebook.com
127.0.0.1 www.amazon.com
#focusmode:end`,
		},
		{
			name: "overwriting domains to hosts file without deleting contents",
			initialData: `# comment
#comment
#focusmode:start
127.0.0.1 www.youtube.com
127.0.0.1 www.instagram.com
127.0.0.1 www.github.com
#focusmode:end`,
			expectedDomains: []string{"www.youtube.com", "www.facebook.com", "www.amazon.com"},
			expectedFileData: `# comment
#comment
#focusmode:start
127.0.0.1 www.youtube.com
127.0.0.1 www.facebook.com
127.0.0.1 www.amazon.com
#focusmode:end`,
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
