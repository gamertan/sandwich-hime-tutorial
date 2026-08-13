#!/usr/bin/env bash
# SPDX-License-Identifier: 0BSD

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_root"

himesan_bin=${HIMESAN_BIN:-himesan}
expected_compiler_version=v1.0.0-beta.2
expected_runtime_version=v1.0.0-beta.1
if ! command -v "$himesan_bin" >/dev/null 2>&1; then
	echo "himesan was not found; install v1.0.0-beta.2 or set HIMESAN_BIN" >&2
	exit 1
fi

export GOWORK=off

actual_version=$("$himesan_bin" version | awk 'NR == 1 { print $2 }')
if [[ "$actual_version" != "$expected_compiler_version" ]]; then
	echo "himesan version is $actual_version; expected $expected_compiler_version" >&2
	exit 1
fi

runtime_version=$(go list -m -f '{{.Version}}' gamertan.com/sandwich-hime/sando)
if [[ "$runtime_version" != "$expected_runtime_version" ]]; then
	echo "sando runtime version is $runtime_version; expected $expected_runtime_version" >&2
	exit 1
fi

generated_digest() {
	find internal/views -type f -name '*.sando.go' -print0 \
		| sort -z \
		| xargs -0 sha256sum
}

"$himesan_bin" check internal/views
before=$(generated_digest)
"$himesan_bin" generate internal/views
after_first=$(generated_digest)
"$himesan_bin" generate internal/views
after_second=$(generated_digest)
"$himesan_bin" check internal/views

if [[ "$before" != "$after_first" || "$after_first" != "$after_second" ]]; then
	echo "generated output was stale or nondeterministic" >&2
	diff -u <(printf '%s\n' "$before") <(printf '%s\n' "$after_second") || true
	exit 1
fi

go test ./...
go vet ./...

build_dir=$(mktemp -d)
trap 'rm -rf "$build_dir"' EXIT
go build -trimpath -o "$build_dir/site" ./cmd/site

dependencies=$(go list -deps ./cmd/site)
if ! grep -qx 'gamertan.com/sandwich-hime/sando' <<<"$dependencies"; then
	echo "production dependency graph does not contain the sando runtime" >&2
	exit 1
fi
if grep -Eq '^gamertan\.com/sandwich-hime$|^gamertan\.com/sandwich-hime/(cmd|internal)(/|$)' <<<"$dependencies"; then
	echo "production dependency graph contains the Sandwich Hime compiler" >&2
	exit 1
fi

echo "verified Beta 2 compiler, Beta 1 runtime, deterministic generation, tests, vet, build, and runtime-only production dependencies"
