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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func NewMaltegoEntity(category, identArchive, ident, prefix, propsPrefix, entName, imgName, description, parent string, isArchive bool, r *RegexConversion, propertyFields ...PropertyField) MaltegoEntity {
	if !strings.Contains(imgName, "/") {
		imgName = ident + "/" + imgName
	}

	var (
		name = prefix + entName
		ent  = MaltegoEntity{
			ID:                name,
			DisplayName:       entName,
			DisplayNamePlural: Pluralize(entName),
			Description:       description,
			SmallIconResource: imgName,
			LargeIconResource: imgName,
			AllowedRoot:       true,
			ConversionOrder:   "2147483647",
			Visible:           true,
			Properties: EntityProperties{
				Value:        propsPrefix + strings.ToLower(entName),
				DisplayValue: propsPrefix + strings.ToLower(entName),
				Fields: Fields{
					Items: []PropertyField{
						{
							Name:        propsPrefix + strings.ToLower(entName),
							Type:        "string",
							Nullable:    true,
							Hidden:      false,
							Readonly:    false,
							Description: "",
							DisplayName: entName,
							SampleValue: "-",
						},
					},
				},
			},
		}
	)

	if r != nil {
		// make sure the regex is valid
		_ = regexp.MustCompile(r.Regex)

		// set converter
		ent.Converter = &Converter{
			Value: r.Regex,
		}

		// add property mappings
		for _, p := range r.Properties {
			ent.Converter.Groups.RegexGroup = append(ent.Converter.Groups.RegexGroup, RegexGroup{
				Property: p,
			})
		}
	}

	if isArchive {
		ent.Category = identArchive
	} else {
		ent.Category = category
	}

	if len(propertyFields) > 0 {
		ent.Properties.Fields.Items = append(ent.Properties.Fields.Items, propertyFields...)
	}

	if len(parent) > 0 {
		ent.Entities = &BaseEntities{
			Entities: []BaseEntity{
				{
					Text: parent,
				},
			},
		}
	}

	return ent
}

func NewStringField(name string, description string) PropertyField {
	return PropertyField{
		Name:        strings.ToLower(name),
		Type:        "string",
		Nullable:    true,
		Hidden:      false,
		Readonly:    false,
		Description: description,
		DisplayName: strings.Title(name),
		SampleValue: "",
	}
}

func NewRequiredStringField(name string, description string) PropertyField {
	return PropertyField{
		Name:        strings.ToLower(name),
		Type:        "string",
		Nullable:    false,
		Hidden:      false,
		Readonly:    false,
		Description: description,
		DisplayName: strings.Title(name),
		SampleValue: "",
	}
}

// TODO: add config struct with default
func GenEntity(category, identArchive, ident, prefix, propsPrefix, outDir string, entName string, imgName string, description string, parent string, isArchive bool, color string, regex *RegexConversion, fields ...PropertyField) {
	imgName = imgName + "_" + color

	var (
		name = prefix + entName
		ent  = NewMaltegoEntity(category, identArchive, ident, prefix, propsPrefix, entName, imgName, description, parent, isArchive, regex, fields...)
	)

	data, err := xml.MarshalIndent(ent, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(filepath.Join(outDir, "Entities", name+".entity"))
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

	// add icon files
	_ = os.MkdirAll(filepath.Join(outDir, "Icons", ident), 0o700)

	var (
		ext  = ".svg"
		dir  = "material-icons"
		base = filepath.Join("/tmp", "icons", dir, "renamed", imgName)
	)

	if _, err = os.Stat(base + "16" + ext); err != nil {
		ext = ".png"
		dir = "material-icons-png"
		base = filepath.Join("/tmp", "icons", dir, "renamed", imgName)
	}

	dstBase := filepath.Join(outDir, "Icons", ident, imgName)

	copyFile(
		filepath.Join("/tmp", "icons", dir, "renamed", imgName+".xml"),
		filepath.Join(outDir, "Icons", ident, imgName+".xml"),
	)

	copyFile(base+"16"+ext, dstBase+ext)
	copyFile(base+"24"+ext, dstBase+"24"+ext)
	copyFile(base+"32"+ext, dstBase+"32"+ext)
	copyFile(base+"48"+ext, dstBase+"48"+ext)
	copyFile(base+"96"+ext, dstBase+"96"+ext)
}

// copyFile the source file contents to destination
// file attributes wont be copied and an existing file will be overwritten.
func copyFile(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if errClose := in.Close(); errClose != nil {
			fmt.Println(errClose)
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(out, in)
	if err != nil {
		log.Fatal(err)
	}

	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// Directory structure:
// .
// ├── entities
// │     ├── ...
// │     └── example.entity
// ├── EntityCategories
// │     └── EntityCategory.category
// ├── Icons
// │     └── EntityCategory
// │           ├── sim_card_alert.png
// │           ├── sim_card_alert.xml
// │           ├── sim_card_alert24.png
// │           ├── sim_card_alert32.png
// │           ├── sim_card_alert48.png
// │           └── sim_card_alert96.png
// └── version.properties.
func GenEntityArchive(entityCategory string) {
	// clean
	_ = os.RemoveAll("entities")

	// create directories
	_ = os.MkdirAll("entities/Entities", 0o700)
	_ = os.MkdirAll("entities/EntityCategories", 0o700)
	_ = os.MkdirAll("entities/Icons", 0o700)

	fVersion, err := os.Create("entities/version.properties")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if errClose := fVersion.Close(); errClose != nil {
			fmt.Println(errClose)
		}
	}()

	fCategory, err := os.Create("entities/EntityCategories/" + entityCategory + ".category")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if errClose := fCategory.Close(); errClose != nil {
			fmt.Println(errClose)
		}
	}()

	_, _ = fVersion.WriteString(`#
#` + time.Now().Format(time.UnixDate) + `
client.version=4.2.12
client.subtitle=
pandora.version=1.4.2
client.name=Maltego Classic Eval
mtz.version=1.0
graph.version=1.2`)

	_, _ = fCategory.WriteString("<EntityCategory name=\"" + entityCategory + "\"/>")

	fmt.Println("generated maltego entity archive")
}

func PackEntityArchive() {
	fmt.Println("packing maltego entity archive")

	// zip and rename to: entities.mtz
	f, err := os.Create("entities" + configFileExtension)
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
	addFiles(w, "entities", "")

	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("packed maltego entity archive")
}

func addFiles(wr *zip.Writer, basePath, baseInZip string) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		path := filepath.Join(basePath, file.Name())
		//fmt.Println(path)

		if !file.IsDir() {
			data, errRead := ioutil.ReadFile(path)
			if errRead != nil {
				fmt.Println(errRead)
			}

			// add files to the archive
			f, errCreate := wr.Create(filepath.Join(baseInZip, file.Name()))
			if errCreate != nil {
				log.Fatal(errCreate)
			}

			_, err = f.Write(data)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			newBase := filepath.Join(basePath, file.Name(), "/")
			//fmt.Println("adding sub directory: " + newBase)
			addFiles(wr, newBase, filepath.Join(baseInZip, file.Name(), "/"))
		}
	}
}
