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

// RequestMessage models a request.
type RequestMessage struct {
	XMLName         xml.Name        `xml:"MaltegoTransformRequestMessage"`
	Entities        Entities        `xml:"Entities"`
	Limits          Limits          `xml:"Limits"`
	TransformFields TransformFields `xml:"TransformFields"`
}

// Limits structure.
type Limits struct {
	XMLName   xml.Name `xml:"Limits"`
	HardLimit string   `xml:"HardLimit,attr"`
	SoftLimit string   `xml:"SoftLimit,attr"`
}

type TransformFields struct {
	Fields []*TransformField `xml:"Field"`
}

// TransformField structure.
type TransformField struct {
	Text string `xml:",chardata"`
	Name string `xml:"Name,attr"`
}
