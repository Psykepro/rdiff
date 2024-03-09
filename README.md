# rdiff - Rolling Hash File Diffing Algorithm

This project implements a rolling hash based file diffing algorithm in Go. It compares an original file and an updated version of the file, and generates a delta describing the changes between the two versions. The delta can be used to upgrade the original file to the new version, reducing bandwidth and storage requirements in scenarios like distributed file storage systems.

## Features

- Generates signatures for the original file using a rolling hash algorithm (Adler-32)
- Computes deltas between the original and updated files, identifying changed, added, or deleted chunks
- Supports configurable chunk sizes for optimizing performance based on file size and network conditions
- Provides a Go package for integration into other projects

## Installation

1. Make sure you have Go installed on your system. You can download and install Go from the official website: [https://golang.org/](https://golang.org/)

2. Clone this repository to your local machine:
   ```bash
   git clone https://github.com/Psykepro/rdiff.git
   ```

3. Navigate to the project directory:
   ```bash
   cd rdiff
   ```

4. Build the project:
   ```bash
   go build ./cmd/..
   ```

## Usage

The project provides a command-line tool called `rdiff` for generating signatures and deltas between files.

### Generating Signatures

To generate signatures for a file, use the `rdiff signature` command:

```bash
./rdiff signature -file <path_to_file> -chunk-size <chunk_size> -output <output_file>
```

- `<path_to_file>`: Path to the file for which signatures will be generated.
- `<chunk_size>`: Size of each chunk in bytes (default: 16).
- `<output_file>`: Path to the output file where the signatures will be stored.

### Generating Delta

To generate a delta between the original file and an updated file, use the `rdiff delta` command:

```bash
./rdiff delta -signature <signature_file> -updated <updated_file> -output <output_file>
```

- `<signature_file>`: Path to the file containing the signatures of the original file.
- `<updated_file>`: Path to the updated version of the file.
- `<output_file>`: Path to the output file where the delta will be stored.

### Printing Delta

To print the delta in a human-readable format, use the `rdiff print` command:

```bash
./rdiff print -delta <delta_file>
```
- `<delta_file>`: Path to the file containing the delta.

## Testing

The project includes unit tests to ensure the correctness of the rolling hash algorithm and the diffing functionality. To run the tests, use the following command:

```bash
go test ./...
```