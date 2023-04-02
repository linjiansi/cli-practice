/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var inputFile string
var outputFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli-practice",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		convertVideoToGif(inputFile, outputFile)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli-practice.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input", "i", "", "Input video file")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Output GIF file")
}

func convertVideoToGif(input string, output string) {
	if input == "" || output == "" {
		fmt.Println("Please specify input and output file paths")
		return
	}

	ffmegCmd := exec.Command(
		"ffmpeg",
		"-i", input, output,
	)

	if checkFFmpegInstalled() {
		installFFmpeg()
	}

	if err := deleteFile(".", output); err != nil {
		fmt.Println(err)
	}

	err := ffmegCmd.Run()
	if err != nil {
		fmt.Printf("Error occured while converting video to GIF: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("GIF created successfully!")
}

func checkFFmpegInstalled() bool {
	_, err := exec.LookPath("ffmpeg")
	return err != nil
}

func installFFmpeg() {
	fmt.Println("Attempting to install FFmpeg...")
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("brew", "install", "ffmpeg")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to install FFmpeg using Homebrew: %v\n", err)
			fmt.Println("Please install FFmpeg manually and try again.")
			os.Exit(1)
		} else {
			fmt.Println("FFmpeg installed successfully!")
		}
	case "linux":
		cmd := exec.Command("sudo", "apt", "install", "ffmpeg")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to install FFmpeg using apt: %v\n", err)
			fmt.Println("Please install FFmpeg manually and try again.")
			os.Exit(1)
		} else {
			fmt.Println("FFmpeg installed successfully!")
		}
	default:
		fmt.Println("Unsupported platform. Please install FFmpeg manually and try again.")
		os.Exit(1)
	}
}

func deleteFile(directory string, filename string) error {
	filePath := directory + string(os.PathSeparator) + filename
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("File %s does not exist\n", filePath)
	}

	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("Failed to delete file %s: %v\n", filePath, err)
	}

	return nil
}
