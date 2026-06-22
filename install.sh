#!/usr/bin/env sh
set -eu

repo="deepcoinapi/agent-cli"
bin_dir="${BIN_DIR:-$HOME/.local/bin}"
version="${DCLI_VERSION:-latest}"

usage() {
  cat <<'EOF'
Install dcli.

Usage:
  install.sh [-b <bin-dir>] [-v <version>]

Options:
  -b <bin-dir>   Install directory (default: ~/.local/bin)
  -v <version>   Release version/tag (default: latest)
  -h             Show help

Environment:
  BIN_DIR        Install directory
  DCLI_VERSION   Release version/tag
EOF
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    -b|--bin-dir)
      bin_dir="$2"
      shift 2
      ;;
    -v|--version)
      version="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown option: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "required command not found: $1" >&2
    exit 1
  fi
}

need uname
need mktemp
need tar

if command -v curl >/dev/null 2>&1; then
  fetch() {
    curl -fsSL "$1" -o "$2"
  }
elif command -v wget >/dev/null 2>&1; then
  fetch() {
    wget -qO "$2" "$1"
  }
else
  echo "required command not found: curl or wget" >&2
  exit 1
fi

os="$(uname -s)"
arch="$(uname -m)"

case "$os" in
  Darwin) os="Darwin" ;;
  Linux) os="Linux" ;;
  *)
    echo "unsupported OS: $os" >&2
    exit 1
    ;;
esac

case "$arch" in
  x86_64|amd64) arch="x86_64" ;;
  arm64|aarch64) arch="arm64" ;;
  *)
    echo "unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

if [ "$version" = "latest" ]; then
  url="https://github.com/$repo/releases/latest/download/dcli_${os}_${arch}.tar.gz"
else
  url="https://github.com/$repo/releases/download/$version/dcli_${os}_${arch}.tar.gz"
fi

tmp_dir="$(mktemp -d)"
cleanup() {
  rm -rf "$tmp_dir"
}
trap cleanup EXIT INT TERM

archive="$tmp_dir/dcli.tar.gz"
echo "Downloading $url"
fetch "$url" "$archive"

tar -xzf "$archive" -C "$tmp_dir"
mkdir -p "$bin_dir"
install_path="$bin_dir/dcli"

if command -v install >/dev/null 2>&1; then
  install -m 0755 "$tmp_dir/dcli" "$install_path"
else
  cp "$tmp_dir/dcli" "$install_path"
  chmod 0755 "$install_path"
fi

echo "Installed dcli to $install_path"
if ! command -v dcli >/dev/null 2>&1; then
  case ":$PATH:" in
    *":$bin_dir:"*) ;;
    *)
      echo "Add this to your shell profile:"
      echo "  export PATH=\"$bin_dir:\$PATH\""
      ;;
  esac
fi

"$install_path" --version >/dev/null
echo "dcli is ready."
