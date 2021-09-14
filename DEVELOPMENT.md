# Development

## Tools

### Go

This repo is written in the [Go](https://golang.org) programming language, which can be installed from their [installation page](https://golang.org/doc/install). The version it was developed with was v1.16, which is the minimum required version for this SDK.

## Shellspec
The SDK developed in this repo is tested by [shellspec](https://github.com/shellspec/shellspec), which can be installed as documented in [their README](https://github.com/shellspec/shellspec#installation).

## Golangci-lint
The SDK is linted by [golangci-lint](https://github.com/golangci/golangci-lint), which can be installed as documented in their [local installation page](https://golangci-lint.run/usage/install/#local-installation).

## How to add a new command
1. Add a new file with the name of the command under `cmd`
2. The minimal file should contain:
```go
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/snyk/snyk-iac-custom-rules/internal"
	"github.com/snyk/snyk-iac-custom-rules/util"
)

var <command>Command = &cobra.Command{
	Use:   "<command>",
	Short: "",
	Long: ``,
	SilenceUsage: true, // disables help from being printed if command fails
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.Run<Command>(args, <command>Params)
	},
}

func new<Command>CommandParams() *internal.<Command>CommandParams {
	return &internal.<Command>CommandParams{
	}
}

var <command>Params = new<Command>CommandParams()

func init() {
	// initialise flags for the command
	RootCommand.AddCommand(<command>Command)
}
```

3. Add a new file with the name of the command under `internal` - this will contain the implementation for `Run<Command>`

```go
package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type <Command>CommandParams struct {
}

func Run<Command>(args []string, params *<Command>CommandParams) error {
	return nil
}
```

4. Add a basic test to with the name of the command under `spec`:

```go
#!/bin/bash
Describe './snyk-iac-custom-rules <command>'
   It 'returns passing test status'
      When call ./snyk-iac-custom-rules <command>
      The status should be success
      The output should include ''
   End
End
```