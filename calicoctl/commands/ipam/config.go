// Copyright (c) 2016 Tigera, Inc. All rights reserved.

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

package ipam

import (
	"errors"
	"os"
	"strings"

	"fmt"

	"github.com/docopt/docopt-go"
	"github.com/projectcalico/calicoctl/calicoctl/commands/argutils"
	"github.com/projectcalico/calicoctl/calicoctl/commands/clientmgr"
	"github.com/projectcalico/calicoctl/calicoctl/commands/constants"
	"github.com/projectcalico/libcalico-go/lib/client"
	"github.com/projectcalico/libcalico-go/lib/numorstring"
	"strconv"
)

func Config(args []string) {
	doc := constants.DatastoreIntro + `Usage:
  calicoctl ipam config set <NAME> <VALUE> [--config=<CONFIG>]
  calicoctl ipam config get <NAME> [--config=<CONFIG>]

Examples:
  # Turn off the full BGP node-to-node mesh
  calicoctl config set nodeToNodeMesh off

  # Set global log level to warning
  calicoctl config set logLevel warning

  # Set log level to info for node "node1"
  calicoctl config set logLevel info --node=node1

  # Display the current setting for the nodeToNodeMesh
  calicoctl config get nodeToNodeMesh

Options:
  -c --config=<CONFIG>  Path to the file containing connection configuration in
                        YAML or JSON format.
                        [default: ` + constants.DefaultConfigPath + `]

Description:

These commands can be used to manage global IPAM settings.  These values can
only be modified when there are no IP reservations.  It is expected that the
IPAM configuration is setup once during system installation.

The table below details the valid config options.

 Name               | Value | Description
--------------------+-------+-------------------------------------------------
 blockSizeIPv4      | 1 - 8 | The size (in bits) of the sub-blocks allocated
                    |       | to a host from the available IPv4 Pools.
                    |       | [Default: 6, i.e. 64 addresses]
 blockSizeIPv4      | 1 - 8 | The size (in bits) of the sub-blocks allocated
                    |       | to a host from the available IPv6 Pools.
                    |       | [Default: 6, i.e. 64 addresses]
`
	parsedArgs, err := docopt.Parse(doc, args, true, "calicoctl", false, false)
	if err != nil {
		fmt.Printf("Invalid option: 'calicoctl %s'. Use flag '--help' to read about a specific subcommand.\n", strings.Join(args, " "))
		os.Exit(1)
	}
	if len(parsedArgs) == 0 {
		return
	}

	// Load the client config and connect.
	cf := parsedArgs["--config"].(string)
	client, err := clientmgr.NewClient(cf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Load the current IPAMConfig.
	ipamConfig, err := client.IPAM().GetIPAMConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	name := argutils.ArgStringOrBlank(parsedArgs, "<NAME>")

	// Handle get.
	if parsedArgs["get"].(bool) {
		switch strings.ToLower(name) {
		case "blocksizeipv4":
			fmt.Println(ipamConfig.AffineBlockBitSizeIPv4)
			return
		case "blocksizeipv6":
			fmt.Println(ipamConfig.AffineBlockBitSizeIPv6)
			return
		default:
			fmt.Printf("Error executing command: unrecognised config name '%s'\n", name)
			os.Exit(1)
		}
	}

	// Handle set.  We currently only support setting the block size - so verify
	// the supplied value is valid.
	value := argutils.ArgStringOrBlank(parsedArgs, "<VALUE>")
	numValue, err := strconv.ParseInt(value, 10, 8)
	if err != nil || numValue < 1 || numValue > 8 {
		fmt.Printf("Error executing command: config value is invalid '%s'\n", value)
		os.Exit(1)
	}

	// Update the config parameter based on the config name.
	switch strings.ToLower(name) {
	case "blocksizeipv4":
		ipamConfig.AffineBlockBitSizeIPv4 = int(numValue)
	case "blocksizeipv6":
		ipamConfig.AffineBlockBitSizeIPv6 = int(numValue)
	default:
		fmt.Printf("Error executing command: unrecognised config name '%s'\n", name)
		os.Exit(1)
	}

	// Apply the config.
	err = client.IPAM().SetIPAMConfig(*ipamConfig)
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		os.Exit(1)
	}

	return
}
