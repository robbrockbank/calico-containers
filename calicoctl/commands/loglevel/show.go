// Copyright (c) 2016-2020 Tigera, Inc. All rights reserved.

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

package loglevel

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/projectcalico/calicoctl/calicoctl/commands/constants"
	"github.com/projectcalico/libcalico-go/lib/clientv3"

	docopt "github.com/docopt/docopt-go"

	"github.com/projectcalico/calicoctl/calicoctl/commands/clientmgr"
)

func produceTable(ctx context.Context, configurationSettings []*componentConfiguration, component string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"COMPONENT", "NODE", "LOG LEVEL"})
	genRow := func(component, node, logLevel string) []string {
		return []string{
			component,
			node,
			logLevel,
		}
	}

	var rows [][]string
	for _, c := range configurationSettings {
		componentName := string(c.component)
		if component != "" {
			if component == componentName {
				rows = append(rows, genRow(componentName, c.node, string(c.logLevel)))
			}
		} else {
			rows = append(rows, genRow(componentName, c.node, string(c.logLevel)))
		}
	}
	table.AppendBulk(rows)
	table.Render()
}

func showComponentConfiguration(ctx context.Context, debuggingConfigurationClient clientv3.DebuggingConfigurationInterface,
	component string) error {

	configurationSettings, _, err := collectLogLevelConfiguration(ctx, debuggingConfigurationClient)
	if err != nil {
		return err
	}

	produceTable(ctx, configurationSettings, component)

	return nil
}

func showConfiguration(ctx context.Context, debuggingConfigurationClient clientv3.DebuggingConfigurationInterface) error {
	configurationSettings, _, err := collectLogLevelConfiguration(ctx, debuggingConfigurationClient)
	if err != nil {
		return err
	}

	produceTable(ctx, configurationSettings, "")

	return nil
}

// Show displays the log level configuration per component/node.
func Show(args []string) error {
	doc := constants.DatastoreIntro + `Usage:
  calicoctl log-level show [--component=<component>] [--config=<CONFIG>]

Options:
  -h --help                     Show this screen.
     --component=<component>    Show log level configuration for the specified component.
  -c --config=<CONFIG>          Path to the file containing connection configuration in
                                YAML or JSON format.
                                [default: ` + constants.DefaultConfigPath + `]

Description:
  The log-level show command prints log level information about a given calico component, or about
  overall calico components.
`
	parsedArgs, err := docopt.Parse(doc, args, true, "", false, false)
	if err != nil {
		return fmt.Errorf("Invalid option: 'calicoctl %s'. Use flag '--help' to read about a specific subcommand.", strings.Join(args, " "))
	}
	if len(parsedArgs) == 0 {
		return nil
	}
	ctx := context.Background()

	// Create a new backend client from env vars.
	cf := parsedArgs["--config"].(string)
	client, err := clientmgr.NewClient(cf)
	if err != nil {
		return err
	}

	debuggingConfigurationClient := client.DebuggingConfiguration()

	passedComponent := parsedArgs["--component"]
	if passedComponent != nil {
		if valid, errorMsg := isValidComponent(ctx, passedComponent.(string)); !valid {
			return fmt.Errorf(errorMsg)
		}
		return showComponentConfiguration(ctx, debuggingConfigurationClient, passedComponent.(string))
	}

	return showConfiguration(ctx, debuggingConfigurationClient)
}
