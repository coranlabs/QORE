/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	//"fmt"

	"os"

	eupf "github.com/coranlabs/HEXA_UPF/cmd"
	"github.com/coranlabs/HEXA_UPF/internal/logger"
	"github.com/coranlabs/HEXA_UPF/pkg/factory"
	server "github.com/coranlabs/HEXA_UPF/src"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
)

//var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "HEXA_UPF",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitializeLogger(logrus.InfoLevel)
		cfg, err := factory.ReadConfig("./config/upfcfg.yaml")
		if err != nil {
			logger.MainLog.Tracef("there was error %v", err)
			server.Service()
		}

		if cfg.Mode == "eupf" {
			logger.MainLog.Tracef("1st if ebpf: %s", cfg.Mode)
			eupf.Emain()

		}
		if cfg.Mode == "ebpf" {
			logger.MainLog.Tracef("1st if ebpf: %s", cfg.Mode)
			server.Service()

		}
		if cfg.Mode == "normal" {
			logger.MainLog.Tracef("2nt if ebpf: %v", cfg.Mode)
			server.Action(cfg)
		} else {
			logger.MainLog.Tracef("else if ebpf: %v", cfg.Mode)
			server.Action(cfg)
		}
		// logger.MainLog.Tracef("rooooooooot")
		// server.Service()
		//Emain()
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

// func init() {
// 	cobra.OnInitialize(initConfig)

// 	// Here you will define your flags and configuration settings.
// 	// Cobra supports persistent flags, which, if defined here,
// 	// will be global for your application.

// 	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.HEXA_UPF.yaml)")

// 	// Cobra also supports local flags, which will only run
// 	// when this action is called directly.
// 	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }

// // initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := os.UserHomeDir()
// 		cobra.CheckErr(err)

// 		// Search config in home directory with name ".HEXA_UPF" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigType("yaml")
// 		viper.SetConfigName(".HEXA_UPF")
// 	}

// 	viper.AutomaticEnv() // read in environment variables that match

// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
// 	}
// }
