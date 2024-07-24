/*
Copyright 2024 The Kubernetes Authors.

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

package apiresources

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
)

type APIResourcesPrinter struct {
	outputFormat string
	noHeaders    bool
}

type PrintFlags struct {
	JSONYamlPrintFlags *genericclioptions.JSONYamlPrintFlags
	OutputFormat       *string
}

func (f *PrintFlags) AllowedFormats() []string {
	formats := f.JSONYamlPrintFlags.AllowedFormats()
	formats = append(formats, []string{"name", "wide"}...)
	return formats
}

func (f *PrintFlags) ToPrinter() (printers.ResourcePrinter, error) {
	outputFormat := ""
	if f.OutputFormat != nil {
		outputFormat = *f.OutputFormat
	}

	p, err := f.JSONYamlPrintFlags.ToPrinter(outputFormat)

	if err != nil && !genericclioptions.IsNoCompatiblePrinterError(err) {
		return nil, err
	}

	return p, nil
}

func (f *PrintFlags) AddFlags(cmd *cobra.Command) {
	f.JSONYamlPrintFlags.AddFlags(cmd)
	if f.OutputFormat != nil {
		cmd.Flags().StringVarP(f.OutputFormat, "output", "o", *f.OutputFormat, fmt.Sprintf("Output format. One of: (%s)", strings.Join(f.AllowedFormats(), ", ")))
	}
}

func NewPrintFlags() *PrintFlags {
	outputFormat := ""
	return &PrintFlags{
		OutputFormat:       &outputFormat,
		JSONYamlPrintFlags: genericclioptions.NewJSONYamlPrintFlags(),
	}
}

func (p APIResourcesPrinter) PrintObj(_ runtime.Object, writer io.Writer) error {
	if p.noHeaders == false && p.outputFormat != "name" {
		if err := printContextHeaders(writer, p.outputFormat); err != nil {
			return err
		}
	}
	for _, r := range resources {
		switch p.outputFormat {
		case "name":
			name := r.APIResource.Name
			if len(r.APIGroup) > 0 {
				name += "." + r.APIGroup
			}
			if _, err := fmt.Fprintf(writer, "%s\n", name); err != nil {
				//errs = append(errs, err)
				return err
			}
		case "wide":
			if _, err := fmt.Fprintf(writer, "%s\t%s\t%s\t%v\t%s\t%v\t%v\n",
				r.APIResource.Name,
				strings.Join(r.APIResource.ShortNames, ","),
				r.APIGroupVersion,
				r.APIResource.Namespaced,
				r.APIResource.Kind,
				strings.Join(r.APIResource.Verbs, ","),
				strings.Join(r.APIResource.Categories, ",")); err != nil {
				//errs = append(errs, err)
				return err
			}
		case "":
			if _, err := fmt.Fprintf(writer, "%s\t%s\t%s\t%v\t%s\n",
				r.APIResource.Name,
				strings.Join(r.APIResource.ShortNames, ","),
				r.APIGroupVersion,
				r.APIResource.Namespaced,
				r.APIResource.Kind); err != nil {
				//errs = append(errs, err)
				return err
			}

		}
	}
	return nil
}

func NewAPIResourcesPrinter(format string, noHeaders bool) *APIResourcesPrinter {
	return &APIResourcesPrinter{
		outputFormat: format,
		noHeaders:    noHeaders,
	}
}
