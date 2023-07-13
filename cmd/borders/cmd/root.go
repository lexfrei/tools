/*
Copyright © 2023 Aleksei Sviridkin <f@lex.la>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

 3. Neither the name of the copyright holder nor the names of its contributors
    may be used to endorse or promote products derived from this software
    without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "borders",
	Short: "Add borders to images to make them square",
	Long: `Add borders to images to make them square (or any other proportions)

Examples:
	to add 10% red border to all images in current directory:
		borders -c #FF0000 -a 10
	to add white border to image.jpg
		borders -f /path/to/image.jpg -c white

Possible colors:
	black, white, red, green, blue, yellow, cyan, magenta, gray,
	darkgray, lightgray, orange, pink, purple, violet, brown
Or any hex color like #ff0000
Or "avg" to get average color of the image`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Your application logic goes here
		return nil
	},
}

var (
	// BorderColor is the color of the border.
	BorderColor string
	// InputFile is the path to the input file.
	InputFile string
	// Directory is the path to the directory with images.
	Directory string
	// AdditionalBorder is the additional border size in percents.
	AdditionalBorder int
	// Prefix is the prefix for the output file.
	Prefix string
	// Minimal is the flag for minimal needed border size for instagram.
	Minimal bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	if rootCmd.Flags().Changed("help") {
		os.Exit(0)
	}
}

//nolint:lll // it's ok
func init() {
	rootCmd.PersistentFlags().StringVarP(&BorderColor, "color", "c", "white", "border color")
	rootCmd.PersistentFlags().StringVarP(&InputFile, "file", "f", "", "input file")
	rootCmd.PersistentFlags().StringVarP(&Directory, "directory", "d", ".", "directory with images")
	rootCmd.PersistentFlags().IntVarP(&AdditionalBorder, "additional-border", "a", 0, "additional border size in percents")
	rootCmd.PersistentFlags().StringVarP(&Prefix, "prefix", "p", "bordered", "prefix for the output file")
	rootCmd.PersistentFlags().BoolVarP(&Minimal, "minimal", "m", false, "minimal needed border size for instagram, non square images will be created")

	rootCmd.Flags().BoolP("help", "h", false, "Help message for toggle")
}
