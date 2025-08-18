package def

import (
	"fmt"
)

type exportedFunction struct {
	OriginalName string
	Rename       string
	Forwarder    string
	Ordinal      uint32
}

type DefFile struct {
	DllName           string
	exportedFunctions []exportedFunction
}

func (d *DefFile) AddExportedFunction(name string, ordinal uint32) {
	d.exportedFunctions = append(d.exportedFunctions, exportedFunction{
		OriginalName: name,
		Ordinal:      ordinal,
	})
}

func (d *DefFile) AddRenamedFunction(originalName string, rename string, ordinal uint32) {
	d.exportedFunctions = append(d.exportedFunctions, exportedFunction{
		OriginalName: originalName,
		Rename:       rename,
		Ordinal:      ordinal,
	})
}

func (d *DefFile) AddForwardedFunction(originalName string, forwarder string, ordinal uint32) {
	d.exportedFunctions = append(d.exportedFunctions, exportedFunction{
		OriginalName: originalName,
		Forwarder:    forwarder,
		Ordinal:      ordinal,
	})
}

func (d *DefFile) GetContent() string {
	var content string

	content += "LIBRARY \"" + d.DllName + "\"\n"
	content += "EXPORTS\n"

	for _, function := range d.exportedFunctions {
		if function.Forwarder == "" && function.Rename == "" {
			content += "\t" + function.OriginalName
		}

		if function.Forwarder != "" {
			// Forwarded functions
			content += "\t" + function.OriginalName + "=" + function.Forwarder
		}

		if function.Rename != "" {
			// Exported-renamed functions
			content += "\t" + function.OriginalName + "=" + function.Rename
		}

		// Add ordinals
		content += " " + "@" + fmt.Sprintf("%d", function.Ordinal)

		content += "\n"
	}

	return content
}
