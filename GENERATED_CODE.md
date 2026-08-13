<!-- SPDX-License-Identifier: 0BSD -->

# Generated code

Every `*.sando.go` file is an owned output of the neighboring `*.sando` source.
The committed neighbors were produced by `himesan v1.0.0-beta.1` and remain
byte-current under `himesan v1.0.0-beta.2`; Beta 2 intentionally preserves an
honest producer marker when generated semantics are unchanged. Commit both
files so production builds need only ordinary Go and the small `sando` runtime.

Never hand-edit a generated neighbor. Run `himesan generate internal/views`,
review the deterministic diff, and use `himesan check internal/views` in local
verification or CI. The generator records its version, runtime ABI, source
digest, and source mappings in each output.
