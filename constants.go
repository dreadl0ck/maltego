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

// constants
const (
	BookMarkColorNone   = "-1"
	BookMarkColorBlue   = "0"
	BookMarkColorGreen  = "1"
	BookMarkColorYellow = "2"
	BookMarkColorOrange = "3"
	BookMarkColorRed    = "4"

	LinkStyleNormal  = "0"
	LinkStyleDashed  = "1"
	LinkStyleDotted  = "2"
	LinkStyleDashdot = "3"

	UIMessageFatal        = "FatalError"
	UIMessagePartialError = "PartialError"
	UIMessageInform       = "Inform"
	UIMessageDebug        = "Debug"

	// Strict is used for enabling strict property matching
	Strict = "strict"

	// Loose enables loose property matching
	Loose = "loose"
)

// LinkDirection determines the direction of node interconnections (links).
type LinkDirection string

const (
	// OutputToInput direction for maltego links
	OutputToInput LinkDirection = "output-to-input"

	// InputToOutput direction for maltego links
	InputToOutput LinkDirection = "input-to-output"

	// bidirectional direction for maltego links
	bidirectional LinkDirection = "bidirectional"
)

// properties
const (
	LinkColor             = "link#maltego.link.color"
	LinkStyle             = "link#maltego.link.style"
	LinkThickness         = "link#maltego.link.thickness"
	Label                 = "link#maltego.link.label"
	PropertyLinkDirection = "link#maltego.link.direction"
	Bookmark              = "bookmark#"
	Notes                 = "notes#"
)
