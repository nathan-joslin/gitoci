# Proposal: Storing Git Repositories as OCI Artifacts

**Author:** Nathan Joslin

**Advisors:** Dr. Kyle Tarplee

## Problem Description

Existing tools for transferring and synchronizing git repositories across air-gap boundaries lack an ability to scale, resulting in wasted time copying existing or nonessential data.

## Existing and Related Solutions

### Zarf

Zarf is a tool that facilitates software delivery on systems without an internet connection. They provide support for transferring git repositories across air-gap boundaries. However, using Zarf requires learning a large tool and they do not appear to support incremental updates. Additionally, they do not support git-lfs.

Sources:
- [Git-repositories docs](https://docs.zarf.dev/ref/components/#git-repositories)
- [Zarf codebase](https://github.com/zarf-dev/zarf/tree/main/src/internal/git)

### Git Bundles

Git bundles, native to git, are archives containing references and objects used for the offline transfer of git objects without an active network connection. With "shallow" bundle support, they can be used for incremental updates. For simple use cases git bundles are a great candidate, however their use at scale requires an in-depth knowledge of git and comes with significant organizational penalities for efficient transfers.

Sources:
- [Git bundle manpage](https://git-scm.com/docs/git-bundle)

## Proposed Solution

Define a specification for storing git repositories as OCI artifacts as well as implement a git remote helper to facilitate pushing and fetching from a git repository stored in an OCI registry.

Using git with an OCI remote helper looks as follows after installing the utility. While adding and cloning the repository requires specifying the remote helper, subsequent git commands may be used in their typical manner:

```bash
# Clone directly from OCI registry
$ git clone oci://example.registry.com/repository/name:tag

# Add an OCI registry as a remote
$ git remote add <remote-name> oci://example.registry.com/repository/name:tag
```

```bash
# Fetch from a remote OCI registry
$ git fetch [<repository>] [<refspec>...]

# Push to a remote OCI registry
$ git push [<repository>] [<refspec>...]
```

### Why OCI

The Open Container Initiative (OCI) has become the industry standard for the storage, transfer, and execution of container images. In recent years, the OCI specification has been used for alternative artifact types, extending beyond container images to: helm repositories, cryptographic signatures, software bill of materials, and more. By defining an OCI specification and implementing an interface for git to interact with OCI registries we can utilize existing data transfer tools.

### Why Git Remote Helper

Git remote helpers allow users to install programs that extend git's native capabilities to interface with alternative remote repository formats; in our case OCI. By using a git remote helper users can use git as they normally would, e.g. with VSCode, avoiding a need to learn new commands from a new tool. The [manual page](https://git-scm.com/docs/gitremote-helpers) includes well-defined instructions for implementing a git remote helper.

No known existing remote helpers for this application exist.

## Proof-of-Concept

A prototype specification and tooling is provided by the ASCE Data Tool, which has shown promising results. Although it is the closest existing solution, it only solves part of the problem. It is capable of storing git repositories as OCI artifacts, with support for git-lfs files. Like Zarf, it uses git references (tags and branches) or hashes to store a part of or an entire git repository in an OCI registry. However, it does not function as a git remote helper and lacks in some areas of speed and efficiency.

### Pros

- Facilitates efficient incremental updates by resolving the minimum data needed to perform updates based on the previous synchronization.
- Supports git-lfs.
- Prototype specification.

### Cons

- Slow to resolve the minimum data needed to perform an incremental update.
- As a subcommand of a larger tool, it requires learning new command groups. Not a git remote helper.

## Goals

- Git OCI Specification.
- Git remote helper implementation.
- Git-LFS support.
- Fast and efficient.
- 50% test coverage.

### Bonus Goals

The following goals are potential enhancements if time permits.

- Homebrew tap and formula, for easy installation.
- Git-lfs bundles.
- go-git bundle support (depending on spec).
- 80% test coverage.

## Implementation Details

The git remote helper will be written in Go for the following reasons:

- Familiarity and Experience with Go and external packages (see below)
- Well-written standard library
- Multi-Architecture support

### External Packages

- [ORAS](https://github.com/oras-project/oras-go): Efficient OCI artifact transfers, with private registry support.
- [sourcegraph/conc](https://github.com/sourcegraph/conc): Structured concurrency, for speed and efficiency.
- [go-git](https://github.com/go-git/go-git): Efficient handling of git objects.

### Git Remote Helper Capabilities

Git remote helpers are not required to support all capabilities and features. The following tables outline a plan for the *necessary* capabilities for an MVP. Additional capabilities labeled as "Maybe" or "Likely-No" may be added if time permits.

The *Support* column indicates if support for a capability is necessary to solve the problem, a nice-to-have, or no plans to support. Descriptions for each capability may be found in the linked manpage references.

#### Pushing

See the [Capabilities for Pushing](https://git-scm.com/docs/gitremote-helpers#_capabilities_for_pushing) section of the manpage.

| Capability     | Support     |
| ------------- | :-------------: |
| connect | Yes |
| stateless-connect | No |
| push | Yes |
| export | Likely-No |
| no-private-update | Likely-No |

#### Fetching

See the [Capabilities for Fetching](https://git-scm.com/docs/gitremote-helpers#_capabilities_for_fetching) section of the manpage.

| Capability     | Support     |
| ------------- | :-------------: |
| connect | Yes |
| stateless-connect | No |
| fetch | Yes |
| import | Likely-No |
| check-connectivity | Likely-No |
| get | Likely-No |

#### Miscellaneous

See the [Miscellaneous Capabilities](https://git-scm.com/docs/gitremote-helpers#_miscellaneous_capabilities) section of the manpage.

| Capability     | Support     |
| ------------- | :-------------: |
| option | Maybe |
| refspec | Maybe |
| bidi-import | Likely-No |
| export-marks | Likely-No |
| import-marks | Likely-No |
| signed-tags | Likely-No |
| object-format | Yes |

## Other Resources

- [Blog post](https://rovaughn.github.io/2015-2-9.html) on implementing a git remote helper