package scp

import (
	"reflect"
	"testing"
)

func Test_toVmPath(t *testing.T) {
	type args struct {
		srcDest string
	}
	tests := []struct {
		name  string
		args  args
		want  vmPath
		want1 bool
	}{
		{
			name: "test case where vmname is there",
			args: struct{ srcDest string }{"vm_01:/path/to/file"},
			want: vmPath{"vm_01", "/path/to/file"}, want1: true,
		},
		{
			name: "test case where vmname is not there there",
			args: struct{ srcDest string }{"/path/to/file"},
			want: vmPath{}, want1: false,
		},
		{
			name: "test case where vmname is there there without separator",
			args: struct{ srcDest string }{"vm_01/path/to/file"},
			want: vmPath{}, want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := toVmPath(tt.args.srcDest)
			if got1 != tt.want1 {
				t.Errorf("toVmPath() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("toVmPath() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
