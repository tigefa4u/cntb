/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"encoding/json"

	"contabo.com/cli/cntb/client"
	contaboCmd "contabo.com/cli/cntb/cmd"
	"contabo.com/cli/cntb/cmd/util"
	"contabo.com/cli/cntb/outputFormatter"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "instancesActions",
	Short: "History of your instance actions",
	Long:  `Show what actions you took on your instance`,
	Run: func(cmd *cobra.Command, args []string) {
		historyRequest := client.ApiClient().InstanceActionsAuditsApi.
			RetrieveInstancesActionsAuditsList(context.Background()).
			XRequestId(uuid.NewV4().String()).
			Page(contaboCmd.Page).
			Size(contaboCmd.Size).
			OrderBy([]string{contaboCmd.OrderBy})

		if cmd.Flags().Changed("instanceId") {
			historyRequest = historyRequest.InstanceId(instanceIdFilter)
		}

		if cmd.Flags().Changed("requestId") {
			historyRequest = historyRequest.RequestId(contaboCmd.RequestIdFilter)
		}

		if cmd.Flags().Changed("changedBy") {
			historyRequest = historyRequest.ChangedBy(contaboCmd.ChangedByFilter)
		}

		resp, httpResp, err := historyRequest.Execute()

		util.HandleErrors(err, httpResp, "while retrieving instance action history")

		responseJson, _ := json.Marshal(resp.Data)

		configFormatter := outputFormatter.FormatterConfig{
			Filter: []string{
				"id", "instanceId", "action", "username", "timestamp",
			},
			WideFilter: []string{
				"id", "instanceId", "action", "username", "changedBy", "requestId",
				"traceId", "tenantId", "customerId", "timestamp", "changes",
			},
			JsonPath: contaboCmd.OutputFormatDetails,
		}

		util.HandleResponse(responseJson, configFormatter)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		contaboCmd.ValidateOutputFormat()

		return nil
	},
}

func init() {
	contaboCmd.HistoryCmd.AddCommand(historyCmd)

	historyCmd.Flags().Int64VarP(&instanceIdFilter, "instanceId", "i", -1,
		`To filter audits using Instance Id`)
	viper.BindPFlag("instanceId", historyCmd.Flags().Lookup("instanceId"))
}