# Manual Installation of Git Remote Helper for OCI Registries

Manual installation of Git Remote Helper for OCI Registries requires two steps:

1. [Download the binary](#1-download-the-binary)
2. [Install the binary](#2-install-the-binary)

## 1. Download the binary

The gitoci binary can be downloaded from the [Releases page](https://github.com/act3-ai/gitoci/-/releases):

- Download the binary corresponding to your operating system from the [Releases page](https://github.com/act3-ai/gitoci/-/releases)

   The options follow the naming scheme `gitoci--<system>--<processor>` as illustrated below.

   | Operating System | Processor | Binary Name |
   | --- | --- | --- |
   | Linux | Intel/AMD | `gitoci--linux--amd64` |
   | Linux (FIPS-compliant) | Intel/AMD | `gitoci-fips--linux--amd64` |
   | macOS | Intel | `gitoci--darwin-amd64` |
   | macOS | Apple | `gitoci--darwin--arm64` |
   | Windows | Intel/AMD | `gitoci--windows--amd64` |

- Rename the binary to `gitoci`
- Move the binary to a directory named `bin` in your home directory (create a `~/bin` directory if one does not already exist)

## 2. Install the Binary

- Make the binary file executable

   ```bash
   sudo chmod +x ~/bin/gitoci
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

Delete the "quarantine attribute" from the binary file. Run the following command in the same directory where the binary named gitoci is located.

```bash
xattr -d com.apple.quarantine gitoci
```

- Return to the [Installation Guide](installation-guide.md) and resume the configuration process after reloading the terminal.
