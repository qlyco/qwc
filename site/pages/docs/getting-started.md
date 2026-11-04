<!--
section: docs
title: Getting Started
description: Install qwc, run it, and learn how a qwc site folder is laid out.
layout: doc
sort_key: 1
tags: setup
-->
# Getting Started

## Install

<pre tabindex="0"><copy-code></copy-code><code># stable numbered releases
go install github.com/qlyco/qwc@latest

# up-to-date development releases
go install github.com/qlyco/qwc@main
</code></pre>

Pre-compiled Linux binaries for stable builds are also attached to the <a href="{{site.repo}}/releases/latest" target="_blank" rel="noopener">Github Releases page</a> if available.

## Run

<pre tabindex="0"><code>qwc                          serve ./public on port 8080
qwc -root mysite -port 3000  serve another folder on another port
qwc -build                   build a static copy and exit</code></pre>

<div class="table-scroll" role="region" aria-label="Flags" tabindex="0">

| Flag | Usage | Default |
|---|---|---|
| `-root` | Folder that holds the site | `public` |
| `-port` | Port to serve on | `8080` |
| `-build` | Build the site as static files, then exit | off |
| `-output` | Folder to build `_qwcsite` in | current folder |
| `-url` | Address the site is served from | none |
| `-cert` | Certificate file. Serves HTTPS together with `-key` | none |
| `-key` | Private key file for `-cert` | none |
| `-config` | JSON file with settings | none |
| `-verbose` | Show full error messages on pages | off |

</div>

`-url` also works when serving. It sets the domain for `absoluteURL` links. A path, like `-url=example.com/docs`, serves the site under `/docs/`.

## HTTPS

Give qwc a certificate and its key, and it serves HTTPS instead of HTTP on the same port:

<pre tabindex="0"><copy-code></copy-code><code>qwc -cert=cert.pem -key=key.pem</code></pre>

Both files are needed. qwc reloads them when they change, for example after a certbot renewal, so no restart is needed. If the new files can't load, qwc keeps the old certificate and logs why.

## Config files

Settings can live in a JSON file instead of flags. Each key is a flag name without the `-`:

<pre tabindex="0"><copy-code></copy-code><code>{
  "root": "public",
  "port": 8080,
  "url": "",
  "output": ".",
  "cert": "",
  "key": "",
  "verbose": false
}</code></pre>

Pass it with `-config`. Use one file per setup, like `dev.json` or `prod.json`:

<pre tabindex="0"><code>qwc -config=prod.json
qwc -config=prod.json -build
qwc -config=prod.json -port 3000</code></pre>

- A flag you type wins over the file, so the last line uses port 3000.
- Paths are read from the folder you run qwc in. Run it from your project folder.
- `build` and `config` cannot be set in the file.
- A missing file, broken JSON or an unknown key stops qwc with an error.

## The site folder

<pre tabindex="0"><code>pages/        one HTML or Markdown file per page, the folders become the URLs
static/       CSS, JS, images and other files, served as they are
layouts/      the outer structure shared by pages
includes/     parts shared between layouts, like a header
components/   reusable markup that takes props
data/         JSON files, available in templates as site.data
site.json     values for the whole site</code></pre>

Only `pages` is needed to start. A page at `pages/blog/hello.html` or `pages/blog/hello.md` is served at `/blog/hello`, and `pages/blog/index.html` at `/blog/`.

A `404.html` in `pages` is used for pages that do not exist. Without one, qwc answers with a bare 404 status. Folders can also have their own 404 page, and a folder can be a single page app. See [missing pages and apps]({{relativeURL "/docs/pages"}}#missing-pages).

## Seeing your changes

qwc builds each page on its first request and caches it in `.qwccache`. After you change something, press a key:

<div class="table-scroll" role="region" aria-label="Keys" tabindex="0">

| Key | Action |
|---|---|
| <kbd>r</kbd> | Reload. Cached files rebuild on their next request. |
| <kbd>R</kbd> | Refresh. Clears the cache, then reloads. |
| <kbd>Control/Ctrl</kbd>+<kbd>C</kbd> | Stop. |

</div>

Deleting the `_qwcbuild` file in the site folder also reloads. This is useful from a deploy script.