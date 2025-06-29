# Git Remote Helper for OCI Registries Installation Guide

This Installation Guide provides the steps necessary to set up Git Remote Helper for OCI Registries so that it is ready to use with an existing ACT3 project or to start a new ACT3 project.

1. [Install Git Remote Helper for OCI Registries](#1-install-git-remote-helper-for-oci-registries)
2. [Create Configuration File](#2-create-configuration-file)

## 1. Install Git Remote Helper for OCI Registries

Git Remote Helper for OCI Registries is available for the following operating systems:

- [MacOS](#macos)
- [Linux](#linux)

> Git Remote Helper for OCI Registries is not supported natively on Windows, but it is installable via the Windows Subsystem for Linux (WSL).

### macOS

MacOS users should complete the [manual installation process](manual-installation-process.md).

### Linux

Linux users should complete the [manual installation process](manual-installation-process.md).

<!-- If your tool is part of dev-tools, uncomment the section below -->
<!-- Linux users who have already installed [Dev Tools](https://www.git.act3-ace.com/getting-started-at-act3/#ubuntu), Git Remote Helper for OCI Registries was already installed as part of the package.

> Linux user can still use the [manual installation process](manual-installation-process.md) at any time. -->

## 2. Create Configuration File

After Git Remote Helper for OCI Registries is installed, you are ready to configure it by following these steps:

- [Create config file](#create-config-file)
- [Open config file](#open-config-file)
- [Edit config file](#edit-config-file)
- [Test the config](#test-the-config)

### Create config file

```bash
gitoci config --write
```

> Read more about the [config command](cli/gitoci/config/index.md)

### Open config file

Open the config file in your default editor by running the config command with the `--open` flag:

```bash
gitoci config --open
```

> [I can't find my config file](#troubleshooting)

### Edit config file

<!-- Walk user through setting necessary configuration values -->

```yml
kind: Configuration
apiVersion: gitoci.act3-ai.io/v1alpha1

name: John Doe
```

> See more [configuration options](apis/out.md#configuration).

### Test the config

<!-- Provide a way for the user to make sure their config is working as expected -->

Otherwise, consult the [troubleshooting](#troubleshooting) section below or the [User Guide](user-guide.md).

## Installation Complete

The installation process is complete when you have successfully tested the config file and confirmed that Git Remote Helper for OCI Registries is correctly installed and working as intended.

If you want to start using Git Remote Helper for OCI Registries right away, move on to the [Quick Start Guide](quick-start-guide.md).

## Troubleshooting

### I can't find my config file

>>>
Run the command:

```bash
gitoci config --location
```
>>>

See the [User Guide](user-guide.md) for information about the default location and default values of the config file, including how to change the location where the config file is stored.

### I would like to put my config file somewhere different from the `XDG_CONFIG_DIR` variable

>>>
Specify a directory to look for your config file by setting the environment variable `GITOCI_CONFIG`.

The `gitoci config --write` command accepts an argument to specify a directory to write the file in: `gitoci config --write <path-to-config-file>`
>>>

If troubleshooting is unsuccessful, see the [additional resources](#additional-resources) section below.

## Additional Resources

- [Documentation](/README.md#documentation)
- [Support](/README.md#support)
