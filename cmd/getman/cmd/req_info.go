/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/KonnorFrik/getman/types"
	"github.com/spf13/cobra"
)

// reqInfoCmd represents the info command
var reqInfoCmd = &cobra.Command{
	Use:   "info <collection> <request name> ...<request name>",
	Short: "Show information about request.",
	Long: `For show info for all requests - don't specify any reqeust's name`,
	Args: cobra.MinimumNArgs(1),
	Run: _ReqInfoCmd,
}

func _ReqInfoCmd(cmd *cobra.Command, args []string) {
	client, err := createClientWithDirectory(cmd)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	collName := args[0]
	coll, err := client.LoadCollection(collName)

	if err != nil {
		PrintfError("Can't load collection: %s: %s\n", collName, err)
		return
	}

	reqNamesSet := convertSlcToSet(args[1:])

	for _, req := range coll.Items {
		if len(args) > 1 {
			if _, exist := reqNamesSet[req.Name]; exist {
				printRequestInfo(req)
			}

		} else {
			printRequestInfo(req)
		}

		println()
	}
}

func init() {
	reqCmd.AddCommand(reqInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func convertSlcToSet[T comparable](slc []T) map[T]struct{} {
	var result = make(map[T]struct{})

	if len(slc) == 0 {
		return result
	}

	for _, v := range slc {
		result[v] = struct{}{}
	}

	return result
}

func printRequestInfo(req *types.RequestItem) {
	fmt.Printf("%s:\n", req.Name)
	fmt.Printf("\t%s %s\n", req.Request.Method, req.Request.URL)
	fmt.Printf("\tHeaders:\n")

	if req.Request.Headers != nil {
		for k, v := range req.Request.Headers {
			fmt.Printf("\t%s:%s\n", k, v)
		}
	}

	if req.Request.Timeout != nil {
		fmt.Printf("Timeout connect: %s\n", req.Request.Timeout.Connect)
		fmt.Printf("Timeout read: %s\n", req.Request.Timeout.Read)
	}
	// TODO: add more info
}
