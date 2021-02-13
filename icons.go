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
	"log"
	"os"
)

var icon = `<Icon>
<Aliases/>
</Icon>`

// CreateXMLIconFile will create the XML structure at the given path.
func CreateXMLIconFile(path string) {
	// create XML info file for maltego
	fXML, err := os.Create(path + ".xml")
	if err != nil {
		log.Fatal(err)
	}

	_, err = fXML.WriteString(icon)
	if err != nil {
		log.Fatal(err)
	}

	err = fXML.Close()
	if err != nil {
		log.Fatal(err)
	}
}

