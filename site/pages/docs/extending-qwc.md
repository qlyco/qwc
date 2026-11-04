<!--
section: docs
title: Extending qwc
description: Add plugins to qwc, make your own qwc command, or run qwc inside your own Go program.
layout: doc
-->
# Extending qwc

qwc is a Go package, and the `qwc` command is a small program built on it. You can extend it in two ways:

<div class="table-scroll" role="region" aria-label="Ways to extend qwc" tabindex="0">

| Way | Good for | You write |
|---|---|---|
| [Your own qwc command](#your-own-command) | qwc as it is, plus plugins or a few flags | A small `main` |
| [qwc inside your program](#inside-your-program) | A bigger program that also serves a site, with its own routes and messages | The program around it |

</div>

Both start with the same thing:

<pre tabindex="0"><copy-code></copy-code><code>go get github.com/qlyco/qwc</code></pre>

## Plugins {#plugins}

A plugin gives templates data from Go, like products from a database. It is any Go value, usually a struct with methods. This one reads products with Go's `database/sql`:

<pre tabindex="0"><copy-code></copy-code><code>{{`type Product struct {
	Name  string
	Price int
}

type store struct {
	database *sql.DB
}

func (shop store) Products() ([]Product, error) {
	rows, err := shop.database.Query("SELECT name, price FROM products")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product

		err = rows.Scan(&product.Name, &product.Price)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, rows.Err()
}

func (shop store) Product(id int) (Product, error) {
	var product Product

	err := shop.database.QueryRow("SELECT name, price FROM products WHERE id = ?", id).Scan(&product.Name, &product.Price)

	return product, err
}`}}</code></pre>

Register it under a name made of letters and digits. Here `openDatabase` stands for however your program connects to its database:

<pre tabindex="0"><copy-code></copy-code><code>{{`database := openDatabase()

err := site.RegisterPlugin("db", store{database: database})`}}</code></pre>

Templates get it with `plugin`, best saved in a variable. Each method becomes something templates can call:

<pre tabindex="0"><copy-code></copy-code><code>{{`{{$db := plugin "db"}}
<ul>{{range $db.Products}}<li>{{.Name}}: {{.Price}}</li>{{end}}</ul>
<p>Featured: {{($db.Product 5).Name}}</p>`}}</code></pre>

An unknown name shows <code>[qwc: error]</code> at that spot. Variables only last within one template, so an include calls `plugin` itself or gets the value through `dict`.

Plugins can be shared as normal Go packages that return a value to register.

### Rules for plugins

- **Pages are cached.** A plugin runs when a page is built, and every visitor gets that result until `site.Reload()`. Use plugins for data that is the same for everyone.
- **Pages build in parallel.** Methods must be safe to call at the same time. A `*sql.DB` is.
- **Errors show on the page** with `-verbose`. Return short messages, and log the details yourself.

## Your own qwc command {#your-own-command}

The `cli` package is the whole qwc command: flags, config files, keys, HTTPS and build mode. Give it your name, your version and a setup function:

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
			return site.RegisterPlugin("db", store{database: openDatabase()})
		},
	})
}`}}</code></pre>

- `Setup` runs before qwc serves or builds, so plugins work in both.
- `Handler` is optional. It gets the site and returns what qwc serves, so you can add routes or wrap every request in middleware. It is not used by `-build`.
- Flags you define with Go's `flag` package before calling `cli.Run` work like qwc's own, including in config files.
- The banner shows your name and version, and the qwc version it is built on, like `mysite [v1.2.0], built on qwc v0.0.1`. Set your version with `-ldflags "-X main.Version=v1.2.0"`.

For example, an API route next to the site, and a header on every response:

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

`saveProduct` and `withSecurityHeaders` stand for your own handler and middleware. Pass requests on to the site with their full path, so keep it at `/` and don't use `http.StripPrefix`.

**Use case:** a shop site whose product pages come from a database. You add a `db` plugin and a route to save products, and keep everything else about qwc: `-build` for the static copy, `-config` for each server, and <kbd>r</kbd> while writing.

## qwc inside your program {#inside-your-program}

A site is an `http.Handler`, so it fits into any Go server, next to your own routes and middleware:

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

`Build` returns the path of the `_qwcsite` folder it wrote. If the site has a path, like `example.com/docs`, pass requests on with their full path. qwc removes the path itself, so don't use `http.StripPrefix`.

**Use case:** an app with logins and an API that also has a blog and docs. Your program owns the routes, sessions and messages, and qwc only serves the pages.

## Changing data

After changing what a plugin reads, call `site.Reload()`. Pages rebuild on their next request:

<pre tabindex="0"><copy-code></copy-code><code>{{`mux.HandleFunc("POST /api/products", func(writer http.ResponseWriter, request *http.Request) {
	// check the login, then save the product
	site.Reload()
})`}}</code></pre>

In your own qwc command, add the route in `Handler`. A separate program can also write the data and then delete `_qwcbuild` in the site folder, which reloads too.

## Data for each visitor

Logins, carts and other per-visitor data do not belong in pages, since pages are shared. Serve them from your own handlers, and fetch them from the page with JavaScript. Same-origin requests send cookies, including `HttpOnly` ones. Send those responses with `Cache-Control: private, no-store`, and protect changes with `SameSite` cookies or a CSRF token.