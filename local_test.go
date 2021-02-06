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
	"fmt"
	"testing"
)


func TestParseLocalArguments(t *testing.T) {
	args := []string{"/var/folders/test", "pãypal.com\nxn--pypal-9qa.com\nregistered", "fqdn=pãypal.com\nxn--pypal-9qa.com\nregistered#unicode=pãypal.com#ascii=xn--pypal-9qa.com#status=registered#ips=34.102.136.180#names=180.136.102.34.bc.googleusercontent.com."}
	lt := ParseLocalArguments(args[1:])
	fmt.Println(lt.Values)
}
