package base

import (
	"reflect"
	"testing"
)

// extract test data using: https://www.ipaddressguide.com/cidr
func TestNewSubnet(t *testing.T) {
	type args struct {
		anIp    string
		netmask string
	}
	tests := []struct {
		name    string
		args    args
		want    *Subnet
		wantErr bool
	}{
		{
			name:    "test subnet 192.168",
			args:    args{anIp: "192.168.100.4", netmask: "255.255.255.0"},
			want:    &Subnet{start: 3232261120, current: 3232261120, Len: 255},
			wantErr: false,
		},
		{
			name:    "test subnet 172.31 with /20 cidr",
			args:    args{anIp: "172.31.17.69", netmask: "255.255.240.0"},
			want:    &Subnet{start: 2887716864, current: 2887716864, Len: 4095},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSubnet(tt.args.anIp, tt.args.netmask)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSubnet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}
