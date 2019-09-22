/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	blogpb "github.com/yibozhuang/go-grpc-mongo/proto"
	"google.golang.org/grpc"
)

var cfgFile string

// Client and context global variables
var client blogpb.BlogServiceClient
var requestCtx context.Context
var requestOpts grpc.DialOption

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "blog",
	Short: "a gRPC client to communicate with the BlogService server",
	Long: `client application to communicate with the BlogService server over gRPC.
Uses mongo as backend database storage.
You can use this client to create and read blogs.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	server_url := "localhost:50051"

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blog.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// intialize our client
	fmt.Println("Starting Blog Service client...")

	requestCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	// no TLS for now
	requestOpts = grpc.WithInsecure()

	conn, err := grpc.Dial(server_url, requestOpts)
	if err != nil {
		log.Fatalf("Unable to establish client connection to %s - %v", server_url, err)
	}

	client = blogpb.NewBlogServiceClient(conn)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".client" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".blog")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
