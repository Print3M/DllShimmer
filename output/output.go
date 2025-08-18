package output

import (
	"dllshimmer/def"
	"dllshimmer/dll"
	embed "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type Output struct {
	Dll         *dll.Dll
	OutputDir   string
	TemplatesFS *embed.FS
}

func (o *Output) GetDefFileName() string {
	return o.Dll.Name + ".def"
}

func (o *Output) GetCodeFileName() string {
	return o.Dll.Name + ".cpp"
}

func (o *Output) GetCompileScriptName() string {
	return "compile.sh"
}

func (o *Output) GetOutputDllName() string {
	return o.Dll.Name
}

func (o *Output) GetLibFileName() string {
	return "original.lib"
}

type CodeFileParams struct {
	Functions    []dll.ExportedFunction
	OriginalPath string
	Mutex        bool
	DllName      string
}

func (o *Output) GetTemplate(filename string) *template.Template {
    path := "templates/" + filename
    content, err := o.TemplatesFS.ReadFile(path)
    if err != nil {
        log.Fatalf("[!] Error while reading embedded template file '%s': %v", path, err)
    }
    return template.Must(template.New("new").Parse(string(content)))
}


func (o *Output) CreateCodeFile(mutex bool, isStaticLinked bool) {
	templateFile := "dynamic-shim.cpp.template"
	if isStaticLinked {
		templateFile = "static-shim.cpp.template"
	}

	tmpl := o.GetTemplate(templateFile)
	outputPath := filepath.Join(o.OutputDir, o.GetCodeFileName())

	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("[!] Error while creating '%s' file: %v", outputPath, err)
	}
	defer f.Close()

	params := CodeFileParams{
		Functions:    o.Dll.ExportedFunctions,
		OriginalPath: sanitizePathForInjection(o.Dll.OriginalPath),
		Mutex:        mutex,
		DllName:      o.Dll.Name,
	}

	err = tmpl.Execute(f, params)
	if err != nil {
		log.Fatalf("[!] Error of template engine: %v", err)
	}

	fmt.Printf("[+] '%s' file created\n", outputPath)
}

func (o *Output) CreateDefFile() {
	var def def.DefFile
	def.DllName = o.Dll.Name

	for _, function := range o.Dll.ExportedFunctions {
		if function.Forwarder == "" {
			def.AddRenamedFunction(function.Name, function.Name+"Fwd", function.Ordinal)
		} else {
			def.AddForwardedFunction(function.Name, function.Forwarder, function.Ordinal)
		}
	}

	content := def.GetContent()
	outputPath := filepath.Join(o.OutputDir, o.GetDefFileName())

	err := os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		log.Fatalf("[!] Error while creating '%s' file: %v", outputPath, err)
	}

	fmt.Printf("[+] '%s' file created\n", outputPath)
}

func (o *Output) CreateLibFile() {
	var def def.DefFile

	// In case of static linking OriginalPath is DLL name itself
	def.DllName = o.Dll.OriginalPath

	for _, function := range o.Dll.ExportedFunctions {
		if function.Forwarder == "" {
			def.AddExportedFunction(function.Name, function.Ordinal)
		} else {
			def.AddForwardedFunction(function.Name, function.Forwarder, function.Ordinal)
		}
	}

	temp, err := os.CreateTemp("", "dllshimmer-*.def")
	if err != nil {
		panic(err)
	}
	defer os.Remove(temp.Name())

	content := def.GetContent()
	err = os.WriteFile(temp.Name(), []byte(content), 0644)
	if err != nil {
		log.Fatalf("[!] Error while creating '%s' file: %v", temp.Name(), err)
	}

	// Convert DLL to .lib file
	outputPath := filepath.Join(o.OutputDir, o.GetLibFileName())
	cmd := exec.Command("dlltool", "-d", temp.Name(), "-l", outputPath, "-m", "i386:x86-64")

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}

	fmt.Printf("[+] '%s' file created\n", outputPath)
}

type CompileScriptParams struct {
	Code           string
	Def            string
	Output         string
	IsStaticLinked bool
}

func (o *Output) CreateCompileScript(isStaticLinked bool) {
	tmpl := o.GetTemplate("compile.sh.template")
	outputPath := filepath.Join(o.OutputDir, o.GetCompileScriptName())

	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("[!] Error while creating '%s' file: %v", outputPath, err)
	}
	defer f.Close()

	params := CompileScriptParams{
		Code:           o.GetCodeFileName(),
		Def:            o.GetDefFileName(),
		Output:         o.GetOutputDllName(),
		IsStaticLinked: isStaticLinked,
	}

	err = tmpl.Execute(f, params)
	if err != nil {
		log.Fatalf("[!] Error of template engine: %v", err)
	}

	fmt.Printf("[+] '%s' file created\n", outputPath)
}
