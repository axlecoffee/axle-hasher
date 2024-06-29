package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/schollz/progressbar"
	"github.com/spf13/cobra"
)

func main() {
	green := color.New(color.FgHiGreen)
	red := color.New(color.FgHiRed)
	yellow := color.New(color.FgHiYellow)

	rootCmd := &cobra.Command{
		Use:   "ahash",
		Short: "Hash files in a directory",
		Run: func(cmd *cobra.Command, args []string) {
			inputDir, _ := cmd.Flags().GetString("input")
			verbose, _ := cmd.Flags().GetBool("verbose")
			clean, _ := cmd.Flags().GetBool("clean")
			outputDir, _ := cmd.Flags().GetString("output")
			ext, _ := cmd.Flags().GetString("ext")

			if inputDir == "" {
				inputDir, _ = os.Getwd()
			}
			if outputDir == "" {
				outputDir, _ = os.Getwd()
			}
			if ext == "" {
				ext = "*"
			}
			if !clean {
				clean = false
			}

			files, err := filepath.Glob(filepath.Join(inputDir, "*."+ext))
			if err != nil {
				log.Fatal(err)
			}
			// Welcome Message
			color.Set(color.FgCyan, color.Bold)
			fmt.Println("Welcome to Axle Hasher - A simple multi-file CLI tool")
			color.Unset()
			fmt.Print("INPUT: ")
			yellow.Print(inputDir)
			color.Unset()
			fmt.Print("\n")
			fmt.Print("OUTPUT: ")
			yellow.Print(outputDir)
			color.Unset()
			if clean {
				fmt.Print("\n")
				fmt.Print("Clean Output Logs: ")
				green.Print(strconv.FormatBool(clean))
				color.Unset()
			} else {
				fmt.Print("\n")
				fmt.Print("Clean Output Logs: ")
				red.Print(strconv.FormatBool(clean))
				color.Unset()
			}
			if verbose {
				fmt.Print("\n")
				fmt.Print("Using Verbose Console Logs: ")
				green.Print(strconv.FormatBool(verbose))
				color.Unset()
			} else {
				fmt.Print("\n")
				fmt.Print("Using Verbose Console Logs: ")
				red.Print(strconv.FormatBool(verbose))
				color.Unset()
			}
			fmt.Print("\n")
			fmt.Print("Filtering Extensions: ")
			yellow.Print(ext)
			color.Unset()
			fmt.Println()

			//green.Println("Welcome To Axle Hasher - A simple multi-file CLI tool\nIMPUT: " + inputDir + "\nOUTPUT: " + outputDir + "\nFilter Extensions: " + ext + "\nClean Output Logs: " + strconv.FormatBool(clean) + "\nUsing Verbose: " + strconv.FormatBool(verbose) + "\n")
			bar := progressbar.New(len(files))
			hashes := make(map[string]string)

			for _, file := range files {
				hash, err := hashFile(file)
				if err != nil {
					log.Fatal(err)
				}
				hashes[filepath.Base(file)] = hash
				bar.Add(1)

			}
			fmt.Printf("\n")
			now := time.Now()
			year := now.Year() % 100 // get the last 2 digits of the year
			month := int(now.Month())
			day := now.Day()
			dateString := fmt.Sprintf("%02d%02d%02d", year, month, day)

			outputFile := filepath.Join(outputDir, fmt.Sprintf("hashes-output-%s.txt", dateString))
			f, err := os.Create(outputFile)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			for file, hash := range hashes {
				if verbose {
					fmt.Printf("%s: %s\n\n", file, hash)
					if clean {
						fmt.Fprintf(f, "%s\n", hash)
					} else {
						fmt.Fprintf(f, "%s: %s\n", file, hash)
					}
				}
				if clean {
					fmt.Fprintf(f, "%s\n", hash)
				} else {
					fmt.Fprintf(f, "%s: %s\n", file, hash)
				}
			}

			fmt.Println("\nHashes saved to", outputFile)
		},
	}

	rootCmd.Flags().StringP("input", "i", "", "Input directory (default: current directory)")
	rootCmd.Flags().BoolP("clean", "c", false, "Clean output (no file names)")
	rootCmd.Flags().BoolP("verbose", "v", false, "Log each hash to console. (This will always use clean)")
	rootCmd.Flags().StringP("output", "o", "", "Output directory (default: current directory)")
	rootCmd.Flags().StringP("ext", "e", "", "File extension to hash (e.g. .jar, default: all files)")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func hashFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write(data)
	hash := hex.EncodeToString(h.Sum(nil))

	return hash, nil
}
