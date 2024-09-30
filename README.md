# DockClean

DockClean is a simple Go-based tool that cleans up unused Docker images from your system.

## Features

- List unused Docker images (images without tags)
- Remove unused Docker images with a single command

## Prerequisites

- Install Go - https://go.dev/dl/

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/mmuazam98/dockclean.git
   cd dockclean
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the binary:
   ```bash
   go build -o dockclean ./cmd/cleaner
   # On Windows Run
   go build -o dockclean.exe ./cmd/cleaner
   ```

## Usage

Once the binary is built, you can use it to clean up your Docker images.

## Basic Command

To list and remove unused Docker images:

```bash
./dockclean
# On Windows Run
dockclean.exe
```

This command will:

- Connect to the Docker daemon.
- List all images that are unused (untagged).
- Remove those images.

### Dry Run Mode

If you want to preview the images that would be deleted without actually removing them, you can use the `--dry-run` flag:

````bash
./dockclean --dry-run
# On Windows Run
dockclean.exe --dry-run

## Sample Output

```bash
Found 3 unused images
Successfully removed image sha256:abc123...
Successfully removed image sha256:def456...
Successfully removed image sha256:ghi789...
````

## Extending the Tool

You can easily extend the functionality by adding features such as:

- Filtering: Clean images based on size, age, or other criteria.
- Interactive mode: Prompt for confirmation before deletion.

## Contributing

Feel free to open an issue or submit a pull request if you:

- Encounter any bugs.
- Have ideas for new features.
- Want to improve the code or documentation.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
