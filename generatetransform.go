/*
 * MALTEGO - Go package that provides datastructures for interacting with the Maltego graphical link analysis tool.
 * Copyright (c) 2021 Philipp Mieden <dreadl0ck [at] protonmail [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package maltego

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

const (
	// can be used to set the debug mode for all generated maltego transforms.
	transformDebug = false

	configFileExtension = ".mtz"
)

// Transforms

// MaltegoTransform models a maltego transformation when exported as configuration.
type MaltegoTransform struct {
	XMLName            xml.Name `xml:"MaltegoTransform"`
	Name               string   `xml:"name,attr"`
	DisplayName        string   `xml:"displayName,attr"`
	Abstract           bool     `xml:"abstract,attr"`
	Template           bool     `xml:"template,attr"`
	Visibility         string   `xml:"visibility,attr"`
	Description        string   `xml:"description,attr"`
	Author             string   `xml:"author,attr"`
	RequireDisplayInfo bool     `xml:"requireDisplayInfo,attr"`

	TransformAdapter string                 `xml:"TransformAdapter"`
	Properties       XMLTransformProperties `xml:"Properties"`
	Constraints      InputConstraints       `xml:"InputConstraints"`
	OutputEntities   string                 `xml:"OutputEntities"`
	DefaultSets      defaultSets            `xml:"defaultSets"`
	StealthLevel     string                 `xml:"StealthLevel"`
}

type defaultSets struct {
	Items []Set `xml:"Set"`
}

type Set struct {
	XMLName xml.Name `xml:"Set"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name,attr"`
}

type XMLTransformProperties struct {
	XMLName xml.Name `xml:"Properties"`
	Text    string   `xml:",chardata"`
	Fields  struct {
		Text     string     `xml:",chardata"`
		Property []Property `xml:"Property"`
	} `xml:"Fields"`
}

// Property structure
type Property struct {
	Text         string `xml:",chardata"`
	Name         string `xml:"name,attr"`
	Type         string `xml:"type,attr"`
	Nullable     bool   `xml:"nullable,attr"`
	Hidden       bool   `xml:"hidden,attr"`
	Readonly     bool   `xml:"readonly,attr"`
	Description  string `xml:"description,attr"`
	Popup        bool   `xml:"popup,attr"`
	Abstract     bool   `xml:"abstract,attr"`
	Visibility   string `xml:"visibility,attr"`
	Auth         bool   `xml:"auth,attr"`
	DisplayName  string `xml:"displayName,attr"`
	DefaultValue string `xml:"DefaultValue,omitempty"`
	SampleValue  string `xml:"SampleValue"`
}

// InputConstraints structure
type InputConstraints struct {
	XMLName xml.Name `xml:"InputConstraints"`
	Text    string   `xml:",chardata"`
	Entity  struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
		Min  int    `xml:"min,attr"`
		Max  int    `xml:"max,attr"`
	} `xml:"Entity"`
}

// TransformCoreInfo describes basic information needed to create a transform.
type TransformCoreInfo struct {
	ID          string // e.g ToAuditRecords
	InputEntity string
	Description string
}

// Settings

type TransformSettingProperty struct {
	XMLName xml.Name `xml:"Property"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Popup   bool     `xml:"popup,attr"`
}

type TransformSettingProperties struct {
	Items []TransformSettingProperty `xml:"Properties"`
}

// TransformSettings structure
type TransformSettings struct {
	XMLName            xml.Name                   `xml:"TransformSettings"`
	Text               string                     `xml:",chardata"`
	Enabled            bool                       `xml:"enabled,attr"`
	DisclaimerAccepted bool                       `xml:"disclaimerAccepted,attr"`
	ShowHelp           bool                       `xml:"showHelp,attr"`
	RunWithAll         bool                       `xml:"runWithAll,attr"`
	Favorite           bool                       `xml:"favorite,attr"`
	Property           TransformSettingProperties `xml:"Properties"`
}

type Server struct {
	XMLName     xml.Name `xml:"MaltegoServer"`
	Text        string   `xml:",chardata"`
	Name        string   `xml:"name,attr"`
	Enabled     bool     `xml:"enabled,attr"`
	Description string   `xml:"description,attr"`
	URL         string   `xml:"url,attr"`
	LastSync    string   `xml:"LastSync"`
	Protocol    struct {
		Text    string `xml:",chardata"`
		Version string `xml:"version,attr"`
	} `xml:"Protocol"`
	Authentication struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
	} `xml:"Authentication"`
	Transforms struct {
		Text      string `xml:",chardata"`
		Transform []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"Transform"`
	} `xml:"Transforms"`
	Seeds string `xml:"Seeds"`
}

type TransformSet struct {
	XMLName     xml.Name `xml:"TransformSet"`
	Text        string   `xml:",chardata"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:"description,attr"`
	Transforms  struct {
		Text      string `xml:",chardata"`
		Transform []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"Transform"`
	} `xml:"Transforms"`
}

// e.g. "ToAuditRecords" -> "To Audit Records [NETCAP]".
func ToTransformDisplayName(in string) string {
	var b strings.Builder

	for i, c := range in {
		switch {
		// if current char is upper case, but the previous is lowercase
		case i > 0 && unicode.IsUpper(c) && unicode.IsLower(rune(in[i-1])):

			b.WriteRune(' ')
			b.WriteRune(c)

		// if current char is upper case, and the next is Lowercase
		case unicode.IsUpper(c) && len(in) > i+1 && unicode.IsLower(rune(in[i+1])):

			// if the next char is followed by an uppercase char
			// or the string ends
			if len(in) > i+2 && unicode.IsUpper(rune(in[i+2])) || len(in) == i+2 {
				b.WriteRune(c)

				continue
			}

			b.WriteRune(' ')
			b.WriteRune(c)

		// else
		default:
			b.WriteRune(c)
		}
	}
	return strings.TrimSpace(b.String() + " [NETCAP]")
}

func NewTransformSettings(id string, debug bool, executable string) TransformSettings {
	trs := TransformSettings{
		Enabled:            true,
		DisclaimerAccepted: false,
		ShowHelp:           true,
		RunWithAll:         true,
		Favorite:           false,
		Property: TransformSettingProperties{
			Items: []TransformSettingProperty{
				{
					Name:  "transform.local.command",
					Type:  "string",
					Popup: false,
					Text:  executable,
				},
				{
					Name:  "transform.local.parameters",
					Type:  "string",
					Popup: false,
					Text:  "transform " + id,
				},
				{
					Name:  "transform.local.working-directory",
					Type:  "string",
					Popup: false,
					Text:  "/usr/local/",
				},
				{
					Name:  "transform.local.debug",
					Type:  "boolean",
					Popup: false,
					Text:  strconv.FormatBool(debug),
				},
			},
		},
	}

	return trs
}

func NewTransform(author, prefix, id string, description string, input string) MaltegoTransform {
	tr := MaltegoTransform{
		Name:               prefix + id,
		DisplayName:        ToTransformDisplayName(id),
		Abstract:           false,
		Template:           false,
		Visibility:         "public",
		Description:        description,
		Author:             author,
		RequireDisplayInfo: false,
		TransformAdapter:   "com.paterva.maltego.transform.protocol.v2api.LocalTransformAdapterV2",
		Properties: XMLTransformProperties{
			Fields: struct {
				Text     string     `xml:",chardata"`
				Property []Property `xml:"Property"`
			}{
				Property: []Property{
					// <Property name="transform.local.command" type="string" nullable="false" hidden="false" readonly="false" description="The command to execute for this transform" popup="false" abstract="false" visibility="public" auth="false" displayName="Command line">
					// <SampleValue></SampleValue>
					// </Property>
					{
						Text:         "",
						Name:         "transform.local.command",
						Type:         "string",
						Nullable:     false,
						Hidden:       false,
						Readonly:     false,
						Description:  "The command to execute for this transform",
						Popup:        false,
						Abstract:     false,
						Visibility:   "public",
						Auth:         false,
						DisplayName:  "Command line",
						SampleValue:  "",
						DefaultValue: "",
					},
					// <Property name="transform.local.parameters" type="string" nullable="true" hidden="false" readonly="false" description="The parameters to pass to the transform command" popup="false" abstract="false" visibility="public" auth="false" displayName="Command parameters">
					// <SampleValue></SampleValue>
					// </Property>
					{
						Text:         "",
						Name:         "transform.local.parameters",
						Type:         "string",
						Nullable:     true,
						Hidden:       false,
						Readonly:     false,
						Description:  "The parameters to pass to the transform command",
						Popup:        false,
						Abstract:     false,
						Visibility:   "public",
						Auth:         false,
						DisplayName:  "Command parameters",
						SampleValue:  "",
						DefaultValue: "",
					},
					// <Property name="transform.local.working-directory" type="string" nullable="true" hidden="false" readonly="false" description="The working directory used when invoking the executable" popup="false" abstract="false" visibility="public" auth="false" displayName="Working directory">
					// <DefaultValue>/</DefaultValue>
					// <SampleValue></SampleValue>
					// </Property>
					{
						Text:         "",
						Name:         "transform.local.working-directory",
						Type:         "string",
						Nullable:     true,
						Hidden:       false,
						Readonly:     false,
						Description:  "The working directory used when invoking the executable",
						Popup:        false,
						Abstract:     false,
						Visibility:   "public",
						Auth:         false,
						DisplayName:  "Working directory",
						SampleValue:  "",
						DefaultValue: "/",
					},
					// <Property name="transform.local.debug" type="boolean" nullable="true" hidden="false" readonly="false" description="When this is set, the transform&apos;s text output will be printed to the output window" popup="false" abstract="false" visibility="public" auth="false" displayName="Show debug info">
					// <SampleValue>false</SampleValue>
					// </Property>
					{
						Text:         "",
						Name:         "transform.local.debug",
						Type:         "boolean",
						Nullable:     true,
						Hidden:       false,
						Readonly:     false,
						Description:  "When this is set, the transform&apos;s text output will be printed to the output window",
						Popup:        false,
						Abstract:     false,
						Visibility:   "public",
						Auth:         false,
						DisplayName:  "Show debug info",
						SampleValue:  "false",
						DefaultValue: "",
					},
				},
			},
		},
		Constraints: InputConstraints{
			Entity: struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
				Min  int    `xml:"min,attr"`
				Max  int    `xml:"max,attr"`
			}{
				Text: "",
				Type: input,
				Min:  1,
				Max:  1,
			},
		},
		OutputEntities: "",
		DefaultSets: defaultSets{Items: []Set{
			{
				Name: "NETCAP",
			},
		}},
		StealthLevel: "0",
	}

	return tr
}

func GenTransform(author, prefix string, outDir string, name string, description string, inputEntity string, executable string) {
	var (
		tr  = NewTransform(author, prefix, name, description, inputEntity)
		trs = NewTransformSettings(strings.ToLower(string(name[0]))+name[1:], transformDebug, executable)
	)

	// write Transform

	data, err := xml.MarshalIndent(tr, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filepath.Join(outDir, "TransformRepositories", "Local", prefix+name+".transform"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// write TransformSettings

	data, err = xml.MarshalIndent(trs, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	f, err = os.Create(filepath.Join(outDir, "TransformRepositories", "Local", prefix+name+".transformsettings"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Directory structure:
// .
// ├── Servers
// │      └── Local.tas
// ├── TransformRepositories
// │      └── Local
// │          ├── corp.Transform1.transform
// │          ├── corp.Transform1.transformsettings
// │          ├── ...
// │          └── ...
// └── version.properties.
func GenTransformArchive() {
	// clean
	_ = os.RemoveAll("transforms")

	// create directories
	_ = os.MkdirAll("transforms/Servers", 0o700)
	_ = os.MkdirAll("transforms/TransformRepositories/Local", 0o700)

	fVersion, err := os.Create("transforms/version.properties")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if errClose := fVersion.Close(); errClose != nil {
			fmt.Println(errClose)
		}
	}()

	_, _ = fVersion.WriteString(`#
#Sat Jun 13 21:48:54 CEST 2020
maltego.client.version=4.2.11.13104
maltego.client.subtitle=
maltego.pandora.version=1.4.2
maltego.client.name=Maltego Classic Eval
maltego.mtz.version=1.0
maltego.graph.version=1.2`)

	fmt.Println("generated maltego transform archive")
}

func PackTransformArchive() {
	fmt.Println("packing maltego transform archive")

	// zip and rename to: transforms.mtz
	f, err := os.Create("transforms" + configFileExtension)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		errClose := f.Close()
		if errClose != nil && !errors.Is(errClose, io.EOF) {
			fmt.Println(errClose)
		}
	}()

	w := zip.NewWriter(f)

	// add files to the archive
	addFiles(w, "transforms", "")

	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("packed maltego transform archive")
}

func PackMaltegoArchive(name string) {
	fmt.Println("packing maltego " + name + " archive")

	// zip and rename to: transforms.mtz
	f, err := os.Create(name + configFileExtension)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		errClose := f.Close()
		if errClose != nil && !errors.Is(errClose, io.EOF) {
			fmt.Println(errClose)
		}
	}()

	w := zip.NewWriter(f)

	// add files to the archive
	addFiles(w, name, "")

	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("packed maltego " + name + " archive")
}
