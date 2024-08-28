package hosts

import (
	"os"
	"reflect"
	"testing"
)

type HostTest []struct {
	name         string
	initialData  string
	expectedData []string
}

func TestGetDomainsFromHost(t *testing.T) {
	tests := HostTest{
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
	tests := HostTest{
		{
			name:         "adding domains to empty hosts file",
			initialData:  ``,
			expectedData: []string{"www.instagram.com", "www.youtube.com"},
		},
		{
			name: "overwriting domains to existing hosts file",
			initialData: `
			#focusmode:start
			www.youtube.com
			www.instagram.com
			www.github.com
			#focusmode:end
			`,
			expectedData: []string{"www.facebook.com", "www.amazon.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempHostsFile, cleanFile := createTempFile(t, tt.initialData)
			defer cleanFile()

			store := &HostsStore{hostsFile: tempHostsFile}
			defer store.Close()

			// add domains to hosts file
			store.AddDomainsToHost(tt.expectedData)

			got, err := store.GetDomainsFromHost()
			if err != nil {
				t.Fatalf("GetDomainsFromHost() error: %v", err)
			}
			assertDomains(t, got, tt.expectedData)
		})
	}
}

func assertDomains(t testing.TB, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %v, want = %v", got, want)
	}
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
