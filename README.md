# DockClean

DockClean is a Go-based tool designed to clean up unused Docker images, enhancing storage management with options like size limits, verbose mode, and concurrent deletion.

## Features

- List and remove unused Docker images (images without tags)
- Clean images exceeding a specified size limit
- Option to remove images associated with stopped containers
- Verbose mode for detailed output
- Concurrent deletion for faster cleanup

## Prerequisites

- Install Go: https://go.dev/dl/
- **Optional**: For testing the application, you can generate a set of unused Docker images with predefined scripts. These scripts will create several untagged images, which can be cleaned by DockClean after each run.

  - Place a `Dockerfile` in the project directory.
  - Use the following scripts to create multiple unlabelled images:
    
    - **Windows (PowerShell)**:  
      ```powershell
      .\scripts\test\generate-unused-images.ps1
      ```

    - **Linux (Bash)**:  
      ```bash
      ./scripts/test/generate-unused-images.sh
      ```

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

Once the binary is built, you can use DockClean to clean up your Docker images.

### 1. Basic Command

To list and remove unused Docker images sequentially:

```bash
./dockclean
# On Windows Run
dockclean.exe
```

This command will:

- Connect to the Docker daemon.
- List all unused (untagged) images.
- Remove these images in a **sequential** order.

**Sample Output (Sequential Execution):**

- **Execution Time:** 611 ms  
  ![Sequential Execution](https://github.com/user-attachments/assets/68eadf90-6171-45cc-891d-d57302ebf7aa)

### 2. Enabling Concurrent Deletion

To remove images concurrently (using Go routines for faster cleanup), add the `--concurrent` flag. This flag enables DockClean to remove images in parallel, resulting in quicker execution times compared to sequential deletion.

**Example:**

```bash
./dockclean --concurrent
# On Windows Run
dockclean.exe --concurrent
```

**Sample Output (Concurrent Execution):**

- **Execution Time:** 143 ms  
  ![Concurrent Execution](https://github.com/user-attachments/assets/0f6bb59d-1d51-4d37-ada9-81dc8601eefb)

### Command-Line Options

#### 3. `--dry-run` : List unused images without deleting them.

**Example:**

```bash
./dockclean --dry-run
# On Windows Run
dockclean.exe --dry-run
```

This option will show which images would be deleted without performing any actual deletion.

#### 4. `--remove-stopped` : Remove images associated with stopped containers.

**Example:**

```bash
./dockclean --remove-stopped
# On Windows Run
dockclean.exe --remove-stopped
```

This command removes Docker images that are only associated with stopped containers, freeing up space from unused images.

**Sample Output:**

![Remove Stopped Images](https://github.com/user-attachments/assets/de9a38fa-625e-4053-b4af-6138661f6cc9)

#### 5. `--verbose` : Enable detailed output for each image during cleanup.

**Example:**

```bash
./dockclean --verbose
# On Windows Run
dockclean.exe --verbose
```

Verbose mode provides additional details about each image, such as its size and creation date, making it easier to see what’s being deleted.

**Sample Output:**

![Verbose Mode](https://github.com/user-attachments/assets/35a8c5fe-4851-4aae-ae92-75e354820b4e)

#### 6. `--size-limit <value>` : Set a size threshold for deleting images (e.g., `--size-limit 500MB`).
  - Specify units with `--B`, `--KB`, `--MB`, or `--GB`.

**Example:**

```bash
./dockclean --size-limit 500MB
# On Windows Run
dockclean.exe --size-limit 500MB
```

This command will only delete images that exceed the specified size limit.

**Sample Output:**

![Size Limit](https://github.com/user-attachments/assets/35a8c5fe-4851-4aae-ae92-75e354820b4e)


## Extending DockClean

You can extend DockClean’s functionality by adding features like:

- **Filtering**: Clean images based on additional criteria such as age.
- **Interactive mode**: Prompt for confirmation before deletion.

## Contributing

Feel free to open an issue or submit a pull request if you:

- Find bugs
- Have feature suggestions
- Want to improve code quality or documentation

## License

This project is licensed under the MIT License. See the LICENSE file for details. 
