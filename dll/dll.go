package dll

import (
	"log"
	"path/filepath"

	peparser "github.com/saferwall/pe"
)

type ExportedFunction struct {
	Name      string
	Forwarder string
	Ordinal   uint32
}

type Dll struct {
	Name              string
	OriginalPath      string
	ExportedFunctions []ExportedFunction
}

func ParseDll(path string, originalPath string) *Dll {
	var dll Dll

	dll.Name = filepath.Base(path)
	dll.OriginalPath = originalPath

	pe, err := peparser.New(path, &peparser.Options{})
	if err != nil {
		log.Fatalf("[!] Error while opening file: %s, reason: %v", path, err)
	}

	err = pe.Parse()
	if err != nil {
		log.Fatalf("[!] Error while parsing file: %s, reason: %v", path, err)
	}

	for _, function := range pe.Export.Functions {
		dll.ExportedFunctions = append(dll.ExportedFunctions, ExportedFunction{
			Name:      function.Name,
			Forwarder: function.Forwarder,
			Ordinal:   function.Ordinal,
		})
	}

	return &dll
}
