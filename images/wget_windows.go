package images

import (
	"fmt"
	"github.com/mhewedy/vermin/cmd"
)

func wget(url string, file string) (string, error) {
	return cmd.ExecuteP(fmt.Sprintf("(New-Object System.Net.WebClient).DownloadFile('%s', '%s')", url, file))
}
