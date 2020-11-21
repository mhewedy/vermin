package command

import (
	"github.com/mhewedy/vermin/hypervisor"
	"github.com/mhewedy/vermin/images"
	"github.com/spf13/cobra"
	"strings"
)

func listRunningVms(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	list, _ := hypervisor.List(false)

	var completions []string
	for _, comp := range list {
		if strings.HasPrefix(comp, toComplete) {
			completions = append(completions, comp)
		}
	}
	return completions, cobra.ShellCompDirectiveDefault
}

func listAllVms(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	list, _ := hypervisor.List(true)

	var completions []string
	for _, comp := range list {
		if strings.HasPrefix(comp, toComplete) {
			completions = append(completions, comp)
		}
	}
	return completions, cobra.ShellCompDirectiveDefault
}

func listStoppedVms(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	for _, comp := range getStoppedVms() {
		if strings.HasPrefix(comp, toComplete) {
			completions = append(completions, comp)
		}
	}
	return completions, cobra.ShellCompDirectiveDefault
}

func listImages(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	list, _ := images.List()

	var completions []string
	for _, comp := range list {
		if strings.HasPrefix(comp, toComplete) {
			completions = append(completions, comp)
		}
	}
	return completions, cobra.ShellCompDirectiveDefault
}

func getStoppedVms() []string {
	stopped := make([]string, 0)
	all, _ := hypervisor.List(true)
	running, _ := hypervisor.List(false)

	for i := range all {
		vm := all[i]
		if !contains(running, vm) {
			stopped = append(stopped, vm)
		}
	}
	return stopped
}

func contains(a []string, s string) bool {
	for i := range a {
		if a[i] == s {
			return true
		}
	}
	return false
}
