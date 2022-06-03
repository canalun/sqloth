package cmd

import (
	"fmt"
	"os"

	"github.com/canalun/sqloth/driver/file_driver"
	"github.com/canalun/sqloth/usecase"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sqloth",
	Short: "generate dummy data for given schema",
	Long:  ``,
	// TODO: good variable name
	Run: func(cmd *cobra.Command, args []string) {
		fp, _ := cmd.Flags().GetString("filePath")
		num, _ := cmd.Flags().GetInt("recordNumber")

		fd := file_driver.NewFileDriver(fp)
		u := usecase.NewUsecase(fd)

		queries := u.GenerateQueryOfDummyData(num)
		for _, query := range queries {
			fmt.Printf("%s\n\n", query)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sqloth.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().IntP("recordNumber", "n", 10, "the # of records you want")
	rootCmd.Flags().StringP("filePath", "f", "./dump.sql", "the path to the schema sql file")
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

		// Search config in home directory with name ".sqloth" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".sqloth")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
