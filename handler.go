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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var transforms []string

// RegisterTransform will register the provided handler in the http.DefaultServeMux
// and collect the name for the route
func RegisterTransform(handlerFunc http.HandlerFunc, name string) {
	transforms = append(transforms, name)
	http.HandleFunc("/run/"+name, handlerFunc)
}

// Home provides a simple greeting together with a listing of supported transforms.
func Home(w http.ResponseWriter, r *http.Request) {

	fmt.Println("RemoteAddr", r.RemoteAddr, "UserAgent", r.UserAgent(), "URI", r.RequestURI)

	var routes string
	for _, t := range transforms {
		routes += "/run/" + t + "<br>"
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hi there! You've reached a Maltego transform server.<br><br>routes:<br>" + routes))
}

// MakeHandler is util to create a http.HandlerFunc, that will get the deserialized MaltegoMessage from a request,
// and can populate the Transform response, which will be written back into the connection as soon as the handler exits.
func MakeHandler(handler func(w http.ResponseWriter, r *http.Request, t *Transform)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("RemoteAddr", r.RemoteAddr, "UserAgent", r.UserAgent(), "URI", r.RequestURI)

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("please send a POST request to this endpoint"))
			return
		}

		// read request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("failed to read request body:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		fmt.Println(r.RemoteAddr, "body contains", len(body), "bytes of data")
		if len(body) == 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("empty body received. please add data"))
			return
		}

		// parse the transform from the request body bytes
		t := &Transform{}
		err = xml.Unmarshal(body, t)
		if err != nil {
			dump(body, request)
			fmt.Println("failed to unmarshal transform:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// request always has the first entity set
		if t.RequestMessage == nil || len(t.RequestMessage.Entities.Items) != 1 {
			dump(body, request)
			if t.RequestMessage == nil {
				fmt.Println("no RequestMessage provided")
			} else {
				fmt.Println("invalid number of entities:", len(t.RequestMessage.Entities.Items))
			}

			http.Error(w, "malformed RequestMessage", http.StatusBadRequest)
			return
		}

		dump(body, request)

		// invoke the user provided handler
		handler(w, r, t)

		if debug {
			formatted, err := xml.MarshalIndent(t, "", "  ")
			if err != nil {
				log.Println("failed to marshal transform: ", err)
			}
			dump(formatted, response)
		}

		t.AddUIMessage("complete", UIMessageInform)

		// write back the response
		_, err = fmt.Fprintf(w, t.ReturnOutput())
		if err != nil {
			fmt.Println("failed to write back response:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
