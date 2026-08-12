# SPDX-License-Identifier: 0BSD

[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$repoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $repoRoot
$env:GOWORK = "off"

$himesan = if ($env:HIMESAN_BIN) { $env:HIMESAN_BIN } else { "himesan" }
$expectedVersion = "v1.0.0-beta.1"
if (-not (Get-Command $himesan -ErrorAction SilentlyContinue)) {
    throw "himesan was not found; install v1.0.0-beta.1 or set HIMESAN_BIN"
}

function Assert-LastExitCode([string]$Step) {
    if ($LASTEXITCODE -ne 0) {
        throw "$Step failed with exit code $LASTEXITCODE"
    }
}

function Get-GeneratedDigest {
    $lines = Get-ChildItem internal/views -Recurse -File -Filter *.sando.go |
        Sort-Object FullName |
        ForEach-Object {
            $hash = (Get-FileHash -Algorithm SHA256 -LiteralPath $_.FullName).Hash.ToLowerInvariant()
            "$($_.FullName):$hash"
        }
    return ($lines -join "`n")
}

$versionLine = (& $himesan version | Select-Object -First 1)
Assert-LastExitCode "himesan version"
if ($versionLine -notmatch '^himesan ([^ ]+) ') {
    throw "could not parse himesan version output: $versionLine"
}
if ($Matches[1] -ne $expectedVersion) {
    throw "himesan version is $($Matches[1]); expected $expectedVersion"
}

$runtimeVersion = (go list -m -f '{{.Version}}' gamertan.com/sandwich-hime/sando | Select-Object -First 1)
Assert-LastExitCode "sando runtime version"
if ($runtimeVersion -ne $expectedVersion) {
    throw "sando runtime version is $runtimeVersion; expected $expectedVersion"
}

& $himesan check internal/views
Assert-LastExitCode "himesan check"
$before = Get-GeneratedDigest
& $himesan generate internal/views
Assert-LastExitCode "first himesan generate"
$afterFirst = Get-GeneratedDigest
& $himesan generate internal/views
Assert-LastExitCode "second himesan generate"
$afterSecond = Get-GeneratedDigest
& $himesan check internal/views
Assert-LastExitCode "final himesan check"

if ($before -ne $afterFirst -or $afterFirst -ne $afterSecond) {
    throw "generated output was stale or nondeterministic"
}

go test ./...
Assert-LastExitCode "go test"
go vet ./...
Assert-LastExitCode "go vet"

$buildDir = Join-Path ([System.IO.Path]::GetTempPath()) ("sandwich-hime-tutorial-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $buildDir | Out-Null
try {
    go build -trimpath -o (Join-Path $buildDir "site.exe") ./cmd/site
    Assert-LastExitCode "go build"
} finally {
    Remove-Item -LiteralPath $buildDir -Recurse -Force -ErrorAction SilentlyContinue
}

$dependencies = @(go list -deps ./cmd/site)
Assert-LastExitCode "go list -deps"
if ($dependencies -notcontains "gamertan.com/sandwich-hime/sando") {
    throw "production dependency graph does not contain the sando runtime"
}
if ($dependencies | Where-Object { $_ -match '^gamertan\.com/sandwich-hime$|^gamertan\.com/sandwich-hime/(cmd|internal)(/|$)' }) {
    throw "production dependency graph contains the Sandwich Hime compiler"
}

Write-Output "verified Beta 1 identity, deterministic generation, tests, vet, build, and runtime-only production dependencies"
