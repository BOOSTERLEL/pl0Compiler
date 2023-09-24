package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"pl0Compiler/build"
)

func main() {
	app := cli.NewApp()
	app.Name = "pl/0"
	app.Usage = "pl0 compiler is a tool for managing pl/0 source code."
	app.Version = "0.0.1-SNAPSHOT"

	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "clang", Usage: "set clang", Value: ""},
		&cli.BoolFlag{Name: "debug", Aliases: []string{"d"}, Usage: "set debug mode"},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "run",
			Usage: "compile and run pl/0 program",
			Action: func(c *cli.Context) error {
				ctx := build.NewContext(buildOptions(c))
				output, _ := ctx.Run(c.Args().First(), nil)
				fmt.Print(string(output))
				return nil
			},
		},
		{
			Name:  "build",
			Usage: "compile pl/0 source code",
			Action: func(c *cli.Context) error {
				ctx := build.NewContext(buildOptions(c))
				ctx.Build(c.Args().First(), nil, "a.out.exe")
				return nil
			},
		},
		{
			Name:  "lex",
			Usage: "lex pl/0 source code and print token list",
			Action: func(c *cli.Context) error {
				ctx := build.NewContext(buildOptions(c))
				tokens, comments, _ := ctx.Lex(c.Args().First(), nil)
				fmt.Println(tokens)
				fmt.Println(comments)
				return nil
			},
		},
		{
			Name:  "ast",
			Usage: "parse pl/0 source code and print ast",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "json", Usage: "output json format"},
			},
			Action: func(c *cli.Context) error {
				ctx := build.NewContext(buildOptions(c))
				f, err := ctx.AST(c.Args().First(), nil)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if c.Bool("json") {
					fmt.Println(f.JSONString())
				} else {
					fmt.Println(f.String())
				}
				return nil
			},
		},
		{
			Name:  "asm",
			Usage: "parse pl/0 source code and print llvm-ir",
			Action: func(c *cli.Context) error {
				ctx := build.NewContext(buildOptions(c))
				ll, _ := ctx.ASM(c.Args().First(), nil)
				fmt.Println(ll)
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func buildOptions(c *cli.Context) *build.Option {
	return &build.Option{
		Debug: c.Bool("debug"),
		Clang: c.String("clang"),
	}
}
