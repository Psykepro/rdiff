package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Psykepro/rdiff/pkg/differ"
	"github.com/Psykepro/rdiff/pkg/fileio"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: rdiff <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  signature -file <path_to_file> -chunk-size <chunk_size> -output <output_file>")
		fmt.Println("  delta -signature <signature_file> -updated <updated_file> -output <output_file>")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "signature":
		signatureCmd := flag.NewFlagSet("signature", flag.ExitOnError)
		file := signatureCmd.String("file", "", "Path to the file for which signatures will be generated")
		chunkSize := signatureCmd.Int("chunk-size", 16, "Size of each chunk in bytes")
		output := signatureCmd.String("output", "", "Path to the output file where the signatures will be stored")
		signatureCmd.Parse(os.Args[2:])

		if *file == "" || *output == "" {
			signatureCmd.Usage()
			os.Exit(1)
		}

		generateSignatures(*file, *chunkSize, *output)
	case "delta":
		deltaCmd := flag.NewFlagSet("delta", flag.ExitOnError)
		signatureFile := deltaCmd.String("signature", "", "Path to the file containing the signatures of the original file")
		updatedFile := deltaCmd.String("updated", "", "Path to the updated version of the file")
		chunkSize := deltaCmd.Int("chunk-size", 16, "Size of each chunk in bytes")
		output := deltaCmd.String("output", "", "Path to the output file where the delta will be stored")
		deltaCmd.Parse(os.Args[2:])

		if *signatureFile == "" || *updatedFile == "" || *output == "" || *chunkSize < 1 {
			deltaCmd.Usage()
			os.Exit(1)
		}

		fileHandler := fileio.NewFileHandler(*chunkSize) // Use the chunkSize from the command line arguments
		generateDelta(*signatureFile, *updatedFile, *output, fileHandler)
	case "print":
		printCmd := flag.NewFlagSet("print", flag.ExitOnError)
		deltaFile := printCmd.String("delta", "", "Path to the file containing the delta")
		printCmd.Parse(os.Args[2:])

		if *deltaFile == "" {
			printCmd.Usage()
			os.Exit(1)
		}

		printDelta(*deltaFile)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func printDelta(deltaFile string) {
	fileHandler := fileio.NewFileHandler(0) // Chunk size is not used for reading delta
	delta, err := fileHandler.ReadDelta(deltaFile)
	if err != nil {
		log.Fatal(err)
	}

	prettyDelta := differ.PrettifyDelta(delta)
	fmt.Printf("Pretty Delta:\n%+v\n", prettyDelta)
}

func generateSignatures(file string, chunkSize int, output string) {
	fileHandler := fileio.NewFileHandler(chunkSize)
	reader, err := fileHandler.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	differ := differ.New(chunkSize)
	signatures := differ.GenerateSignatures(reader)

	err = fileHandler.WriteSignatures(signatures, output)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Signatures generated and saved to: %s\n", output)
}

func generateDelta(signatureFile, updatedFile, output string, fileHandler fileio.FileHandler) {
	signatures, err := fileHandler.ReadSignatures(signatureFile)
	if err != nil {
		log.Fatal(err)
	}

	reader, err := fileHandler.Open(updatedFile)
	if err != nil {
		log.Fatal(err)
	}

	differ := differ.New(fileHandler.ChunkSize())
	delta := differ.GenerateDelta(signatures, reader)

	err = fileHandler.WriteDelta(delta, output)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Delta generated and saved to: %s\n", output)
}
