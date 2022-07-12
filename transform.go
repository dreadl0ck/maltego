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
	"log"
)

// Transform models a maltego transformation message.
type Transform struct {
	XMLName          xml.Name          `xml:"MaltegoMessage"`
	ResponseMessage  *ResponseMessage  `xml:"MaltegoTransformResponseMessage,omitempty"`
	ExceptionMessage *ExceptionMessage `xml:"MaltegoTransformExceptionMessage"`
	RequestMessage   *RequestMessage   `xml:"MaltegoTransformRequestMessage,omitempty"`
}

// ResponseMessage models a maltego response message.
type ResponseMessage struct {
	Entities   Entities   `xml:"Entities"`
	UIMessages UIMessages `xml:"UIMessages"`
}

// Entities is a container for maltego entities.
type Entities struct {
	Items []*Entity `xml:"Entity"`
}

// UIMessages is a container for maltego UIMessages.
type UIMessages struct {
	Items []*UIMessage `xml:"UIMessage"`
}

// UIMessage models a maltego UI message.
type UIMessage struct {
	Text        string `xml:",chardata"`
	MessageType string `xml:"MessageType,attr"`
}

// ExceptionMessage contains one or more exceptions.
type ExceptionMessage struct {
	Exceptions Exceptions `xml:"Exceptions"`
}

// Exceptions is a container for maltego exceptions.
type Exceptions struct {
	Items []*Exception `xml:"Exception"`
}

// Exception models a maltego exception.
type Exception struct {
	Text string `xml:",chardata"`
	Code string `xml:"code,attr"`
}

// AddEntity adds an entity to the transform.
func (tr *Transform) AddEntity(typ, value string) *Entity {

	// ensure response message is initialized
	if tr.ResponseMessage == nil {
		tr.ResponseMessage = &ResponseMessage{}
	}

	ent := NewEntity(typ, EscapeText(value), "100")
	tr.ResponseMessage.Entities.Items = append(tr.ResponseMessage.Entities.Items, ent)

	return ent
}

// AddUIMessage adds a UI message to the transform.
func (tr *Transform) AddUIMessage(message, messageType string) {

	// ensure response message is initialized
	if tr.ResponseMessage == nil {
		tr.ResponseMessage = &ResponseMessage{}
	}

	// add UIMessage
	tr.ResponseMessage.UIMessages.Items = append(tr.ResponseMessage.UIMessages.Items, &UIMessage{
		Text:        message,
		MessageType: messageType,
	})
}

// AddException adds an exception to the transform.
func (tr *Transform) AddException(exceptionString, code string) {

	// ensure response message is initialized
	if tr.ExceptionMessage == nil {
		tr.ExceptionMessage = &ExceptionMessage{}
	}

	// add exception
	tr.ExceptionMessage.Exceptions.Items = append(tr.ExceptionMessage.Exceptions.Items, &Exception{
		Text: exceptionString,
		Code: code,
	})
}

// DisplayInformation models maltego display information.
type DisplayInformation struct {
	Labels []*DisplayLabel `xml:"Label"`
}

// DisplayLabel models a label for display information.
type DisplayLabel struct {
	XMLName xml.Name `xml:"Label"`
	Text    string   `xml:",cdata"`
	Name    string   `xml:"Name,attr"`
	Type    string   `xml:"Type,attr"`
}

func NewDisplayLabel(text string, name string) *DisplayLabel {
	return &DisplayLabel{
		Text: text,
		Name: name,
		Type: "text/html",
	}
}

// ReturnOutput returns the transformations XML representation.
func (tr *Transform) ReturnOutput() string {

	data, err := xml.Marshal(tr)
	if err != nil {
		log.Println("failed to marshal transform: ", err)
	}

	return string(data)
}

// ThrowExceptions generates an exception message.
func (tr *Transform) ThrowExceptions() string {

	tr.ResponseMessage = nil

	data, err := xml.Marshal(tr)
	if err != nil {
		log.Println("failed to marshal transform: ", err)
	}

	return string(data)
}
