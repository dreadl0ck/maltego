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
	"encoding/xml"
	"strconv"
	"strings"
)

/*
 *	Entity
 */

// Entity models a transform entity.
type Entity struct {
	XMLName   xml.Name            `xml:"Entity"`
	Type      string              `xml:"Type,attr"`
	Genealogy *Genealogy          `xml:"Genealogy,omitempty"`
	Value     string              `xml:"Value"`
	Weight    string              `xml:"Weight"`
	Info      *DisplayInformation `xml:"DisplayInformation,omitempty"`
	IconURL   string              `xml:"IconURL,omitempty"`
	Fields    *AdditionalFields   `xml:"AdditionalFields,omitempty"`
}

// AdditionalFields is a container for fields.
type AdditionalFields struct {
	XMLName xml.Name `xml:"AdditionalFields"`
	Items   []*Field  `xml:"Field"`
}

// Genealogy structure.
type Genealogy struct {
	Type GenealogyType `xml:"Type"`
}

// GenealogyType structure.
type GenealogyType struct {
	Name    string `xml:"Name,attr"`
	OldName string `xml:"OldName,attr"`
}

// Field structure.
type Field struct {
	Text         string `xml:",chardata"`
	MatchingRule string `xml:"MatchingRule,attr"`
	Name         string `xml:"Name,attr"`
	DisplayName  string `xml:"DisplayName,attr"`
}

// NewEntity is the constructor for an Entity.
func NewEntity(typ, value string, weight string) *Entity {
	return &Entity{
		Type:   typ,
		Value:  value,
		Weight: weight,
	}
}

func (tre *Entity) GetFieldByName(name string) string {
	for _, f := range tre.Fields.Items {
		if f.Name == name {
			return f.Text
		}
	}
	return ""
}

// AddProperty adds a property.
func (tre *Entity) AddProperty(fieldName, displayName, matchingRule, value string) {

	if tre.Fields == nil {
		tre.Fields = &AdditionalFields{}
	}

	// collect field
	tre.Fields.Items = append(tre.Fields.Items, &Field{
		Text:         EscapeText(value),
		MatchingRule: matchingRule,
		Name:         fieldName,
		DisplayName:  displayName,
	})
}

// AddProp is shorthand for a strict AddProperty, that uses the title version of the fieldName as displayName.
func (tre *Entity) AddProp(fieldName, value string) {

	if tre.Fields == nil {
		tre.Fields = &AdditionalFields{}
	}

	// collect field
	tre.Fields.Items = append(tre.Fields.Items, &Field{
		Text:         EscapeText(value),
		MatchingRule: Strict,
		Name:         fieldName,
		DisplayName:  strings.Title(fieldName),
	})
}

// AddDisplayInformation adds display information.
func (tre *Entity) AddDisplayInformation(text, name string) {
	if tre.Info == nil {
		tre.Info = &DisplayInformation{}
	}
	tre.Info.Labels = append(tre.Info.Labels, NewDisplayLabel(text, name))
}

// SetLinkColor sets the link color.
func (tre *Entity) SetLinkColor(color string) {
	tre.AddProperty(LinkColor, "LinkColor", Loose, color)
}

// SetLinkStyle sets the link style.
func (tre *Entity) SetLinkStyle(style string) {
	tre.AddProperty(LinkStyle, "LinkStyle", Loose, style)
}

// SetLinkThickness sets the link thickness.
func (tre *Entity) SetLinkThickness(thick int) {
	thickInt := strconv.Itoa(thick)
	tre.AddProperty(LinkThickness, "LinkThickness", Loose, thickInt)
}

// SetLinkLabel sets the link label.
func (tre *Entity) SetLinkLabel(label string) {
	tre.AddProperty(Label, "Label", Loose, label)
}

// SetBookmark sets a bookmark on the entity.
func (tre *Entity) SetBookmark(bookmark string) {
	tre.AddProperty(Bookmark, "Bookmark", Loose, bookmark)
}

// SetNote sets a note on the entity.
func (tre *Entity) SetNote(note string) {
	tre.AddProperty(Notes, "Notes", Loose, note)
}

// SetLinkDirection sets the link direction
func (tre *Entity) SetLinkDirection(dir LinkDirection) {
	tre.AddProperty(PropertyLinkDirection, "Direction", Loose, string(dir))
}
