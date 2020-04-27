#!/bin/bash
# based on: https://raw.githubusercontent.com/habitat-sh/habitat/master/components/hab/install.sh
set -eou pipefail

# If the variable `$DEBUG` is set, then print the shell commands as we execute.
if [ -n "${DEBUG:-}" ]; then set -x; fi

readonly pcio_root="https://github.com/mhewedy/vermin/releases/download"

main() {
  # Set an empty version variable, signaling we want the latest release
  version=""

  # Parse command line flags and options.
  while getopts "hv:t:" opt; do
    case "${opt}" in
      h)
        print_help
        exit 0
        ;;
      v)
        version="${OPTARG}"
        ;;
      t)
        target="${OPTARG}"
        ;;
      \?)
        echo "" >&2
        print_help >&2
        exit_with "Invalid option" 1
        ;;
    esac
  done

  info "Installing vermin"
  create_workdir
  get_platform
  validate_target
  download_archive "$version" "$target"
  verify_archive
  extract_archive
  install_vermin
  configure_vermin
  info "Installation of vermin complete."
  info "********************************"
  info "You can start creating VMs right away using 'vermin images' then 'vermin create <image>'"
  info "********************************"
}

print_help() {
  need_cmd cat
  need_cmd basename

  local _cmd
  _cmd="$(basename "${0}")"
  cat <<USAGE
${_cmd}

Authors: The Vermin Maintainers <mhewedy@gmail.com>

Installs the Vermin 'vermin' program.

USAGE:
    ${_cmd} [FLAGS]

FLAGS:
    -h    Prints help information
    -v    Specifies a version (ex: 0.15.0, 0.15.0/20161222215311)
    -t    Specifies the ActiveTarget of the 'vermin' program to download.
            [values: x86_64-linux, x86_64-linux-kernel2] [default: x86_64-linux]
            This option is only valid on Linux platforms

ENVIRONMENT VARIABLES:
     SSL_CERT_FILE   allows you to verify against a custom cert such as one
                     generated from a corporate firewall

USAGE
}

create_workdir() {
  need_cmd mktemp
  need_cmd rm
  need_cmd mkdir

  if [ -n "${TMPDIR:-}" ]; then
    local _tmp="${TMPDIR}"
  elif [ -d /var/tmp ]; then
    local _tmp=/var/tmp
  else
    local _tmp=/tmp
  fi

  workdir="$(mktemp -d -p "$_tmp" 2> /dev/null || mktemp -d "${_tmp}/vermin.XXXX")"
  # Add a trap to clean up any interrupted file downloads
  # shellcheck disable=SC2154
  trap 'code=$?; rm -rf $workdir; exit $code' INT TERM EXIT
  cd "${workdir}"
}

get_platform() {
  need_cmd uname
  need_cmd tr

  local _ostype
  _ostype="$(uname -s)"

  case "${_ostype}" in
    Darwin|Linux)
      sys="$(uname -s | tr '[:upper:]' '[:lower:]')"
      ;;
    *)
      exit_with "Unrecognized OS type when determining platform: ${_ostype}" 2
      ;;
  esac

  case "${sys}" in
    darwin)
      need_cmd shasum

      ext=tar.gz
      shasum_cmd="shasum -a 256"
      ;;
    linux)
      need_cmd sha256sum

      ext=tar.gz
      shasum_cmd="sha256sum"
      ;;
    *)
      exit_with "Unrecognized sys type when determining platform: ${sys}" 3
      ;;
  esac

  if [ -z "${target:-}" ]; then
    target="${sys}"
  fi
}

# Validate the CLI Target requested.  In most cases ${arch}-${sys}
# for the current system is the only valid Target.  In the case of
# x86_64-linux systems we also need to support the x86_64-linux-kernel2
# Target. Creates an array of valid Targets for the current system,
# adding any valid alternate Targets, and checks if the requested
# Target is present in the array.
validate_target() {
  local valid_targets=("${sys}")
  case "${sys}" in
   linux)
    valid_targets+=("x86_64-linux-kernel2")
    ;;
  esac

  if ! (_array_contains "${target}" "${valid_targets[@]}") ; then
    local _vts
    printf -v _vts "%s, " "${valid_targets[@]}"
    _e="${target} is not a valid target for this system. Please specify one of: [${_vts%, }]"
    exit_with "$_e" 7
  fi
}

download_archive() {
  need_cmd mv
  need_cmd curl

  local _version="${1:-latest}"
  local -r _target="${2:?}"
  local url

  if [ "$_version" == "latest" ]; then
    _version="$(curl -s https://github.com/mhewedy/vermin/releases/latest -I | grep 'location: ' | cut -d '/' -f8 | tr -d '\r' | tr -d '\n')"
  fi

  url="${pcio_root}/${_version}/vermin-${_version}.${_target}.${ext}"

  dl_file "${url}" "${workdir}/vermin-${_version}.${_target}.${ext}"
  dl_file "${url}.sha256sum" "${workdir}/vermin-${_version}.${_target}.${ext}.sha256sum"

  archive="vermin-${_version}.${_target}.${ext}"
  sha_file="vermin-${_version}.${_target}.${ext}.sha256sum"
}

verify_archive() {
  info "Verifying the shasum digest matches the downloaded archive"
  ${shasum_cmd} -c "${sha_file}"
}

extract_archive() {
  need_cmd sed

  info "Extracting ${archive}"
  case "${ext}" in
    tar.gz)
      need_cmd tar

      archive_dir="${archive%.tar.gz}"
      mkdir "${archive_dir}"
      tar xzf "${archive}" -C "${archive_dir}"
      ;;
    *)
      exit_with "Unrecognized file extension when extracting: ${ext}" 4
      ;;
  esac
}

install_vermin() {
  case "${sys}" in
    darwin|linux)
      need_cmd mkdir
      need_cmd install

      info "Installing vermin into /usr/local/bin"
      mkdir -pv /usr/local/bin
      install -v "${archive_dir}"/vermin /usr/local/bin/vermin
      ;;
    *)
      exit_with "Unrecognized sys when installing: ${sys}" 5
      ;;
  esac
}

configure_vermin() {
  case "${sys}" in
    darwin|linux)
      need_cmd mkdir
      need_cmd curl
      need_cmd chmod
      need_cmd vboxmanage

      info "Configuring vermin"
      mkdir -pv "$HOME/.vermin/vms"
      mkdir -pv "$HOME/.vermin/images"
      curl -s https://raw.githubusercontent.com/mhewedy/vermin/master/etc/keys/vermin_rsa > "$HOME/.vermin/vermin_rsa"
      chmod 0600 "$HOME/.vermin/vermin_rsa"
      ;;
    *)
      exit_with "Unrecognized sys when installing: ${sys}" 5
      ;;
  esac
}

need_cmd() {
  if ! command -v "$1" > /dev/null 2>&1; then
    exit_with "Required command '$1' not found on PATH" 127
  fi
}

info() {
  echo "--> vermin-install: $1"
}

warn() {
  echo "xxx vermin-install: $1" >&2
}

exit_with() {
  warn "$1"
  exit "${2:-10}"
}

_array_contains() {
  local e
  for e in "${@:2}"; do
    if [[ "$e" == "$1" ]]; then
      return 0
    fi
  done
  return 1
}

dl_file() {
  local _url="${1}"
  local _dst="${2}"
  local _code
  local _wget_extra_args=""
  local _curl_extra_args=""

  # Attempt to download with wget, if found. If successful, quick return
  if command -v wget > /dev/null; then
    info "Downloading via wget: ${_url}"

    if [ -n "${SSL_CERT_FILE:-}" ]; then
      wget ${_wget_extra_args:+"--ca-certificate=${SSL_CERT_FILE}"} -q -O "${_dst}" "${_url}"
    else
      wget -q -O "${_dst}" "${_url}"
    fi

    _code="$?"

    if [ $_code -eq 0 ]; then
      return 0
    else
      local _e="wget failed to download file, perhaps wget doesn't have"
      _e="$_e SSL support and/or no CA certificates are present?"
      warn "$_e"
    fi
  fi

  # Attempt to download with curl, if found. If successful, quick return
  if command -v curl > /dev/null; then
    info "Downloading via curl: ${_url}"

    if [ -n "${SSL_CERT_FILE:-}" ]; then
      curl ${_curl_extra_args:+"--cacert ${SSL_CERT_FILE}"} -sSfL "${_url}" -o "${_dst}"
    else
      curl -sSfL "${_url}" -o "${_dst}"
    fi

    _code="$?"

    if [ $_code -eq 0 ]; then
      return 0
    else
      local _e="curl failed to download file, perhaps curl doesn't have"
      _e="$_e SSL support and/or no CA certificates are present?"
      warn "$_e"
    fi
  fi

  # If we reach this point, wget and curl have failed and we're out of options
  exit_with "Required: SSL-enabled 'curl' or 'wget' on PATH with" 6
}

main "$@" || exit 99
