<!--
layout: default
-->
# qwc

<p class="lead">A small web composer for HTML and Markdown sites.</p>

<p><a class="button" href="{{site.repo}}/releases/latest" target="_blank" rel="noopener">Download</a></p>

qwc (pronounced as quick) is a small and simple tool for turning a folder of pages into a website. Run it as a server and it works like a flat file CMS: each page is built on its first request and kept until you reload. Or run one command and get a static copy for any host.

Pages are HTML or Markdown with Go templates. There is no build tool, no bundler, no database, and no JavaScript unless you add it.

## What you get

- Layouts, includes and components, written as Go templates
- Markdown pages that can mix in HTML, includes and components
- Front matter in an HTML comment at the top of each page
- Collections of pages, sorted by date or `sort_key`, with real paginated URLs
- Site wide values from `site.json` and a `data` folder
- A cache that holds pages until you press <kbd>r</kbd>
- A build mode that writes a static site for any host
- A Go library, so the same engine can run inside your own program

## Quick start

<pre tabindex="0"><copy-code></copy-code><code>go install github.com/qlyco/qwc@latest
qwc</code></pre>

That serves the `public` folder on <a href="http://localhost:8080" target="_blank" rel="noopener">localhost:8080</a>, creating the folder if it is not there yet. Add a page:

<pre tabindex="0"><copy-code></copy-code><code>{{`<!--
title: Hello
-->
# Hello, qwc

This page is written in **Markdown**.`}}</code></pre>

Save it as `public/pages/hello.md` and open `/hello`.

## Small on purpose

qwc has minimal dependencies. Its templates are Go's <a href="https://pkg.go.dev/html/template" target="_blank" rel="noopener"><code>html/template</code></a>, so output is escaped by default and there is nothing new to learn.

[Read the docs]({{relativeURL "/docs/getting-started"}}) or <a href="{{site.repo}}" target="_blank" rel="noopener">see the source</a>.

## Disclaimer

qwc is in active pre-1.0 development. Although the current implementation is stable and free of known major issues, the core API and feature sets are still evolving. Expect potential breaking changes in updates leading up to v1.0.0.