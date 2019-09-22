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

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a blog post",
	Long: `Update a blog post
ID here is the mongo unique identifier
All other parameters need to be provided as this is an overwrite rather than only update changed.`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("id")
		author, err := cmd.Flags().GetString("author")
		title, err := cmd.Flags().GetString("title")
		content, err := cmd.Flags().GetString("content")

        blog := &blogpb.Blog{
			Id:       id,
            AuthorId: author,
            Title:    title,
            Content:  content,
        }

		req := &blogpb.UpdateBlogReq{
		    Blog: blog,
        }

		res, err := client.UpdateBlog(context.Background(), req)

		if err != nil {
			fmt.Sprintf("Error - %v", err)
		}

		fmt.Println(res.GetBlog())
	},
}

func init() {
	updateCmd.Flags().StringP("id", "i", "", "The ID of the blog")
	updateCmd.Flags().StringP("author", "a", "", "Update the author")
	updateCmd.Flags().StringP("title", "t", "", "Update the title of the blog")
	updateCmd.Flags().StringP("content", "c", "", "Update the content of the blog")
	updateCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
