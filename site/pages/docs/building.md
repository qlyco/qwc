<!--
section: docs
title: Building
description: Build a static copy of a qwc site for any static host.
layout: doc
sort_key: 5
tags: deployment
-->
# Building

<pre tabindex="0"><copy-code></copy-code><code>qwc -build</code></pre>

This writes the whole site to `_qwcsite`, including every page number, the 404 pages and static files, then exits. Upload that folder to any static host.

<div class="table-scroll" role="region" aria-label="Build commands" tabindex="0">

| Command | Output | Path |
|---|---|---|
| `qwc -build` | `./_qwcsite` | `/` |
| `qwc -build -output=deploy` | `./deploy/_qwcsite` | `/` |
| `qwc -build -url=user.github.io/project` | `./_qwcsite` | `/project/` |

</div>

`-url` is the address the site will live at. Its path is where links start, so `relativeURL` and `absoluteURL` use it. Without it, links start at `/`.

Markdown pages are written as HTML files, so `pages/blog/hello.md` becomes `_qwcsite/blog/hello.html`.

## Apps in a static build

The build writes each `200.html` as both `index.html` and `404.html` in its folder, so the app loads for the folder itself and for missing URLs in it.

## Serving instead

You can also serve the site without a build. Pages are cached, and a reload picks up changes without a restart.