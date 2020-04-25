package vms

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPortForwardArgs(t *testing.T) {
	type args struct {
		ports string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "<vm port>",
			args:    struct{ ports string }{ports: "3000"},
			want:    []string{"0.0.0.0:3000:localhost:3000"},
			wantErr: false},
		{name: "<vm port>:<local port>",
			args:    struct{ ports string }{ports: "3000:4000"},
			want:    []string{"0.0.0.0:4000:localhost:3000"},
			wantErr: false},
		{name: "<vm port>:<local port> <vmport>",
			args:    struct{ ports string }{ports: "3000:4000 8080"},
			want:    []string{"0.0.0.0:4000:localhost:3000", "0.0.0.0:8080:localhost:8080"},
			wantErr: false},
		{name: "<vmport> <vm port>:<local port>",
			args:    struct{ ports string }{ports: "8080 3000:4000"},
			want:    []string{"0.0.0.0:8080:localhost:8080", "0.0.0.0:4000:localhost:3000"},
			wantErr: false},
		{name: "<vmport> <vm port> <vm port>",
			args:    struct{ ports string }{ports: "3000 4000 5000"},
			want:    []string{"0.0.0.0:3000:localhost:3000", "0.0.0.0:4000:localhost:4000", "0.0.0.0:5000:localhost:5000"},
			wantErr: false},
		{name: "<vmport> <vm port1>-<vm port2>",
			args:    struct{ ports string }{ports: "3000 4000-4001"},
			want:    []string{"0.0.0.0:3000:localhost:3000", "0.0.0.0:4000:localhost:4000", "0.0.0.0:4001:localhost:4001"},
			wantErr: false},
		{name: "<vm port1>-<vm port2>:<local port1>-<local port2>",
			args:    struct{ ports string }{ports: "4000-4001:8080-8081"},
			want:    []string{"0.0.0.0:8080:localhost:4000", "0.0.0.0:8081:localhost:4001"},
			wantErr: false},
		{name: "<vm port1>-<vm port2>:<local port1>-<local port2> <vmport>",
			args:    struct{ ports string }{ports: "4000-4001:8080-8081 9000"},
			want:    []string{"0.0.0.0:8080:localhost:4000", "0.0.0.0:8081:localhost:4001", "0.0.0.0:9000:localhost:9000"},
			wantErr: false},
		{name: "invalid <vm port1>-<vm port2>:<local port1>-<local port2>",
			args:    struct{ ports string }{ports: "4000-4001:8080-8082"},
			want:    nil,
			wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapPorts(tt.args.ports)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapPorts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			want := tt.want

			if !assert.ElementsMatch(t, got, tt.want) {
				t.Errorf("mapPorts() got = %v, want %v, %v %v", got, want, len(got), len(want))
			}
		})
	}
}
