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
	"sort"

	v3 "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/clientv3"
	"github.com/projectcalico/libcalico-go/lib/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/projectcalico/libcalico-go/lib/options"
)

type componentConfiguration struct {
	component v3.Component
	node      string
	logLevel  v3.LogLevel
}

// byComponent sorts componentConfiguration by name. When node is specified is higher priority
// than when it is not
type byComponent []*componentConfiguration

func (c byComponent) Len() int      { return len(c) }
func (c byComponent) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c byComponent) Less(i, j int) bool {
	if c[i].component == c[j].component {
		if c[i].node != "" && c[j].node != "" {
			return c[i].node < c[j].node
		} else if c[i].node != "" {
			return true
		} else {
			return false
		}
	}
	return c[i].component < c[j].component
}

func collectLogLevelConfiguration(ctx context.Context, debuggingConfigurationClient clientv3.DebuggingConfigurationInterface) ([]*componentConfiguration, v1.ObjectMeta, error) {
	dc, err := debuggingConfigurationClient.Get(ctx, "default", options.GetOptions{})
	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			return nil, v1.ObjectMeta{}, nil
		}

		return nil, v1.ObjectMeta{}, err
	}

	configuration := dc.Spec.Configuration

	configurationSettings := make([]*componentConfiguration, len(configuration))

	for i, c := range configuration {
		configurationSettings[i] = &componentConfiguration{component: c.Component, logLevel: c.LogSeverity, node: c.Node}
	}

	// Sort this by component name first. Component/node is higher priority than Component
	sort.Sort(byComponent(configurationSettings))

	return configurationSettings, dc.ObjectMeta, nil
}

func isValidComponent(ctx context.Context, component string) (bool, string) {
	return v3.IsValidDebuggingConfigurationComponent(v3.Component(component))
}
