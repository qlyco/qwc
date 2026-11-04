<!--
section: docs
title: Text Pages
description: Make sitemaps, feeds and other text files from templates in qwc.
layout: doc
sort_key: 4
tags: guides
-->
# Text Pages

A text page is any text file in `pages` that is not `.html` or `.md`, like `sitemap.xml`, `feed.xml`, `robots.txt` or `data.json`. qwc runs it as a template and serves it under its own name.

## How text pages work

A text page works like other pages, with a few differences:

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

qwc doesn't make a sitemap itself. A text page can. Save this as `pages/sitemap.xml`:

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
{{- range tags}}
  <url>
    <loc>{{absoluteURL .URL}}</loc>
  </url>
{{- end}}
  <url>
    <loc>{{absoluteURL (printf "/%s/" (or site.tags_slug "tags"))}}</loc>
  </url>
</urlset>`}}</code></pre>

`pages` is like `collection`, but it also lists index pages, so `/` and `/docs/` are included. 404 and 200 pages are left out, and so are `hidden` and `none` pages. Use `-url` so the links have the domain.

Tag pages aren't in `pages`, so the sitemap lists them with `tags`, using each tag's `.URL`. The last tag line is added by hand, and uses `tags_slug` from `site.json`, or `tags` if it is not set. Use these lines only if you have `layouts/term.html` and `layouts/taxonomy.html`. To skip tags with one post, wrap the tag's `<url>` in <code>{{`{{if gt .Count 1}}`}}</code> and <code>{{`{{end}}`}}</code>.

Feeds, a `robots.txt`, a JSON search index or a `site.webmanifest` made from `site.json` work the same way.