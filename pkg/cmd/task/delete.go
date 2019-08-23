// Copyright © 2019 The Tekton Authors.
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

package task

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cliopts "k8s.io/cli-runtime/pkg/genericclioptions"
)

func deleteCommand(p cli.Params) *cobra.Command {
	f := cliopts.NewPrintFlags("delete")
	eg := `
# Delete a Task of name 'foo' in namespace 'bar'
tkn task delete foo -n bar

tkn t rm foo -n bar",
`

	c := &cobra.Command{
		Use:          "delete",
		Aliases:      []string{"rm"},
		Short:        "Delete a task resource in a namespace",
		Example:      eg,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		Annotations: map[string]string{
			"commandType": "main",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := &cli.Stream{
				Out: cmd.OutOrStdout(),
				Err: cmd.OutOrStderr(),
			}
			return deleteTask(s, p, args[0])
		},
	}
	f.AddFlags(c)
	return c
}

func deleteTask(s *cli.Stream, p cli.Params, tName string) error {
	cs, err := p.Clients()
	if err != nil {
		return fmt.Errorf("Failed to create tekton client")
	}

	if err := cs.Tekton.TektonV1alpha1().Tasks(p.Namespace()).Delete(tName, &metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("Failed to delete task %q: %s", tName, err)
	}

	fmt.Fprintf(s.Out, "Task deleted: %s\n", tName)
	return nil
}
