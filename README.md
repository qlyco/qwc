# qwc: web composer

## overview

[![Release](https://github.com/qlyco/qwc/actions/workflows/release.yml/badge.svg)](https://github.com/qlyco/qwc/actions/workflows/release.yml)
[![Deploy](https://github.com/qlyco/qwc/actions/workflows/site.yml/badge.svg)](https://github.com/qlyco/qwc/actions/workflows/site.yml)

qwc is a static site generator / flat-file cms written in go.

Pages are HTML or Markdown files with Go templates. qwc builds each page on its first request and caches it until you reload. It can also build the whole site into static files for any host.

## install

```sh
go install github.com/qlyco/qwc@latest
```

## usage

```sh
qwc                          # serve ./public on port 8080
qwc -root mysite -port 3000  # serve another folder on another port
qwc -build                   # build ./_qwcsite and exit
qwc -build -output=deploy    # build ./deploy/_qwcsite and exit
qwc -build -url=user.github.io/project   # build for a site in a subfolder
qwc -url=example.com         # serve with full links for example.com
qwc -cert=cert.pem -key=key.pem          # serve HTTPS
qwc -config=prod.json        # read settings from a JSON file, see example.qwcconfig.json
qwc -verbose                 # show full error messages on pages
```

While serving:

| key | action |
|---|---|
| `r` | reload, cached files rebuild on their next request |
| `R` | refresh, clear the cache, then reload |
| `Control/Ctrl+C` | stop |

Deleting `_qwcbuild` in the site folder also reloads.

## site folder

```
pages/        HTML or Markdown pages, one per URL
static/       CSS, JS, images and other files, served as they are
layouts/      shared page structure
includes/     parts shared between layouts
components/   reusable markup that takes props
data/         JSON files, available as site.data
site.json     values for the whole site
```

qwc creates `.qwccache/` and `_qwcbuild` itself.

## pages

A page is an `.html` or `.md` file. If both exist, the `.html` one is used. A page can start with front matter in an HTML comment:

```html
<!--
title: Hello
date: 2026-11-06
layout: post
-->
<h1>Hello</h1>
```

`layout: post` puts the page body into the `content` block of `layouts/post.html`. A layout can build on another one with its own `layout:` key. In a `.md` page, templates run first, then the Markdown is turned into HTML. `templates: false` turns templates off for a page.

`visibility: hidden` keeps a page out of collections and `pages`, but its URL still works. `visibility: none` treats the page as missing, like a draft.

Templates can use:

| function | returns |
|---|---|
| `site` | values from `site.json` and `data/` |
| `collection "blog/**"` | pages matching a pattern, newest first |
| `pages "**"` | the same, with index pages too, for a sitemap |
| `slice list 0 4` | part of a list, never fails on short lists |
| `dict "key" value` | props for a component |
| `relativeURL "/blog/"` | a link inside the site, with the `-url` path added |
| `absoluteURL "/blog/"` | the same with the domain, for canonical links and feeds |

Front matter `paginate: blog/**` and `per_page: 4` turn on `/blog/page/2` and so on.

A mistake in a template shows `[qwc: error]` where it happened, and the full message goes to the log. Front matter mistakes fall back to safe defaults, and an unknown `visibility` is treated as `none`.

Other text files in `pages/`, like `sitemap.xml` or `robots.txt`, are text pages. They work like other pages, without HTML escaping, are served under their own name, and are never listed. `type: text` makes an `.html` or `.md` page a text page.

For missing URLs, qwc walks up to the nearest folder with one of these:

| file | serves | status |
|---|---|---|
| `404.html` | that page, for anything missing in the folder | 404 |
| `200.html` | that page, for the folder and every missing URL in it, as a single page app | 200 |

A `200.html` only counts in a folder without an index or 404 page. Both can also be `.md` files. If nothing matches, qwc sends a 404 status with no page.

## library

```go
import "github.com/qlyco/qwc/core"

site := core.New("public", nil)

http.ListenAndServe(":8080", site)
```

| method | same as |
|---|---|
| `site.Reload()` | pressing `r` |
| `site.Refresh()` | pressing `R` |
| `site.SetURL("example.com")` | `qwc -url=example.com` |
| `site.Build("deploy", "")` | `qwc -build -output=deploy` |
| `site.SetVerbose(true)` | `qwc -verbose` |
| `site.RegisterPlugin("db", value)` | adds a plugin, used in templates as `plugin "db"` |

To make your own qwc command with plugins, keep `main` small and pass your name, version and a setup function to `cli.Run`:

```go
func main() {
	cli.Run(cli.App{Name: "mysite", Version: Version, Setup: func(site *core.Server) error {
		return site.RegisterPlugin("db", store{})
	}})
}
```

`cli.App` also takes an optional `Handler` that wraps the site, for extra routes or middleware.
