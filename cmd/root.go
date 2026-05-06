package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	showVersion bool
	Version     = "dev"
	Commit      = "none"
	Date        = "today"
	BuiltBy     = "go build"
)

var rootCmd = &cobra.Command{
	Use:     "passgen",
	Short:   "Генератор паролей с настраиваемыми классами символов",
	Version: Version,
}

func Execute() {
	if showVersion {
		fmt.Println("passgen version:", Version)
		fmt.Println("from commit:", Commit)
		fmt.Println("built date:", Date)
		fmt.Println("built by:", BuiltBy)
		os.Exit(0)
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.passgen.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false, "Show version")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigName(".passgen")
		viper.SetConfigType("yaml")
	}
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
}
