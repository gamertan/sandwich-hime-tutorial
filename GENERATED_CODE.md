<!-- SPDX-License-Identifier: 0BSD -->

# Generated code

Every `*.sando.go` file is an owned output of the neighboring `*.sando` source.
Commit both files so production builds need only ordinary Go and the small
`sando` runtime.

Never hand-edit a generated neighbor. Run `himesan generate internal/views`,
review the deterministic diff, and use `himesan check internal/views` in local
verification or CI. The generator records its version, runtime ABI, source
digest, and source mappings in each output.
