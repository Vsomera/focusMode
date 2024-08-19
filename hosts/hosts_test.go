package hosts

import (
	"reflect"
	"testing"
)

func TestExtractDomainsFromData(t *testing.T) {

	tests := []struct {
		name        string
		data        string
		wantDomains []string
		wantErr     bool
	}{
		{
			name:        "no domains",
			data:        "",
			wantDomains: []string{},
			wantErr:     false,
		},
		{
			name: "extracting domains with comments",
			data: `# Copyright (c) 1993-2009 Microsoft Corp.
			127.0.0.1 example.host
			# comment
			127.0.0.1 example.host
			# comment`,
			wantDomains: []string{
				"example.host",
				"example.host",
			},
			wantErr: false,
		},
		{
			name: "test file with no comments",
			data: `0.0.0.0 www.instagram.com
			0.0.0.0 www.facebook.com
			192.168.1.83 host.docker.internal
			192.168.1.83 gateway.docker.internal
			127.0.0.1 kubernetes.docker.internal`,
			wantDomains: []string{
				"www.instagram.com",
				"www.facebook.com",
				"host.docker.internal",
				"gateway.docker.internal",
				"kubernetes.docker.internal",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotDomains, err := extractDomainsFromData(tt.data)

			if !reflect.DeepEqual(gotDomains, tt.wantDomains) {
				t.Errorf("gotDomains = %v, want = %v", gotDomains, tt.wantDomains)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error %v, wantError %v", err, tt.wantErr)
			}

		})
	}

}
