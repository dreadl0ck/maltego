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
	"net/http"
)

var lookupIP = maltego.MakeHandler(func(w http.ResponseWriter, r *http.Request, t *maltego.Transform) {

	// get host name from query
	host := t.RequestMessage.Entities.Items[0].Value

	fmt.Println("got request from", r.RemoteAddr, "to lookup IPs for:", host)

	// perform lookup
	ips, err := net.LookupIP(host)
	if err != nil {
		fmt.Println("failed to lookup IPs:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("results for", host, "=", ips)

	// process results
	for _, ip := range ips {
		t.AddEntity(maltego.IPv4Address, ip.String())
	}
})
