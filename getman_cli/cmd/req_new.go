/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/KonnorFrik/getman"
	"github.com/spf13/cobra"
)

// reqNewCmd represents the new command
var reqNewCmd = &cobra.Command{
	Use:   "new <collection> ...<collection>",
	Short: "Create a new request for at least one specified collection",
	Long: `For add variable - use pattern: "{{url}}/api/path"`,
	Args: cobra.MinimumNArgs(1),
	Run: _ReqNewCmd,
}

var (
	flagReqBuilderMethod string
	flagReqBuilderUrl string
	// TODO: need syntax for many header's "key=value" in one string, because can't use one flag many times :c
	flagReqBuilderHeader string
	flagReqBuilderBodyString string
	flagReqBuilderBodyFile string
	flagReqBuilderBodyBinary string
	// TODO: check is this realy need and how it work
	flagReqBuilderBodyRaw string

	flagReqBuilderAuthBasic string
	flagReqBuilderAuthBearer string
	flagReqBuilderAuthApiKey string

	flagReqBuilderCookieJar bool
	flagReqBuilderCookieDisable bool
)

func _ReqNewCmd(cmd *cobra.Command, args []string) {
	// TODO: validate neccesary flags for first (if empty - exit)
	client, err := createClientWithDirectory(cmd)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	reqBuilder := getman.NewRequestBuilder()
	// build the request 
	// load collections
	// append request into each loaded collection
	// save collection

}

func init() {
	reqCmd.AddCommand(reqNewCmd)
	reqNewCmd.Flags().StringVarP(&flagReqBuilderMethod, "method", "X", "", "HTTP method for request.")
	reqNewCmd.Flags().StringVar(&flagReqBuilderUrl, "url", "", "Url for request.")
	// reqNewCmd.Flags().StringVarP(&flagReqBuilderHeader, "header", "H", "", "HTTP header key-value for request.")
	reqNewCmd.Flags().StringVar(&flagReqBuilderBodyString, "data", "", "Text body for request.")
	reqNewCmd.Flags().StringVar(&flagReqBuilderBodyFile, "data-file", "", "Body for request from file as text.")
	reqNewCmd.Flags().StringVar(&flagReqBuilderBodyBinary, "data-binary", "", "Body for request from file as binary data.")
	reqNewCmd.Flags().StringVar(&flagReqBuilderBodyRaw, "data-raw", "", "Body for request from file as raw data.")

	reqNewCmd.Flags().StringVarP(&flagReqBuilderAuthBasic, "user", "U", "", "Basic Auth for request. Syntax: \"user:password\"")
	reqNewCmd.Flags().StringVarP(&flagReqBuilderAuthBearer, "bearer", "B", "", "Bearer Auth for request.")
	reqNewCmd.Flags().StringVarP(&flagReqBuilderAuthApiKey, "api-key", "K", "", "API Key Auth for request. Syntax: \"X-API-Key:value\"")

	// TODO: not implemented.
	// reqNewCmd.Flags().BoolVarP(&flagReqBuilderCookieJar, "cookie-jar", "c", false, "Enable auto manage for cookies.")
	// reqNewCmd.Flags().BoolVarP(&flagReqBuilderCookieDisable, "no-cookie", "C", false, "Disable cookies.")
}
