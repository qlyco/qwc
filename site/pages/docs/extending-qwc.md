<!--
section: docs
title: Extending qwc
description: Add plugins to qwc, make your own qwc command, or run qwc inside your own Go program.
layout: doc
sort_key: 7
tags: setup
-->
# Extending qwc

The `qwc` command is a small Go program built on the qwc package. There are two ways to extend it:

<div class="table-scroll" role="region" aria-label="Ways to extend qwc" tabindex="0">

| Way | Good for | You write |
|---|---|---|
| [Your own qwc command](#your-own-command) | qwc as it is, plus plugins or a few flags | A small `main` |
| [qwc inside your program](#inside-your-program) | A bigger program that also serves a site, with its own routes and messages | The program around it |

</div>

Both start with the same thing:

<pre tabindex="0"><copy-code></copy-code><code>go get github.com/qlyco/qwc</code></pre>

## Plugins {#plugins}

A plugin gives templates data from Go, like products from a database. Any Go value works, usually a struct with methods. This one uses `database/sql`:

<pre tabindex="0"><copy-code></copy-code><code>{{`type Product struct {
	Name  string
	Price int
}

type store struct {
	// insert your data source here
}

func (shop store) Products() ([]Product, error) {
	// return all products
}

func (shop store) Product(id int) (Product, error) {
	// return the product with this id
}`}}</code></pre>

Register it under a name of letters and digits:

<pre tabindex="0"><copy-code></copy-code><code>{{`err := site.RegisterPlugin("db", store{})`}}</code></pre>

Templates get it with `plugin`. Each method can be called from the template:

<pre tabindex="0"><copy-code></copy-code><code>{{`{{$db := plugin "db"}}
<ul>{{range $db.Products}}<li>{{.Name}}: {{.Price}}</li>{{end}}</ul>
<p>Featured: {{($db.Product 5).Name}}</p>`}}</code></pre>

An unknown name shows <code>[qwc: error]</code> there. Variables only last within one template, so an include calls `plugin` itself, or gets the value with `dict`.

Share plugins as normal Go packages.

### Rules for plugins

- **Pages are cached.** A plugin runs when a page builds, and everyone sees that result until `site.Reload()`. Use plugins for data that is the same for everyone.
- **Pages build in parallel.** Methods must be safe to call at the same time. A `*sql.DB` is.
- **Errors show on the page** with `-verbose`. Return short messages, and log the details yourself.

## Your own qwc command {#your-own-command}

The `cli` package is the whole qwc command. Give it your name, version and a setup function:

<pre tabindex="0"><copy-code></copy-code><code>{{`package main

import (
	"github.com/qlyco/qwc/cli"
	"github.com/qlyco/qwc/core"
)

var Version = "development"

func main() {
	cli.Run(cli.App{
		Name:    "mysite",
		Version: Version,
		Setup: func(site *core.Server) error {
			return site.RegisterPlugin("db", store{})
		},
	})
}`}}</code></pre>

- `Setup` runs before qwc serves or builds, so plugins work in both.
- `Handler` is optional. Use it to add routes or middleware. It isn't used by `-build`.
- Flags you define with Go's `flag` package before calling `cli.Run` work like qwc's own, including in config files.
- The banner shows your name and version, and the qwc version it is built on, like `mysite [v1.2.0], built on qwc v0.0.1`. Set your version with `-ldflags "-X main.Version=v1.2.0"`.

For example, an API route and a header on every response:

<pre tabindex="0"><copy-code></copy-code><code>{{`cli.Run(cli.App{
	Name:    "shop",
	Version: Version,
	Handler: func(site *core.Server) http.Handler {
		mux := http.NewServeMux()
		mux.Handle("POST /api/products", saveProduct(site))
		mux.Handle("/", site)

		return withSecurityHeaders(mux)
	},
})`}}</code></pre>

`saveProduct` and `withSecurityHeaders` are your own code. Mount the site at `/`, and don't use `http.StripPrefix`.

**Example:** a shop whose product pages come from a database. Add a `db` plugin and a save route. Everything else stays the same: `-build`, `-config` and <kbd>r</kbd>.

## qwc inside your program {#inside-your-program}

A site is an `http.Handler`, so it works in any Go server, next to your own routes:

<pre tabindex="0"><copy-code></copy-code><code>{{`site := core.New("public", nil)

mux := http.NewServeMux()
mux.Handle("/api/", api)
mux.Handle("/", site)

http.ListenAndServe(":8080", mux)`}}</code></pre>

<div class="table-scroll" role="region" aria-label="Library methods" tabindex="0">

| Method | Same as |
|---|---|
| `core.New(root, logger)` | Opening a site folder. Pass `nil` to keep qwc quiet |
| `site.Reload()` | Pressing <kbd>r</kbd> |
| `site.Refresh()` | Pressing <kbd>R</kbd> |
| `site.SetURL("example.com")` | `-url=example.com` while serving |
| `site.SetVerbose(true)` | `-verbose` |
| `site.RegisterPlugin("db", value)` | Adds a [plugin](#plugins) |
| `site.Build("deploy", "example.com")` | `-build -output=deploy -url=example.com` |

</div>

`Build` returns the path of the `_qwcsite` folder it wrote. If the site has a path, like `example.com/docs`, pass the full path through, and don't use `http.StripPrefix`.

**Example:** an app with logins and an API, plus a blog and docs. Your program handles routes and logins. qwc serves the pages.

## Changing data

After a plugin's data changes, call `site.Reload()`. Pages rebuild on their next request:

<pre tabindex="0"><copy-code></copy-code><code>{{`mux.HandleFunc("POST /api/products", func(writer http.ResponseWriter, request *http.Request) {
	// check the login, then save the product
	site.Reload()
})`}}</code></pre>

In your own command, add the route in `Handler`. A separate program can delete `_qwcbuild` in the site folder instead, which also reloads.

## Data for each visitor

Logins, carts and other per-visitor data don't belong in pages, because pages are shared. Serve them from your own handlers, and fetch them with JavaScript. Send those responses with `Cache-Control: private, no-store`. Protect changes with `SameSite` cookies or a CSRF token.