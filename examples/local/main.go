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

package main

import (
	"fmt"
	"github.com/dreadl0ck/maltego"
	"net"
	"os"
)

// This is an example for a local transformation that does a reverse name lookup for a given address.
// It will take an IP address and return the hostnames associated with it as maltego entities.
func main() {

	// parse arguments
	lt := maltego.ParseLocalArguments(os.Args)

	// ensure the provided address is valid
	ip := net.ParseIP(lt.Value)
	if ip == nil {
		maltego.Die("invalid ip", lt.Value+" is not a valid IP address")
	}

	// lookup provided ip address
	names, err := net.LookupAddr(lt.Value)
	if err != nil {
		maltego.Die(err.Error(), "failed to lookup address")
	}

	// create new transform
	t := maltego.Transform{}

	// iterate over lookup results
	for _, host := range names {
		e := t.AddEntity("maltego.DNSName", host)
		e.AddProperty("hostname", "Hostname", maltego.Strict, host)
	}

	t.AddUIMessage("complete", maltego.UIMessageInform)

	// return output to stdout and exit cleanly (exit code 0)
	fmt.Println(t.ReturnOutput())
}
