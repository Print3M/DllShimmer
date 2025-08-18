package main

import (
	"dllshimmer/cli"
	"dllshimmer/dll"
	"dllshimmer/output"
	"embed"
	"fmt"
	"path/filepath"
)

//go:embed templates/*
var templatesFS embed.FS

func main() {
	flags := cli.ParseCli()

	cli.PrintBanner()

	out := output.Output{
		Dll:         dll.ParseDll(flags.Input, flags.OriginalPath),
		OutputDir:   filepath.Clean(flags.Output),
		TemplatesFS: &templatesFS,
	}

	out.CreateCodeFile(flags.Mutex, flags.Static)
	out.CreateDefFile()
	out.CreateCompileScript(flags.Static)

	if flags.Static {
		out.CreateLibFile()
	}

	fmt.Println()
	fmt.Println("Success! What to do next?")
	fmt.Println()
	fmt.Printf("  1. Jump into the '%s/' directory.\n", out.OutputDir)
	fmt.Printf("  2. Add your backdoor to the '%s' file.\n", out.GetCodeFileName())
	fmt.Printf("  3. Compile project using the '%s' script.\n", out.GetCompileScriptName())
	fmt.Println()
}
