package hosts

import (
	"os"
	"reflect"
	"testing"
)

func TestHostsFileStore(t *testing.T) {
	t.Run("get domains between markers", func(t *testing.T) {
		tempHostsFile, cleanFile := createTempFile(t, `
		# comment
		# comment
		0.0.0.0 www.youtube.com
		#focusmode:start
		127.0.0.1 www.instagram.com
		0.0.0.0 www.facebook.com
		#focusmode:end
		`)
		defer cleanFile()

		store := &HostsStore{hostsFile: tempHostsFile}
		got, _ := store.GetDomainsFromHost()
		want := []string{"www.instagram.com", "www.facebook.com"}

		assertDomains(t, got, want)
	})
	t.Run("domains outside of markers", func(t *testing.T) {
		tempHostsFile, cleanFile := createTempFile(t, `
		# comment
		# comment
		0.0.0.0 www.youtube.com
		127.0.0.1 www.instagram.com
		0.0.0.0 www.facebook.com
		# comment
		`)
		defer cleanFile()

		store := &HostsStore{hostsFile: tempHostsFile}
		got, _ := store.GetDomainsFromHost()
		want := []string{}

		assertDomains(t, got, want)
	})
	t.Run("adding domains to hosts file", func(t *testing.T) {
		tempHostsFile, cleanFile := createTempFile(t, `
		# comment
		# comment
		`)
		defer cleanFile()
		store := &HostsStore{hostsFile: tempHostsFile}

		got := []string{"www.instagram.com", "www.youtube.com", "www.facebook.com"}
		store.AddDomainsToHost(got)
		want, _ := store.GetDomainsFromHost()

		assertDomains(t, got, want)
	})
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
