<!--
section: docs
title: Text pages
description: Make sitemaps, feeds and other text files from templates in qwc.
layout: doc
-->
# Text pages

A text page is any text file in `pages` that is not `.html` or `.md`, like `sitemap.xml`, `feed.xml`, `robots.txt` or `data.json`. qwc runs it as a template and serves it under its own name.

## How text pages work

A text page works like any other page. It can have front matter, a `layout`, includes and components, `templates: false` and `visibility`. Only a few things differ:

<div class="table-scroll" role="region" aria-label="Text pages" tabindex="0">

| Feature | Text page |
|---|---|
| Files | Any text type that is not `.html` or `.md`, like `.xml`, `.txt`, `.json` or `.csv` |
| URL | The file name as it is, so `pages/sitemap.xml` is `/sitemap.xml` |
| Templates | No HTML escaping, values are written exactly as they are |
| Collections and `pages` | Text pages are never listed |

</div>

## Sending a page as text

`type: text` in the front matter turns an `.html` or `.md` page into a text page. For example, this sends Markdown as it is, at `/notes.md`, without turning it into HTML:

<pre tabindex="0"><code>{{`<!--
title: Notes
type: text
-->
# {{.title}}

Raw Markdown, with templates still run.`}}</code></pre>

Add `templates: false` to send the file exactly as written.

## Example: a sitemap

qwc does not make a sitemap by itself, but a text page can. Save this as `pages/sitemap.xml`:

<pre tabindex="0"><copy-code></copy-code><code>{{`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
{{- range pages "**"}}
  <url>
    <loc>{{absoluteURL .url}}</loc>
    {{- with .date}}
    <lastmod>{{.Format "2006-01-02"}}</lastmod>
    {{- end}}
  </url>
{{- end}}
</urlset>`}}</code></pre>

`pages` works like `collection`, but also lists index pages, so `/` and `/docs/` are in the sitemap. 404 and 200 pages are left out, and so are pages with `visibility: hidden` or `none`. Build or serve with `-url` so the links have the domain.

Feeds, a `robots.txt`, a JSON search index or a `site.webmanifest` made from `site.json` work the same way.
