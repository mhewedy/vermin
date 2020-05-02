// +build !windows

package images

import "github.com/mhewedy/vermin/cmd"

func wget(url string, file string) (string, error) {
	return cmd.Execute("wget", "-O", file, url)
}

func wgetP(title string, url string, file string) (string, error) {
	return cmd.ExecuteP(title, "wget", "-O", file, url)
}
