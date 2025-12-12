/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	getman "github.com/KonnorFrik/getman/client"
	"github.com/KonnorFrik/getman/types"
	"github.com/spf13/cobra"
)

// reqNewCmd represents the new command
var reqNewCmd = &cobra.Command{
	Use:   "new <collection> ...<collection>",
	Short: "Create a new request for at least one specified collection",
	Long:  `For use variables from env - use pattern: "{{url}}/api/path"`,
	Args:  cobra.MinimumNArgs(1),
	Run:   _ReqNewCmd,
}

var (
	flagRequestName          string
	flagReqBuilderMethod     string
	flagReqBuilderUrl        string
	flagReqBuilderHeader     []string = []string{}
	flagReqBuilderBodyString string
	flagReqBuilderBodyFile   string
	flagReqBuilderBodyBinary string

	flagReqBuilderAuthBasic  string
	flagReqBuilderAuthBearer string
	// flagReqBuilderAuthApiKey string

	// TODO: how to store cookies for requests? Files in storage ? External files ?
	// flagReqBuilderCookieJar bool
	// flagReqBuilderCookieDisable bool
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

		reqBuilder = reqBuilder.Header(
			strings.TrimSpace(parts[0]),
			strings.TrimSpace(parts[1]),
		)
	}

	// TODO: add priority or exluding flags if one already setted
	if flagReqBuilderBodyString != "" {
		reqBuilder = reqBuilder.BodyString(flagReqBuilderBodyString)
	}

	if flagReqBuilderBodyFile != "" {
		content, err := readFileAsText(flagReqBuilderBodyFile)

		if err != nil {
			PrintfError("Can't read file: %s: %s\n", flagReqBuilderBodyFile, err)
			return
		}

		reqBuilder = reqBuilder.BodyString(content)
	}

	if flagReqBuilderBodyBinary != "" {
		content, err := readFileAsBytes(flagReqBuilderBodyBinary)

		if err != nil {
			PrintfError("Can't read file: %s: %s\n", flagReqBuilderBodyFile, err)
			return
		}

		reqBuilder = reqBuilder.BodyBinary(content, "")
	}

	// TODO: add priority or exluding flags if one already setted
	if flagReqBuilderAuthBasic != "" {
		parts := strings.Split(flagReqBuilderAuthBasic, ":")

		if len(parts) != 2 {
			PrintfError("BasicAuth '%s': invalid syntax\n", flagReqBuilderAuthBasic)
			return
		}

		reqBuilder = reqBuilder.AuthBasic(parts[0], parts[1])
	}

	if flagReqBuilderAuthBearer != "" {
		reqBuilder = reqBuilder.AuthBearer(flagReqBuilderAuthBearer)
	}

	// TODO: find a way for: how to set location for api key auth
	// if flagReqBuilderAuthBearer != "" {
	// 	reqBuilder.AuthAPIKey(flagReqBuilderAuthBearer)
	// }

	req, err := reqBuilder.Build()

	if err != nil {
		PrintfError("Can't build request: %s\n", err)
		return
	}

	for _, collName := range args {
		coll, err := client.LoadCollection(collName)

		if err != nil {
			PrintfError("Can't load collection: %s: %s\n", collName, err)

			if flagExitOnError {
				fmt.Printf("Exit on error.")
				return
			}

			continue
		}

		coll.Items = append(coll.Items, &types.RequestItem{
			Name:    flagRequestName,
			Request: req,
		})

		err = client.SaveCollection(coll)

		if err != nil {
			PrintfError("Can't save collection: %s: %s\n", collName, err)

			if flagExitOnError {
				fmt.Printf("Exit on error.")
				return
			}
		}
	}
}

func init() {
	reqCmd.AddCommand(reqNewCmd)

	reqNewCmd.Flags().StringVarP(&flagRequestName, "name", "n", "default-name", "Name for request")

	reqNewCmd.Flags().StringVarP(&flagReqBuilderMethod, "method", "X", "", "HTTP method for request.")
	reqNewCmd.Flags().StringVar(&flagReqBuilderUrl, "url", "", "Url for request.")

	reqNewCmd.Flags().StringArrayVarP(&flagReqBuilderHeader, "header", "H", []string{}, "HTTP header key-value for request. Syntax: \"Header: Value\"")
	reqNewCmd.Flags().StringVar(&flagReqBuilderBodyString, "data", "", "Text body for request.")
	reqNewCmd.Flags().StringVar(&flagReqBuilderBodyFile, "data-file", "", "Body for request from file as text.")
	reqNewCmd.Flags().StringVar(&flagReqBuilderBodyBinary, "data-binary", "", "Body for request from file as binary data.")

	reqNewCmd.Flags().StringVarP(&flagReqBuilderAuthBasic, "user", "U", "", "Basic Auth for request. Syntax: \"user:password\"")
	reqNewCmd.Flags().StringVarP(&flagReqBuilderAuthBearer, "bearer", "B", "", "Bearer Auth for request.")
	// reqNewCmd.Flags().StringVarP(&flagReqBuilderAuthApiKey, "api-key", "K", "", "API Key Auth for request. Syntax: \"X-API-Key:value\"")

	// TODO: not implemented.
	// reqNewCmd.Flags().BoolVarP(&flagReqBuilderCookieJar, "cookie-jar", "c", false, "Enable auto manage for cookies.")
	// reqNewCmd.Flags().BoolVarP(&flagReqBuilderCookieDisable, "no-cookie", "C", false, "Disable cookies.")
}

func readFileAsText(filename string) (string, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func readFileAsBytes(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
