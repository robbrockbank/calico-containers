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
	"strings"

	"github.com/projectcalico/calicoctl/calicoctl/commands/constants"
	v3 "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/clientv3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	docopt "github.com/docopt/docopt-go"

	"github.com/projectcalico/calicoctl/calicoctl/commands/clientmgr"
	"github.com/projectcalico/libcalico-go/lib/options"
)

func setConfiguration(ctx context.Context, debuggingConfigurationClient clientv3.DebuggingConfigurationInterface,
	component string, node string, logLevel string) error {

	configurationSettings, objectMeta, err := collectLogLevelConfiguration(ctx, debuggingConfigurationClient)
	if err != nil {
		return err
	}

	found := false
	spec := make([]v3.ComponentConfiguration, len(configurationSettings))

	for i, c := range configurationSettings {
		if string(c.component) == component && c.node == node {
			found = true
			spec[i] = v3.ComponentConfiguration{Component: c.component, Node: c.node, LogSeverity: v3.LogLevel(logLevel)}
		} else {
			spec[i] = v3.ComponentConfiguration{Component: c.component, Node: c.node, LogSeverity: c.logLevel}
		}
	}

	if !found {
		spec = append(spec, v3.ComponentConfiguration{Component: v3.Component(component), Node: node, LogSeverity: v3.LogLevel(logLevel)})
	}

	dc := &v3.DebuggingConfiguration{
		ObjectMeta: v1.ObjectMeta{
			Name: "default",
		},
		Spec: v3.DebuggingConfigurationSpec{
			Configuration: spec,
		},
	}

	if objectMeta.ResourceVersion == "" {
		_, err = debuggingConfigurationClient.Create(ctx, dc, options.SetOptions{})
	} else {
		dc.ObjectMeta = objectMeta
		dc.ObjectMeta.CreationTimestamp = v1.Now()
		_, err = debuggingConfigurationClient.Update(ctx, dc, options.SetOptions{})
	}
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	return nil
}

// Set configures log level setting per component.
func Set(args []string) error {
	doc := constants.DatastoreIntro + `Usage:
  calicoctl log-level set --component=<component> [--node=<node>] --severity=<Info|Debug> [--config=<CONFIG>]

Options:
  -h --help                     Show this screen.
     --component=<component>    Set log level severity for the specified component.
     --node=<node>              Set log level severity for the specified component running on the specified node.
     --severity=<info|debug>    Indicates the log severity to be set.
  -c --config=<CONFIG>          Path to the file containing connection configuration in
                                YAML or JSON format.
                                [default: ` + constants.DefaultConfigPath + `]

Description:
  The log-level set command allows setting log severity for a given calico component, or for a given component on
  a given node.
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

	passedNode := ""
	passedComponent := parsedArgs["--component"]
	ps := parsedArgs["--node"]
	if ps != nil {
		passedNode = parsedArgs["--node"].(string)
	}
	logSeverity := parsedArgs["--severity"]

	// Avoid making a call to get DebuggingConfiguration and one to try to set it with newer valid
	// since we already have an utility that validates component is valid.
	if valid, errorMsg := isValidComponent(ctx, passedComponent.(string)); !valid {
		return fmt.Errorf(errorMsg)
	}

	return setConfiguration(ctx, debuggingConfigurationClient, passedComponent.(string), passedNode, logSeverity.(string))
}
