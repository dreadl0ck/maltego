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

import "encoding/xml"

const (
	AS                    = "maltego.AS"
	Affiliation           = "maltego.Affiliation"
	Alias                 = "maltego.Alias"
	Banner                = "maltego.Banner"
	BuiltWithRelationship = "maltego.BuiltWithRelationship"
	BuiltWithTechnology   = "maltego.BuiltWithTechnology"
	CircularArea          = "maltego.CircularArea"
	Company               = "maltego.Company"
	DNSName               = "maltego.DNSName"
	DateTime              = "maltego.DateTime"
	Device                = "maltego.Device"
	Document              = "maltego.Document"
	Domain                = "maltego.Domain"
	EmailAddress          = "maltego.EmailAddress"
	File                  = "maltego.File"
	GPS                   = "maltego.GPS"
	Hash                  = "maltego.Hash"
	IPv4Address           = "maltego.IPv4Address"
	Image                 = "maltego.Image"
	Location              = "maltego.Location"
	MXRecord              = "maltego.MXRecord"
	NSRecord              = "maltego.NSRecord"
	Netblock              = "maltego.Netblock"
	Organization          = "maltego.Organization"
	Person                = "maltego.Person"
	PhoneNumber           = "maltego.PhoneNumber"
	Phrase                = "maltego.Phrase"
	Port                  = "maltego.Port"
	Sentiment             = "maltego.Sentiment"
	Service               = "maltego.Service"
	Twit                  = "maltego.Twit"
	URL                   = "maltego.URL"
	UniqueIdentifier      = "maltego.UniqueIdentifier"
	WebTitle              = "maltego.WebTitle"
	Website               = "maltego.Website"
)

// MaltegoEntity represents an exported entity model on disk
type MaltegoEntity struct {
	XMLName xml.Name `xml:"MaltegoEntity"`
	ID      string   `xml:"id,attr"`

	DisplayName       string `xml:"displayName,attr"`
	DisplayNamePlural string `xml:"displayNamePlural,attr"`
	Description       string `xml:"description,attr"`
	Category          string `xml:"category,attr"`

	SmallIconResource string `xml:"smallIconResource,attr"`
	LargeIconResource string `xml:"largeIconResource,attr"`

	AllowedRoot     bool   `xml:"allowedRoot,attr"`
	ConversionOrder string `xml:"conversionOrder,attr"`
	Visible         bool   `xml:"visible,attr"`

	Entities   *BaseEntities    `xml:"BaseEntities,omitempty"`
	Properties EntityProperties `xml:"Properties"`

	Converter *Converter `xml:"Converter,omitempty"`
}

// Converter contains information how to detect entities based on a regular expression.
type Converter struct {
	XMLName xml.Name    `xml:"Converter"`
	Text    string      `xml:",chardata"`
	Value   string      `xml:"Value"`
	Groups  RegexGroups `xml:"RegexGroups"`
}

// RegexGroups is a container for regex groups.
type RegexGroups struct {
	Text       string       `xml:",chardata"`
	RegexGroup []RegexGroup `xml:"RegexGroup"`
}

// RegexGroup structure
type RegexGroup struct {
	Text     string `xml:",chardata"`
	Property string `xml:"property,attr"`
}

// BaseEntities structure
type BaseEntities struct {
	Text     string `xml:",chardata"`
	Entities []BaseEntity
}

// BaseEntity structure
type BaseEntity struct {
	Text string `xml:",chardata"`
}

// EntityProperties contain property metadata
type EntityProperties struct {
	XMLName      xml.Name `xml:"Properties"`
	Value        string   `xml:"value,attr"`
	DisplayValue string   `xml:"displayValue,attr"`
	Groups       string   `xml:"Groups"`
	Fields       Fields   `xml:"Fields"`
}

// Fields hold property items.
type Fields struct {
	Items []PropertyField
}

// PropertyField are set on entities.
type PropertyField struct {
	XMLName     xml.Name `xml:"Field"`
	Text        string   `xml:",chardata"`
	Name        string   `xml:"name,attr"`
	Type        string   `xml:"type,attr"`
	Nullable    bool     `xml:"nullable,attr"`
	Hidden      bool     `xml:"hidden,attr"`
	Readonly    bool     `xml:"readonly,attr"`
	Description string   `xml:"description,attr"`
	DisplayName string   `xml:"displayName,attr"`
	SampleValue string   `xml:"SampleValue"`
}

// EntityCoreInfo describes an entity.
type EntityCoreInfo struct {
	Name        string
	Icon        string
	Description string
	Parent      string
	Fields      []PropertyField
}

// RegexConversion contains conversion information.
type RegexConversion struct {
	Regex      string
	Properties []string
}
