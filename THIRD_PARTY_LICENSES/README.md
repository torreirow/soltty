# Third-party licenses

Soltty is distributed under the MIT License (see [`../LICENSE`](../LICENSE)).

Soltty is statically linked, so every released binary contains compiled code from
the Go modules below. Their license texts are reproduced in this directory and
ship inside every release archive, as their terms require.

| Module                                 | Version               | License      | Shipped in          |
|----------------------------------------|-----------------------|--------------|---------------------|
| `github.com/pkg/browser`               | v0.0.0-20240102092130 | BSD-2-Clause | all builds          |
| `github.com/spf13/cobra`               | v1.8.0                | Apache-2.0   | all builds          |
| `github.com/spf13/pflag`               | v1.0.5                | BSD-3-Clause | all builds          |
| `github.com/inconshreveable/mousetrap` | v1.1.0                | Apache-2.0   | Windows builds only |
| `golang.org/x/sys`                     | v0.1.0                | BSD-3-Clause | Windows builds only |

The two Windows-only modules are pulled in by `cobra` on that platform. They are
listed and shipped unconditionally so that one archive layout serves every target.

## Notes on compliance

- **Apache-2.0** (`cobra`, `mousetrap`): the license text is included as §4(a)
  requires. Neither module ships a `NOTICE` file, so §4(d) does not apply. Soltty
  does not modify either module, so there are no changes to state under §4(b).
- **BSD-2-Clause / BSD-3-Clause** (`browser`, `pflag`, `x/sys`): the copyright
  notice, condition list and disclaimer are reproduced in full, as redistribution
  in binary form requires.
- None of these licenses is copyleft; all permit redistribution within an
  MIT-licensed work.

## Not included, and why

`go-md2man`, `blackfriday`, `gopkg.in/yaml.v3` and `gopkg.in/check.v1` appear in
the module graph as documentation-generation and test dependencies of `cobra`.
They are not linked into any released binary — verified with `go list -deps` for
each target platform — so their licenses are not reproduced here.

## Regenerating

Run [`../scripts/update-third-party-licenses.sh`](../scripts/update-third-party-licenses.sh)
after changing dependencies. It re-reads the license texts from the Go module
cache and fails if a linked module is missing from this directory.
