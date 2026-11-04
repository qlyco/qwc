<!--
section: docs
title: Getting started
description: Install qwc, run it, and learn how a qwc site folder is laid out.
layout: doc
-->
# Getting started

## Install

<pre tabindex="0"><copy-code></copy-code><code>go install github.com/qlyco/qwc@latest</code></pre>

Linux builds are also attached to the <a href="{{site.repo}}/releases/latest" target="_blank" rel="noopener">latest release</a>.

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

`-url` works when serving too. It adds the domain to `absoluteURL` links, like the ones in a sitemap. If it has a path, like `-url=example.com/docs`, qwc serves the site under `/docs/`.

## HTTPS

Give qwc a certificate and its key, and it serves HTTPS instead of HTTP on the same port:

<pre tabindex="0"><copy-code></copy-code><code>qwc -cert=cert.pem -key=key.pem</code></pre>

Both files are needed. When they are replaced, for example when certbot renews the certificate, qwc picks up the new ones by itself, with no restart. If the new files cannot be loaded, qwc keeps using the old certificate and logs why.

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

Pass the file with `-config`. Keep one file for each setup, like `dev.json` and `prod.json`:

<pre tabindex="0"><code>qwc -config=prod.json
qwc -config=prod.json -build
qwc -config=prod.json -port 3000</code></pre>

- A flag you type wins over the file, so the last line uses port 3000.
- Paths in the file are read from the folder you run qwc in, like the same flags typed. Run qwc from your project folder.
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

qwc builds each page the first time it is asked for and keeps it in `.qwccache`. Later requests get the saved copy, so they are about as fast as static files. After you change something, reload:

<div class="table-scroll" role="region" aria-label="Keys" tabindex="0">

| Key | Action |
|---|---|
| <kbd>r</kbd> | Reload. Cached files rebuild on their next request. |
| <kbd>R</kbd> | Refresh. Clears the cache, then reloads. |
| <kbd>Control/Ctrl</kbd>+<kbd>C</kbd> | Stop. |

</div>

Deleting the `_qwcbuild` file in the site folder also reloads, which is handy from a deploy script.