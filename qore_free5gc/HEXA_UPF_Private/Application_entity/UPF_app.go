/*
Copyright © 2025 NAME HERE satyam012005@gmail.com
*/
package upf_app

import (
	//"fmt"

	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/coranlabs/CORAN_UPF_eBPF/Messages_handling_entity/PFCP_server"
	UPF_config "github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/config"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/service"
	ebpf_datapath "github.com/coranlabs/CORAN_UPF_eBPF/eBPF_Datapath_entity"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
)

// var cfgFile string
func GracefulShutdown() context.Context {
	// Create a context with cancel to manage graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Create a channel to receive OS signals
	stopper := make(chan os.Signal, 1)

	// Register for specific signals
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)

	// Goroutine to wait for signals and cancel the context
	go func() {
		sig := <-stopper
		logger.InitLog.Infof("\nReceived signal: %v. Initiating shutdown...\n", sig)
		cancel()
	}()

	return ctx
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "CORAN_UPF_eBPF",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		UPF_config.Initialize()
		logger.Initialize(0)

		//logger.SetLogLevel(logrus.InfoLevel)
		//Handle interrupt

		ctx := GracefulShutdown()
		//all code logic goes here
		EBPF_controller,err := ebpf_datapath.Setup_eBPF(&UPF_config.Conf , ctx)
		if err != nil {
			logger.UPF_MAIN.Errorf("Failed to setup eBPF: %s", err)
		}
		defer EBPF_controller.Unload_and_detach(&EBPF_controller.Coran_ebpf_datapathObjects)
		ResourceManager, err := service.NewResourceManager(UPF_config.Conf.UEIPPool, UPF_config.Conf.FTEIDPool)
	if err != nil {
		logger.UPF_MAIN.Errorf("failed to create ResourceManager - err: %v", err)
	}

		PFCP_server ,err:= PFCP_server.Setup_PFCP_server(&UPF_config.Conf,  EBPF_controller, ResourceManager)
		if err != nil {
			logger.UPF_MAIN.Errorf("Failed to setup PFCP_server: %s", err)
		}
		defer PFCP_server.Close()
		
		logger.UPF_MAIN.Infof("Application started")
		<-ctx.Done()

		logger.UPF_MAIN.Infof("Application stopped")

		//server.Emain()

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

// 	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.CORAN_UPF_eBPF.yaml)")

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

// 		// Search config in home directory with name ".CORAN_UPF_eBPF" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigType("yaml")
// 		viper.SetConfigName(".CORAN_UPF_eBPF")
// 	}

// 	viper.AutomaticEnv() // read in environment variables that match

// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
// 	}
// }
