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
	"context"
	"flag"
	"fmt"
	"github.com/dreadl0ck/maltego"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/foomo/simplecert"
	"github.com/foomo/tlsconfig"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

var (
	flagAddr = flag.String("addr", ":8081", "server listen address")
	flagTLS  = flag.String("tls", "", "use tls")
)

func initTrx() {

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
}

// This example demonstrates how spin up a custom HTTPS webserver for production deployment.
// It shows how to configure and start your service in a way that the certificate can be automatically renewed via the TLS challenge, before it expires.
// For this to succeed, we need to temporarily free port 443 (on which your service is running) and complete the challenge.
// Once the challenge has been completed the service will be restarted via the DidRenewCertificate hook.
// Requests to port 80 will always be redirected to the TLS secured version of your site.
func main() {

	// OAuth manager
	manager := manage.NewDefaultManager()

	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)

	s := server.NewDefaultServer(manager)
	s.SetAllowGetAccessRequest(true)
	s.SetClientInfoHandler(server.ClientFormHandler)

	s.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	s.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := s.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		s.HandleTokenRequest(w, r)
	})

	var (
		// the structure that handles reloading the certificate
		certReloader *simplecert.CertReloader
		err          error
		numRenews    int
		ctx, cancel  = context.WithCancel(context.Background())

		// init strict tlsConfig (this will enforce the use of modern TLS configurations)
		// you could use a less strict configuration if you have a customer facing web application that has visitors with old browsers
		tlsConf = tlsconfig.NewServerTLSConfig(tlsconfig.TLSModeServerStrict)

		// a simple constructor for a http.Server with our Handler
		makeServer = func() *http.Server {
			return &http.Server{
				Addr:      *flagAddr,
				Handler:   http.DefaultServeMux,
				TLSConfig: tlsConf,

				// prevent timeout on long running requests
				ReadTimeout:  0,
				WriteTimeout: 0,
				IdleTimeout:  0,
			}
		}

		// init server
		srv = makeServer()

		// init simplecert configuration
		cfg = simplecert.Default
	)

	initTrx()

	// check if a domain was provided, otherwise run without TLS
	if *flagTLS == "" {
		s := makeServer()
		log.Fatal(s.ListenAndServe())
	}

	// configure
	cfg.Domains = []string{*flagTLS}
	cfg.CacheDir = "letsencrypt"
	cfg.SSLEmail = "you@emailprovider.com"

	// disable HTTP challenges - we will only use the TLS challenge for this example.
	cfg.HTTPAddress = ""

	// this function will be called just before certificate renewal starts and is used to gracefully stop the service
	// (we need to temporarily free port 443 in order to complete the TLS challenge)
	cfg.WillRenewCertificate = func() {
		// stop server
		cancel()
	}

	// this function will be called after the certificate has been renewed, and is used to restart your service.
	cfg.DidRenewCertificate = func() {

		numRenews++

		// restart server: both context and server instance need to be recreated!
		ctx, cancel = context.WithCancel(context.Background())
		srv = makeServer()

		// force reload the updated cert from disk
		certReloader.ReloadNow()

		// here we go again
		go serve(ctx, srv)
	}

	log.Println("hello world")

	// init simplecert configuration
	// this will block initially until the certificate has been obtained for the first time.
	// on subsequent runs, simplecert will load the certificate from the cache directory on disk.
	certReloader, err = simplecert.Init(cfg, func() {
		os.Exit(0)
	})
	if err != nil {
		log.Fatal("simplecert init failed: ", err)
	}

	// redirect HTTP to HTTPS
	log.Println("starting HTTP Listener on Port 80")
	go http.ListenAndServe(":80", http.HandlerFunc(simplecert.Redirect))

	// enable hot reload
	tlsConf.GetCertificate = certReloader.GetCertificateFunc()

	// start serving
	log.Println("will serve at: https://" + cfg.Domains[0])
	serve(ctx, srv)

	fmt.Println("waiting forever")
	<-make(chan bool)
}

func serve(ctx context.Context, srv *http.Server) {

	// lets go
	go func() {
		if err := srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %+s\n", err)
		}
	}()

	log.Printf("server started")
	<-ctx.Done()
	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	err := srv.Shutdown(ctxShutDown)
	if err == http.ErrServerClosed {
		log.Printf("server exited properly")
	} else if err != nil {
		log.Printf("server encountered an error on exit: %+s\n", err)
	}
}
