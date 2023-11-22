// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"github.com/unicornultrafoundation/go-u2u/libs/p2p/enode"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var (
	app = &cli.App{
		Name:        filepath.Base(os.Args[0]),
		Usage:       "go-u2u crawler",
		Version:     "v.0.0.1",
		Writer:      os.Stdout,
		HideVersion: true,
	}
)

func init() {
	app.Flags = append(app.Flags, Flags...)
	app.Before = func(ctx *cli.Context) error {
		return Setup(ctx)
	}
	// Add subcommands.
	app.Commands = []*cli.Command{
		apiCommand,
		crawlerCommand,
		discv4Command,
		discv5Command,
	}
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-127)
	}
}

// commandHasFlag returns true if the current command supports the given flag.
func commandHasFlag(ctx *cli.Context, flag cli.Flag) bool {
	names := flag.Names()
	set := make(map[string]struct{}, len(names))
	for _, name := range names {
		set[name] = struct{}{}
	}
	for _, fn := range ctx.FlagNames() {
		if _, ok := set[fn]; ok {
			return true
		}
	}
	return false
}

// getNodeArg handles the common case of a single node descriptor argument.
func getNodeArg(ctx *cli.Context) *enode.Node {
	if ctx.NArg() < 1 {
		exit("missing node as command-line argument")
	}
	n, err := parseNode(ctx.Args().First())
	if err != nil {
		exit(err)
	}
	return n
}

func exit(err interface{}) {
	if err == nil {
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
