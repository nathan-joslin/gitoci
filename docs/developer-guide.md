# Git Remote Helper for OCI Registries Developer Guide

## Design Patterns

Git Remote Helper for OCI Registries is organized into three layers:

- [`cmd` Package](#cmd-package): CLI commands defined using the `cobra` package
- [`actions` Package](#actions-package): Main functionality of Git Remote Helper for OCI Registries
- [Other Packages](#other-packages): Purpose-separated components of Git Remote Helper for OCI Registries functionality

### `cmd` Package

The `cmd` package uses [`cobra`](https://pkg.go.dev/github.com/spf13/cobra) to define the command line interface for Git Remote Helper for OCI Registries.

> [`cmd` Package](./../cmd/gitoci/cmd)

### `actions` Package

The `actions` package contains the core functionality of Git Remote Helper for OCI Registries. The commands defined in `cmd` run and "action" in the `actions` package.

> [`actions` Package](./../pkg/actions)

### Other Packages

The other packages in the `pkg` folder contain smaller components of the functionality of Git Remote Helper for OCI Registries.

> [Other Packages](./../pkg)

## Testing

### Unit Tests

Run the following command from the root directory of the repository. This will run all unit tests

```bash
go test ./...

# or

make test
```

### Functional Tests

<!-- Describe how to run functional tests -->

## Releasing

The act3-pt CLI contains a `act3-pt ci release` command that automates this process.

## Code Generation

### Generate CLI Documentation (automatically done in CI/CD pipeline)

```bash
gitoci util gendocs <data output location>
```

Generate markdown documents from command usage descriptions and placed in the specified directory.

### Generate API Documentation (automatically done in CI/CD pipeline)

```bash
make apidoc
```

Generate API documentation using `crd-ref-docs`

### Generate API Deep Copy Functions

```bash
make generate
```

Generate deep copy functions for APIs using `controller-gen`
