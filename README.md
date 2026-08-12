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

## Run the current source preview

Sandwich Hime does not have immutable public release tags yet. Do not invent a
version-shaped install command: clone the compiler and this starter side by
side, then use a local Go workspace as an explicit preview bridge.

```sh
mkdir sandwich-hime-walk
cd sandwich-hime-walk

git clone https://gitea.speelman.ca/gamertan/sandwich-hime.git
git clone https://gitea.speelman.ca/gamertan/sandwich-hime-tutorial.git

cd sandwich-hime
go install ./cmd/himesan

cd ../sandwich-hime-tutorial
go work init .
go work edit -replace=gamertan.com/sandwich-hime/sando=../sandwich-hime/sando

./scripts/verify.sh
go run ./cmd/site
```

Make sure `$(go env GOPATH)/bin` is on `PATH`, or set `HIMESAN_BIN` to the
compiler executable when running the verification script. `go.work` and
`go.work.sum` are intentionally ignored: they are local preview wiring, not a
claim that `v0.0.0` was published.

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
```

The application owns the server, routing, headers, data, and deployment.
Sandwich Hime owns template compilation; `sando` owns the tiny runtime render
contract. Read the [language and security documentation](https://sandwichhime.com/docs/)
before accepting real user content.

## What the verification gate proves

`./scripts/verify.sh` checks committed output, generates twice and compares
digests, runs all tests and `go vet`, builds the server into a temporary
directory, and inspects its Go dependency graph. The only production
Sandwich Hime package allowed by that graph is
`gamertan.com/sandwich-hime/sando`.

The human-authored starter is [0BSD](LICENSES.md), specifically so copying it
does not drag a complicated license conversation into your application.
