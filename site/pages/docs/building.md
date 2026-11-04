<!--
section: docs
title: Building
description: Build a static copy of a qwc site for any static host.
layout: doc
-->
# Building

<pre tabindex="0"><copy-code></copy-code><code>qwc -build</code></pre>

This writes the whole site into `_qwcsite`: every page, every page number, the 404 pages and the static files. Then it exits. Upload the folder to any static host.

<div class="table-scroll" role="region" aria-label="Build commands" tabindex="0">

| Command | Output | Path |
|---|---|---|
| `qwc -build` | `./_qwcsite` | `/` |
| `qwc -build -output=deploy` | `./deploy/_qwcsite` | `/` |
| `qwc -build -url=user.github.io/project` | `./_qwcsite` | `/project/` |

</div>

`-url` is the address the site will live at. Path is what links start with. It sets the path that `relativeURL` adds, and the domain that `absoluteURL` adds. Without it, links start at `/`.

Markdown pages are written as HTML files, so `pages/blog/hello.md` becomes `_qwcsite/blog/hello.html`.

## Apps in a static build

The build writes each `200.html` as both `index.html` and `404.html` in its folder, so the app loads for the folder itself and for missing URLs in it.

## Serving instead

qwc can also run as the site's server, with no build step. Pages are built once and cached, and a reload brings in new changes without a restart.
