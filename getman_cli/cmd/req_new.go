/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

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
	flagReqBuilderHeader []string = []string{}
	flagReqBuilderBodyString string
	flagReqBuilderBodyFile string
	flagReqBuilderBodyBinary string
	// TODO: check is this realy need and how it work
	flagReqBuilderBodyRaw string

	flagReqBuilderAuthBasic string
	flagReqBuilderAuthBearer string
	flagReqBuilderAuthApiKey string

	// TODO: how to store cookies for requests? Files in storage ? External files ?
	flagReqBuilderCookieJar bool
	flagReqBuilderCookieDisable bool
)

func _ReqNewCmd(cmd *cobra.Command, args []string) {
	// TODO: validate neccesary flags for first (if empty - exit)
	if flagReqBuilderMethod == "" {
		PrintfError("flag '--method' is required\n")
		return
	}

	if flagReqBuilderUrl == "" {
		PrintfError("flag '--url' is required\n")
		return
	}

	client, err := createClientWithDirectory(cmd)

	if err != nil {
		PrintfError("%s\n", err)
		return
	}

	reqBuilder := getman.NewRequestBuilder()
	reqBuilder = reqBuilder.Method(strings.ToUpper(flagReqBuilderMethod))
	reqBuilder = reqBuilder.URL(flagReqBuilderUrl)

	for ind, hdr := range flagReqBuilderHeader {
		parts := strings.Split(hdr, ":")

		if len(parts) != 2 {
			PrintfError("Header #%d: '%s': invalid syntax\n", ind, hdr)
			continue
		}

		reqBuilder.Header(
			strings.TrimSpace(parts[0]),
			strings.TrimSpace(parts[1]),
		)
	}

	if flagReqBuilderBodyString != "" {
		reqBuilder.BodyString()
	}

	if flagReqBuilderBodyFile != "" {
		// open file and read as text
	}

	if flagReqBuilderBodyBinary != "" {
		// open file and read as bytes
	}

	// build the request 
	// load collections
	// append request into each loaded collection
	// save collection

}

func init() {
	reqCmd.AddCommand(reqNewCmd)
	reqNewCmd.Flags().StringVarP(&flagReqBuilderMethod, "method", "X", "", "HTTP method for request.")
	reqNewCmd.Flags().StringVar(&flagReqBuilderUrl, "url", "", "Url for request.")

	reqNewCmd.Flags().StringArrayVarP(&flagReqBuilderHeader, "header", "H", []string{}, "HTTP header key-value for request. Syntax: \"Header: Value\"")
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
