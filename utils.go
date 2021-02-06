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
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var newlineReplacer = strings.NewReplacer("&#xA;", "\n")

type messageType string

const debug = true

const (
	response messageType = "RESPONSE"
	request  messageType = "REQUEST"
)

func dump(data []byte, typ messageType) {
	if debug {
		fmt.Println("================== " + typ + " ====================")
		fmt.Println(string(data))
		fmt.Println("===============================================")
	}
}


// EscapeText ensures that the input text is safe to embed within XML.
func EscapeText(text string) string {
	var buf bytes.Buffer

	err := xml.EscapeText(&buf, []byte(text))
	if err != nil {
		fmt.Println(err)
	}

	return newlineReplacer.Replace(buf.String())
}

// Die will create a new transform with an error message and signal an error and the output to maltego.
func Die(err string, msg string) {
	trx := Transform{}
	// add error message for the user
	trx.AddUIMessage(msg+": "+err, UIMessageFatal)
	fmt.Println(trx.ReturnOutput())
	log.Println(msg, err)
	os.Exit(0) // don't signal an error for the transform invocation
}

// GetThickness can be used to calculate the line thickness.
func GetThickness(val, min, max uint64) int {
	if min == max {
		min = 0
	}

	delta := max - min

	// log.Println("delta=", delta, "float64(delta)*0.01 = ", float64(delta)*0.01)
	// log.Println("delta=", delta, "float64(delta)*0.1 = ", float64(delta)*0.1)
	// log.Println("delta=", delta, "float64(delta)*0.5 = ", float64(delta)*0.5)
	// log.Println("delta=", delta, "float64(delta)*1 = ", float64(delta)*1)
	// log.Println("delta=", delta, "float64(delta)*2 = ", float64(delta)*2)

	switch {
	case float64(val) <= float64(delta)*0.01:
		return 1
	case float64(val) <= float64(delta)*0.1:
		return 2
	case float64(val) <= float64(delta)*0.3:
		return 3
	case float64(val) <= float64(delta)*0.6:
		return 4
	case float64(val) <= float64(delta)*0.9:
		return 5
	default:
		return 5
	}
}

// PrintProgress sets the progressbar in Maltego
// this is documented in the old versions of the Maltego manual
// but does not seem to work with the current version
func (tr *Transform) PrintProgress(percentage int) {

	if percentage < 0 || percentage > 100 {
		fmt.Println("invalid percentage value:", percentage)
		return
	}

	_, err := os.Stderr.WriteString("%" + strconv.Itoa(percentage) + "\n")
	if err != nil {
		log.Println("failed to write progress update: ", err)
	}
}

// GetThicknessInterval returns a value for the line thickness.
// Calculation happens based on the values provided for min and max.
func GetThicknessInterval(val, min, max uint64) int {
	if min == max {
		min = 0
	}

	interval := (max - min) / 5

	switch {
	case val <= interval:
		return 1
	case val <= interval*2:
		return 2
	case val <= interval*3:
		return 3
	case val <= interval*4:
		return 4
	case val <= interval*5:
		return 5
	default: // bigger than interval*5
		return 5
	}
}

// noPluralsMap contains words for which to make an exception when pluralizing nouns.
var NoPluralsMap = map[string]struct{}{
	"Software": {},
	"Ethernet": {},
}

// Pluralize returns the plural for a given noun.
func Pluralize(name string) string {
	if strings.HasSuffix(name, "e") || strings.HasSuffix(name, "w") {
		if _, ok := NoPluralsMap[name]; !ok {
			name += "s"
		}
	}

	if strings.HasSuffix(name, "y") {
		name = name[:len(name)-1] + "ies"
	}

	if strings.HasSuffix(name, "t") {
		if _, ok := NoPluralsMap[name]; !ok {
			name += "s"
		}
	}

	if strings.HasSuffix(name, "n") {
		if _, ok := NoPluralsMap[name]; !ok {
			name += "s"
		}
	}

	return name
}

func GenServerListing(prefix, outDir string, trs []TransformCoreInfo) {
	srv := Server{
		Name:        "Local",
		Enabled:     true,
		Description: "Local transforms hosted on this machine",
		URL:         "http://localhost",
		LastSync:    time.Now().Format("2006-01-02 15:04:05.000 MST"), // example: 2020-06-23 20:47:24.433 CEST"
		Protocol: struct {
			Text    string `xml:",chardata"`
			Version string `xml:"version,attr"`
		}{
			Version: "0.0",
		},
		Authentication: struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		}{
			Type: "none",
		},
		Seeds: "",
	}

	for _, t := range trs {
		srv.Transforms.Transform = append(srv.Transforms.Transform, struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		}{
			Name: prefix + t.ID,
		})
	}

	data, err := xml.MarshalIndent(srv, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filepath.Join(outDir, "Servers", "Local.tas"))
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

func GenTransformSet(name string, description string, prefix string, outDir string, trs []TransformCoreInfo) {
	tSet := TransformSet{
		Name:        name,
		Description: description,
	}

	for _, t := range trs {
		tSet.Transforms.Transform = append(tSet.Transforms.Transform, struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		}{
			Name: prefix + t.ID,
		})
	}

	data, err := xml.MarshalIndent(tSet, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	_ = os.MkdirAll(filepath.Join(outDir, "TransformSets"), 0o700)
	f, err := os.Create(filepath.Join(outDir, "TransformSets", name + ".set"))
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

func GenFullConfigArchive(ident string) {
	// clean
	_ = os.RemoveAll(ident)

	// create directories
	_ = os.MkdirAll(filepath.Join(ident, "Servers"), 0o700)
	_ = os.MkdirAll(filepath.Join(ident, "TransformRepositories", "Local"), 0o700)

	// create directories
	_ = os.MkdirAll(filepath.Join(ident, "Entities"), 0o700)
	_ = os.MkdirAll(filepath.Join(ident, "EntityCategories"), 0o700)
	_ = os.MkdirAll(filepath.Join(ident, "Icons"), 0o700)

	fVersion, err := os.Create(filepath.Join(ident, "version.properties"))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if errClose := fVersion.Close(); errClose != nil {
			fmt.Println(errClose)
		}
	}()

	fCategory, err := os.Create(filepath.Join(ident, "EntityCategories", ident + ".category"))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if errClose := fCategory.Close(); errClose != nil {
			fmt.Println(errClose)
		}
	}()

	// Sat Jun 13 21:48:54 CEST 2020
	_, _ = fVersion.WriteString(`#
#` + time.Now().Format(time.UnixDate) + `
maltego.client.version=4.2.12
maltego.client.subtitle=
maltego.pandora.version=1.4.2
maltego.client.name=Maltego Classic Eval
maltego.mtz.version=1.0
maltego.graph.version=1.2`)

	_, _ = fCategory.WriteString("<EntityCategory name=\"" + strings.ToTitle(ident) + "\"/>")

	fmt.Println("bootstrapped configuration archive for Maltego")
}

func GenMachines(ident string, machinePrefix string) {
	path := filepath.Join(ident, "Machines")

	err := os.Mkdir(path, 0700)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir("machines")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {

		// Machine Properties
		propFile, errCompile := os.Create(
			filepath.Join(
				path,
				machinePrefix+strings.Replace(
					filepath.Base(f.Name()),
					".machine",
					".properties",
					1,
				),
			),
		)
		if errCompile != nil {
			log.Fatal(errCompile)
		}

		_, _ = propFile.WriteString(`#` + time.Now().Format(time.UnixDate) + `
favorite=true
enabled=true`)

		err = propFile.Close()
		if err != nil {
			log.Fatal(err)
		}

		// Machine

		copyFile(
			filepath.Join("machines", f.Name()),
			filepath.Join(
				path,
				machinePrefix+filepath.Base(f.Name()),
			),
		)
	}
}