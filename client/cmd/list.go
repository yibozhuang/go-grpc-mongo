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
	"io"

	"github.com/spf13/cobra"
	blogpb "github.com/yibozhuang/go-grpc-mongo/proto"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all blog posts",
	Long: `Lists all blog posts over gRPC.
Uses the streaming feature of gRPC.`,
	Run: func(cmd *cobra.Command, args []string) {
		req := &blogpb.ListBlogsReq{}

		stream, err := client.ListBlogs(context.Background(), req)

		if err != nil {
			fmt.Sprintf("Error - %v", err)
		}

		for {
			// stream.Recv returns a pointer to the Res object for current iteration
			res, err := stream.Recv()

			// If end of stream, break the loop
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Sprintf("Error - %v", err)
			}

			fmt.Println(res.GetBlog())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
