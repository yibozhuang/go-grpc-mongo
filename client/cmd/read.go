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

	"github.com/spf13/cobra"
	blogpb "github.com/yibozhuang/go-grpc-mongo/proto"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Get the blog post by the ID",
	Long: `Find a blog post by ID
The ID is the mongo unique identifier.`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("id")

		if err != nil {
			fmt.Sprintf("Error - %v", err)
		}

		req := &blogpb.ReadBlogReq{
			Id: id,
		}

		res, err := client.ReadBlog(context.Background(), req)

		if err != nil {
			fmt.Sprintf("Error - %v", err)
		}

		fmt.Println(res.GetBlog())
	},
}

func init() {
	readCmd.Flags().StringP("id", "i", "", "The ID of the blog")
	readCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
