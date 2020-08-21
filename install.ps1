<#
.SYNOPSIS
Installs vermin the smart virtual machines manager.

Authors: The vermin Maintainers <mhewedy@gmail.com>

.DESCRIPTION
This script builds vermin and ensures that all necesary prerequisites are installed.

.Parameter Version
Specifies a version (ex: 0.75.0, 0.75.0/20190219232208)
#>

param (
    [Alias("v")]
    [string]$Version
)

# based on https://raw.githubusercontent.com/habitat-sh/habitat/master/components/hab/install.ps1

$ErrorActionPreference="stop"

Set-Variable packagesRootUrl -value "https://github.com/mhewedy/vermin/releases/download"

Function Get-File($url, $dst) {
    Write-Host "Downloading $url"
    # Can't use [System.Net.SecurityProtocolType]::Tls12 on older .NET versions
    # Need to use 3072. Un patched older versions of windows will fail even on 3072
    try {
        [System.Net.ServicePointManager]::SecurityProtocol = [Enum]::ToObject([System.Net.SecurityProtocolType], 3072)
    } catch {
        Write-Error "TLS 1.2 is not supported on this operating system. Upgrade or patch your Windows installation."
    }
    $wc = New-Object System.Net.WebClient
    $wc.DownloadFile($url, $dst)
}

Function Get-WorkDir {
    $parent = [System.IO.Path]::GetTempPath()
    [string] $name = [System.Guid]::NewGuid()
    New-Item -ItemType Directory -Path (Join-Path $parent $name)
}

# Downloads the requested archive from packages.chef.io
Function Get-Archive($channel, $version) {

    if(!$version -Or $version -eq "latest") {
        $version=(Get-RedirectedUrl "https://github.com/mhewedy/vermin/releases/latest").Split("/")[7]
    }

    $vermin_url="$packagesRootUrl/${version}/vermin-${version}-windows.zip"
    $sha_url="$vermin_url.sha256sum"
    $hab_dest = (Join-Path ($workdir) "vermin.zip")
    $sha_dest = (Join-Path ($workdir) "vermin.zip.shasum256")

    Get-File $vermin_url $hab_dest
    $result = @{ "zip" = $hab_dest }

    # Note that this will fail on versions less than 0.71.0
    # when we did not upload shasum files to bintray.
    # NOTE: This is left in place because, while we don't ship <0.71.0
    # from s3 today, the intent is to move old releases over
    try {
        Get-File $sha_url $sha_dest
        $result["shasum"] = (Get-Content $sha_dest).Split()[0]
    } catch {
        Write-Warning "No shasum exists for $version. Skipping validation."
    }
    $result
}

function Get-SHA256Converter {
    if($PSVersionTable.PSEdition -eq 'Core') {
        [System.Security.Cryptography.SHA256]::Create()
    } else {
        New-Object -TypeName Security.Cryptography.SHA256Managed
    }
}

Function Get-Sha256($src) {
    $converter = Get-SHA256Converter
    try {
        $bytes = $converter.ComputeHash(($in = (Get-Item $src).OpenRead()))
        return ([System.BitConverter]::ToString($bytes)).Replace("-", "").ToLower()
    } finally {
        # Older .Net versions do not expose Dispose()
        if($PSVersionTable.PSEdition -eq 'Core' -Or ($PSVersionTable.CLRVersion.Major -ge 4)) {
            $converter.Dispose()
        }
        if ($null -ne $in) { $in.Dispose() }
    }
}

Function Assert-Shasum($archive) {
    Write-Host "Verifying the shasum digest matches the downloaded archive"
    $actualShasum = Get-Sha256 $archive.zip
    if($actualShasum -ne $archive.shasum) {
        Write-Error "Checksum '$($archive.shasum)' invalid."
    }
}

Function Install-Vermin {
    $habPath = Join-Path $env:ProgramData Vermin
    if(Test-Path $habPath) { Remove-Item $habPath -Recurse -Force }
    New-Item $habPath -ItemType Directory | Out-Null
    $folder = (Get-ChildItem (Join-Path ($workdir) "vermin.exe"))
    Copy-Item "$($folder.FullName)" $habPath
    $env:PATH = New-PathString -StartingPath $env:PATH -Path $habPath
    $machinePath = [System.Environment]::GetEnvironmentVariable("PATH", "Machine")
    $machinePath = New-PathString -StartingPath $machinePath -Path $habPath
    [System.Environment]::SetEnvironmentVariable("PATH", $machinePath, "Machine")
    $folder.Name.Replace("vermin.exe","")
}

Function New-PathString([string]$StartingPath, [string]$Path) {
    if (-not [string]::IsNullOrEmpty($path)) {
        if (-not [string]::IsNullOrEmpty($StartingPath)) {
            [string[]]$PathCollection = "$path;$StartingPath" -split ';'
            $Path = ($PathCollection |
                    Select-Object -Unique |
                    Where-Object {-not [string]::IsNullOrEmpty($_.trim())} |
                    Where-Object {Test-Path "$_"}
            ) -join ';'
        }
        $path
    } else {
        $StartingPath
    }
}

Function Expand-Zip($zipPath) {
    $dest = $workdir
    try {
        # Works on .Net 4.5 and up (as well as .Net Core)
        # Yes on PS v5 and up we have Expand-Archive but this works on PS v4 too
        [System.Reflection.Assembly]::LoadWithPartialName("System.IO.Compression.FileSystem") | Out-Null
        [System.IO.Compression.ZipFile]::ExtractToDirectory($zipPath, $dest)
    } catch {
        try {
            # Works on all GUI enabled versions. Will fail
            # On Server Core editions
            $shellApplication = New-Object -com shell.application
            $zipPackage = $shellApplication.NameSpace($zipPath)
            $destinationFolder = $shellApplication.NameSpace($dest)
            $destinationFolder.CopyHere($zipPackage.Items())
        } catch{
            Write-Error "Unable to unzip files on this OS"
        }
    }
}

Function Assert-Vermin($ident) {
    Write-Host "Checking installed vermin version $ident"

	$actual = vermin --version
	if (!$actual) {
		Write-Error "Unable to verify vermin was succesfully installed"
	}

}

Function Configure-Vermin() {
    Write-Host "Configuring vermin"

    if(! (Test-Path "$HOME/.vermin/vms")) { New-Item "$HOME/.vermin/vms" -ItemType Directory | Out-Null}
    if(! (Test-Path "$HOME/.vermin/images")) { New-Item "$HOME/.vermin/images" -ItemType Directory | Out-Null}

    wget https://raw.githubusercontent.com/mhewedy/vermin/master/etc/keys/vermin_rsa -o "$HOME/.vermin/vermin_rsa"
    wget https://raw.githubusercontent.com/hashicorp/vagrant/master/keys/vagrant -o "$HOME/.vermin/vagrant_insecure_private_key"
}

Function Configure-Virtualbox() {
    $vboxPath = "C:\Program Files\Oracle\VirtualBox"
    if (!(Test-Path $vboxPath)) {
        Write-Error "Unable to find VirtualBox. Install VirtualBox then re-Run the installation script..."
    }

    $env:PATH = New-PathString -StartingPath $env:PATH -Path $vboxPath
    $machinePath = [System.Environment]::GetEnvironmentVariable("PATH", "Machine")
    $machinePath = New-PathString -StartingPath $machinePath -Path $vboxPath
    [System.Environment]::SetEnvironmentVariable("PATH", $machinePath, "Machine")
}

Function Print-Howto() {

  Write-Host ""
  Write-Host "To list all available images:"
  Write-Host "PS > vermin images"
  Write-Host ""
  Write-Host "To create a VM from an image:"
  Write-Host "PS > vermin create <image>"
  Write-Host ""

}

function Get-RedirectedUrl() {
  param(
    [Parameter(Mandatory = $true, Position = 0)]
    [uri] $url,
    [Parameter(Position = 1)]
    [Microsoft.PowerShell.Commands.WebRequestSession] $session = $null
  )

  $request_url = $url
  $retry = $false

  do {
    try {
      $response = Invoke-WebRequest -UseBasicParsing -Method Head -WebSession $session -Uri $request_url

      if($response.BaseResponse.ResponseUri -ne $null)
      {
        # PowerShell 5
        $result = $response.BaseResponse.ResponseUri.AbsoluteUri
      } elseif ($response.BaseResponse.RequestMessage.RequestUri -ne $null) {
        # PowerShell Core
        $result = $response.BaseResponse.RequestMessage.RequestUri.AbsoluteUri
      }

      $retry = $false
    } catch {
      if(($_.Exception.GetType() -match "HttpResponseException") -and
        ($_.Exception -match "302"))
      {
        $request_url = $_.Exception.Response.Headers.Location.AbsoluteUri
        $retry = $true
      } else {
        throw $_
      }
    }
  } while($retry)

  return $result
}



Write-Host "Installing vermin the smart virtual machines manager"

$workdir = Get-WorkDir
New-Item $workdir -ItemType Directory -Force | Out-Null
try {
    $archive = Get-Archive $version
    if($archive.shasum) {
        Assert-Shasum $archive
    }
    Expand-zip $archive.zip
    $fullIdent = Install-Vermin
    Assert-Vermin $fullIdent
    Configure-Vermin
    Configure-Virtualbox
    Print-Howto

    Write-Host "Installation of vermin program complete."
} finally {
    try { Remove-Item $workdir -Recurse -Force } catch {
        Write-Warning "Unable to delete $workdir"
    }
}
