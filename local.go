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
	"strings"
)

// handles some characters that would cause errors
var replacer = strings.NewReplacer("&amp;", "&", "\\=", "=")

// LocalTransform is used to handle a local transform from stdin.
type LocalTransform struct {
	Value  string
	Values map[string]string
}

// ParseLocalArguments parses the arguments supplied on the commandline.
func ParseLocalArguments(args []string) LocalTransform {
	if len(args) < 2 {
		log.Fatal("need at least 2 arguments, got ", len(args), ": ", args)
	}

	var (
		value  = args[0]
		values = make(map[string]string)
	)

	if len(args) > 1 {
		// search the remaining arguments for variables
		for _, arg := range args[1:] {

			// remove any newlines
			arg = strings.ReplaceAll(arg, "\n", " ")

			if len(arg) > 0 {
				vars := strings.Split(arg, "#")
				for _, x := range vars {
					kv := strings.Split(x, "=")
					if len(kv) == 2 {
						values[kv[0]] = replacer.Replace(kv[1])
					} else {
						values[kv[0]] = replacer.Replace(strings.Join(kv[1:], "="))
					}
				}
			}
		}
	}

	return LocalTransform{
		Value:  value,
		Values: values,
	}
}
