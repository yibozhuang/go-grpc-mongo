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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new blog post",
	Long: `Create a new blog post over gRPC
A blog post requires an AuthorId, Title, and Content.`,
	Run: func(cmd *cobra.Command, args []string) {
		author, err := cmd.Flags().GetString("author")
		title, err := cmd.Flags().GetString("title")
		content, err := cmd.Flags().GetString("content")

		if err != nil {
			fmt.Sprintf("Error - %v", err)
		}

		blog := &blogpb.Blog{
			AuthorId: author,
			Title:    title,
			Content:  content,
		}

		res, err := client.CreateBlog(
			context.TODO(),
			&blogpb.CreateBlogReq{
				Blog: blog,
			},
		)

		if err != nil {
			fmt.Sprintf("Error - %v", err)
		}

		fmt.Printf("Blog created: %s", res.Blog.Id)
        fmt.Println()
	},
}

func init() {
	createCmd.Flags().StringP("author", "a", "", "An author ID")
	createCmd.Flags().StringP("title", "t", "", "A title for the blog")
	createCmd.Flags().StringP("content", "c", "", "The content of the blog")

	createCmd.MarkFlagRequired("author")
	createCmd.MarkFlagRequired("title")
	createCmd.MarkFlagRequired("content")

	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
