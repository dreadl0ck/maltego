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
	"flag"
	"fmt"
	"github.com/dreadl0ck/maltego"
	"log"
	"net/http"
)

var (
	flagAddr   = flag.String("addr", ":8081", "server listen address")
)

func main() {

	flag.Parse()

	// register transforms to http.DefaultServeMux
	maltego.RegisterTransform(lookupAddr, "lookupAddr")
	maltego.RegisterTransform(lookupMX, "lookupMX")
	maltego.RegisterTransform(lookupNS, "lookupNS")
	maltego.RegisterTransform(lookupIP, "lookupIP")
	maltego.RegisterTransform(lookupTXT, "lookupTXT")
	maltego.RegisterTransform(lookupPort, "lookupPort")
	maltego.RegisterTransform(lookupCNAME, "lookupCNAME")
	maltego.RegisterTransform(lookupSRV, "lookupSRV")

	// register catch all handler to serve home page
	http.HandleFunc("/", maltego.Home)

	fmt.Println("serving at", *flagAddr)

	// start server
	err := http.ListenAndServe(*flagAddr, nil) // no handler passed => http.DefaultServeMux will be used!
	if err != nil {
		log.Fatal("failed to serve HTTP: ", err)
	}
}
