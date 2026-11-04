package cli

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/qlyco/qwc/core"
	"github.com/qlyco/qwc/internal/console"
	"github.com/qlyco/qwc/internal/res"
)

type App struct {
	Name    string
	Version string
	Setup   func(site *core.Server) error
	Handler func(site *core.Server) http.Handler
}

func (app App) banner() string {
	version := app.Version
	info, isAvailable := debug.ReadBuildInfo()

	if !isAvailable {
		return fmt.Sprintf(res.CustomBannerMessage, app.Name, version)
	}

	if (version == "" || version == res.VersionDefault) && isStamped(info.Main) {
		version = info.Main.Version
	}

	if info.Main.Path == res.ModulePath {
		return fmt.Sprintf(res.BannerMessage, res.AppName, res.AppDescription, version)
	}

	for _, dependency := range info.Deps {
		if dependency.Path == res.ModulePath {
			return fmt.Sprintf(res.BuiltOnBannerMessage, app.Name, version, res.AppName, dependency.Version)
		}
	}

	return fmt.Sprintf(res.CustomBannerMessage, app.Name, version)
}

func (app App) handler(server *core.Server) http.Handler {
	if app.Handler == nil {
		return server
	}

	return app.Handler(server)
}

func isStamped(module debug.Module) bool {
	return module.Version != "" && module.Version != res.UnstampedVersion
}

func Run(app App) {
	root := flag.String(res.RootFlag, res.RootDefault, res.RootUsage)
	port := flag.Int(res.PortFlag, res.PortDefault, res.PortUsage)
	build := flag.Bool(res.BuildFlag, false, res.BuildUsage)
	output := flag.String(res.OutputFlag, res.OutputDefault, res.OutputUsage)
	siteURL := flag.String(res.URLFlag, "", res.URLUsage)
	certFile := flag.String(res.CertFlag, "", res.CertUsage)
	keyFile := flag.String(res.KeyFlag, "", res.KeyUsage)
	configPath := flag.String(res.ConfigFlag, "", res.ConfigUsage)
	verbose := flag.Bool(res.VerboseFlag, false, res.VerboseUsage)

	flag.Parse()

	if flag.NArg() > 0 {
		log.Printf(res.UnexpectedArgMessage, flag.Arg(0))
		flag.Usage()
		os.Exit(res.UsageExitCode)
	}

	if *configPath != "" {
		err := applyConfig(*configPath)

		if err != nil {
			log.Fatal(err)
		}
	}

	if (*certFile == "") != (*keyFile == "") {
		log.Fatal(res.CertPairMessage)
	}

	fmt.Print(app.banner())

	err := prepareSiteFolder(*root)

	if err != nil {
		log.Fatal(err)
	}

	server := core.New(*root, log.Default())

	if *verbose {
		server.SetVerbose(true)
	}

	if app.Setup != nil {
		err = app.Setup(server)

		if err != nil {
			log.Fatal(err)
		}
	}

	if *build {
		siteDir, err := server.Build(*output, *siteURL)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(res.BuiltMessage, siteDir)
		return
	}

	err = server.SetURL(*siteURL)

	if err != nil {
		log.Fatal(err)
	}

	address := fmt.Sprintf(res.ListenAddressFormat, *port)
	web := &http.Server{Addr: address, Handler: app.handler(server)}
	scheme := res.HTTPScheme

	if *certFile != "" {
		certificate, err := loadCertificate(*certFile, *keyFile)

		if err != nil {
			log.Fatal(err)
		}

		web.TLSConfig = &tls.Config{GetCertificate: certificate.get}
		scheme = res.HTTPSScheme
	}

	fmt.Printf(res.ServingMessage, *root, scheme, address)
	fmt.Print(res.KeyHintMessage)

	terminal := console.OpenTerminal()

	console.Listen(server, terminal)

	if *certFile != "" {
		err = web.ListenAndServeTLS("", "")
	} else {
		err = web.ListenAndServe()
	}

	terminal.Restore()
	log.Fatal(err)
}

func prepareSiteFolder(root string) error {
	if root == res.RootDefault {
		err := os.MkdirAll(root, os.ModePerm)

		if err != nil {
			return err
		}
	}

	info, err := os.Stat(root)

	if err != nil || !info.IsDir() {
		return fmt.Errorf(res.MissingSiteDirMessage, root)
	}

	return nil
}
