/*
Copyright © 2022 Aleksey Sviridkin <f@lex.la>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	DataType string
	cfgFile  string
	Parallel uint
)

// rootCmd represents the base command when called without any subcommands.
//nolint:exhaustivestruct // Not needed here
var rootCmd = &cobra.Command{
	Use: "mtgdsgenerator",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mtgdsgenerator.yaml)")
	rootCmd.PersistentFlags().StringVar(&DataType, "datatype", "all_cards",
		"type of cards archive (all_cards, oracle_cards, etc)")
	rootCmd.PersistentFlags().UintVar(&Parallel, "parallel", 10, "number of parallel downloaders")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mtgdsgenerator" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mtgdsgenerator")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
