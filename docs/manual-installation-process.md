# Manual Installation of Git Remote Helper for OCI Registries

Manual installation of Git Remote Helper for OCI Registries requires two steps:

1. [Download the binary](#1-download-the-binary)
2. [Install the binary](#2-install-the-binary)

## 1. Download the binary

The git-remote-oci binary can be downloaded from the [Releases page](https://github.com/act3-ai/gitoci/-/releases):

- Download the binary corresponding to your operating system from the [Releases page](https://github.com/act3-ai/gitoci/-/releases)

   The options follow the naming scheme `git-remote-oci--<system>--<processor>` as illustrated below.

   | Operating System | Processor | Binary Name |
   | --- | --- | --- |
   | Linux | Intel/AMD | `git-remote-oci--linux--amd64` |
   | Linux (FIPS-compliant) | Intel/AMD | `git-remote-oci-fips--linux--amd64` |
   | macOS | Intel | `git-remote-oci--darwin-amd64` |
   | macOS | Apple | `git-remote-oci--darwin--arm64` |
   | Windows | Intel/AMD | `git-remote-oci--windows--amd64` |

- Rename the binary to `git-remote-oci`
- Move the binary to a directory named `bin` in your home directory (create a `~/bin` directory if one does not already exist)

## 2. Install the Binary

- Make the binary file executable

   ```bash
   sudo chmod +x ~/bin/git-remote-oci
   ```

- Add the `bin` directory to your PATH

   ```bash
   echo 'export PATH=$HOME/bin:$PATH' >> ~/.bashrc
   ```

- Reload your terminal

   ```bash
   source ~/.bashrc
   ```

- Additional macOS Installation Step

> Linux users can skip this step

Delete the "quarantine attribute" from the binary file. Run the following command in the same directory where the binary named git-remote-oci is located.

```bash
xattr -d com.apple.quarantine git-remote-oci
```

- Return to the [Installation Guide](installation-guide.md) and resume the configuration process after reloading the terminal.
