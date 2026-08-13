<!-- SPDX-License-Identifier: 0BSD -->

# Sandwich Hime tutorial starter

This is the runnable companion to the official
[Walk the path tutorial](https://sandwichhime.com/docs/tutorial/). The website
owns the lesson; this repository is the small application you can clone, run,
take apart, and turn into something of your own. It takes the lesson's minimal
`Link` model one small step further with a described `Trail` and adds visible
per-request proof, while preserving the same component structure.

**Use this, change it, or don’t—we’re just glad you’re here with us.**

The starter demonstrates the boundary plainly:

- Hime-san compiles three typed `.sando` templates into committed Go during
  development.
- `Badge`, `Home`, and `Layout` nest through the same `sando.Component`
  contract.
- An ordinary `net/http` handler builds the trail list and typed view data on
  every request.
- `?name=` is untrusted input and is contextually escaped by generated code.
- Each successful response contains a fresh UTC timestamp and process-local
  request number, uses `Cache-Control: no-store`, and reports buffered render
  time through `Server-Timing`.
- Production imports the Apache-2.0 `sando` runtime, not the compiler.

## Walk the path with the Beta 1 runtime and Beta 2 compiler

Install Go 1.25 or newer, then resolve the tiny runtime before installing the
immutable classroom compiler:

```sh
git clone https://gitea.speelman.ca/gamertan/sandwich-hime-tutorial.git
cd sandwich-hime-tutorial
GOWORK=off go mod download gamertan.com/sandwich-hime/sando@v1.0.0-beta.1
go install gamertan.com/sandwich-hime/cmd/himesan@v1.0.0-beta.2
./scripts/verify.sh
GOWORK=off go run ./cmd/site
```

On Windows PowerShell, use the native verifier:

```powershell
git clone https://gitea.speelman.ca/gamertan/sandwich-hime-tutorial.git
Set-Location sandwich-hime-tutorial
$env:GOWORK = "off"
go mod download gamertan.com/sandwich-hime/sando@v1.0.0-beta.1
go install gamertan.com/sandwich-hime/cmd/himesan@v1.0.0-beta.2
.\scripts\verify.ps1
go run ./cmd/site
```

Make sure Go's install directory—normally `$(go env GOPATH)/bin`—is on
`PATH`, or set `HIMESAN_BIN` to the full compiler path before running either
verifier. The starter deliberately runs with `GOWORK=off`: it proves the
application resolves the published Apache-2.0 runtime rather than a neighboring
development checkout. The verification scripts require compiler
`v1.0.0-beta.2` and runtime `v1.0.0-beta.1` exactly.

If an earlier compiler-first attempt reports that the parent module does not
contain `gamertan.com/sandwich-hime/sando`, repair only that module selection:

```sh
GOWORK=off go mod download gamertan.com/sandwich-hime/sando@v1.0.0-beta.1
GOWORK=off go get gamertan.com/sandwich-hime/sando@v1.0.0-beta.1
```

Then run the verifier again. You do not need to clear your whole Go module
cache.

The Beta 2 compiler and Beta 1 runtime are intended for classrooms, learning,
prototypes, and evaluation. Their
interfaces may still change before final v1. Linux and Windows have been
maintainer-tested; macOS is provisional while native maintainer testing is
pending. Useful Mac compatibility reports are welcome on canonical Gitea.

Open [http://127.0.0.1:8080/?name=Hime-san](http://127.0.0.1:8080/?name=Hime-san),
refresh it, and watch the request number and UTC time change. Then try:

```text
http://127.0.0.1:8080/?name=<script>alert("no")</script>
```

The browser displays those characters as text. They do not become markup or
script. The tests also prove that an ordinary `javascript:` URL is rejected at
render time.

## Project map

```text
cmd/site/main.go              application-owned listener
internal/server/handler.go    router, typed request data, and HTTP policy
internal/views/views.go       typed template contracts
internal/views/*.sando        templates people edit
internal/views/*.sando.go     committed generated Go
scripts/verify.sh             generation, tests, build, and dependency gate
scripts/verify.ps1            the same gate for native Windows PowerShell
```

The application owns the server, routing, headers, data, and deployment.
Sandwich Hime owns template compilation; `sando` owns the tiny runtime render
contract. Read the [language and security documentation](https://sandwichhime.com/docs/)
before accepting real user content.

## What the verification gate proves

`./scripts/verify.sh` and `scripts/verify.ps1` check committed output, generate
twice and compare digests, run all tests and `go vet`, build the server into a
temporary directory, and inspect its Go dependency graph. The only production
Sandwich Hime package allowed by that graph is
`gamertan.com/sandwich-hime/sando`.

The human-authored starter is [0BSD](LICENSES.md), specifically so copying it
does not drag a complicated license conversation into your application.
