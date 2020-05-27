package vms

import (
	"reflect"
	"testing"
)

func Test_filters_apply(t *testing.T) {
	type args struct {
		list vmInfoList
	}
	tests := []struct {
		name string
		f    filters
		args args
		want vmInfoList
	}{
		{
			name: "test 3 entries in the list, and 2 of them matches on tags",
			f:    []filter{{name: "tags", value: "k8s"}},
			args: args{list: vmInfoList{vmInfo{tags: "k8s-master"}, vmInfo{tags: "k8s-slave"}, vmInfo{tags: "redis"}}},
			want: vmInfoList{vmInfo{tags: "k8s-master"}, vmInfo{tags: "k8s-slave"}},
		},
		{
			name: "when no filter matches, return empty list",
			f:    []filter{{name: "name", value: "not found"}},
			args: args{list: vmInfoList{vmInfo{tags: "k8s-master"}, vmInfo{tags: "k8s-slave"}, vmInfo{tags: "redis"}}},
			want: vmInfoList{},
		},
		{
			name: "when no filter provided return the list as is",
			f:    []filter{},
			args: args{list: vmInfoList{vmInfo{tags: "k8s-master"}, vmInfo{tags: "k8s-slave"}, vmInfo{tags: "redis"}}},
			want: vmInfoList{vmInfo{tags: "k8s-master"}, vmInfo{tags: "k8s-slave"}, vmInfo{tags: "redis"}},
		},
		{
			name: "test two filters provided - AND ops will be applied",
			f:    []filter{{name: "tags", value: "k8s"}, {name: "image", value: "ubuntu"}, {name: "name", value: "vm"}},
			args: args{list: vmInfoList{vmInfo{tags: "k8s-master"}, vmInfo{tags: "k8s-slave", image: "ubuntu/focal", name: "vm_01"}}},
			want: vmInfoList{vmInfo{tags: "k8s-slave", image: "ubuntu/focal", name: "vm_01"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.apply(tt.args.list); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
