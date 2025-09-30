// Copyright © 2024 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"path"

	"github.com/Fraktal-PM3/firefly-cli/internal/constants"
	"github.com/Fraktal-PM3/firefly-cli/internal/docker"
	"github.com/Fraktal-PM3/firefly-cli/internal/log"
	"github.com/Fraktal-PM3/firefly-cli/internal/stacks"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:               "stop <stack_name>",
	Short:             "Stop a stack",
	Long:              `Stop a stack`,
	ValidArgsFunction: listStacks,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := log.WithVerbosity(context.Background(), verbose)
		ctx = log.WithLogger(ctx, logger)

		version, err := docker.CheckDockerConfig()
		if err != nil {
			return err
		}
		ctx = context.WithValue(ctx, docker.CtxComposeVersionKey{}, version)

		stackManager := stacks.NewStackManager(ctx)
		if len(args) == 0 {
			return fmt.Errorf("no stack specified")
		}
		stackName := args[0]
		
		stackDirectory, err := cmd.Flags().GetString("stack-dir")
		if err != nil {
			return err
		}		

		if err := stackManager.LoadStack(stackName, stackDirectory); err != nil {
			return err
		}

		fmt.Printf("stopping stack '%s'... ", stackName)
		if err := stackManager.StopStack(); err != nil {
			return err
		}
		fmt.Print("done\n")
		return nil
	},
}

func init() {
	stopCmd.Flags().StringP("stack-dir", "d", path.Join(constants.StacksDir), "Path to the stack directory")
	rootCmd.AddCommand(stopCmd)
}
