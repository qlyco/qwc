<!--
section: docs
title: Pages and templates
description: Front matter, layouts, includes, components, site data, collections and pagination in qwc.
layout: doc
-->
# Pages and templates

Pages, layouts, includes and components are Go templates. See the <a href="https://pkg.go.dev/text/template" target="_blank" rel="noopener">Go template docs</a> for the syntax. Pages can be HTML or [Markdown]({{relativeURL "/docs/markdown"}}), and layouts, includes and components are HTML.

## Front matter

A page can start with an HTML comment holding `key: value` lines. The values are available in the page and its layout as `.key`.

<pre tabindex="0"><code>{{`<!--
title: Hello
date: 2026-11-06
summary: The first post.
layout: post
-->`}}</code></pre>

Any key can be used. These ones mean something to qwc:

<div class="table-scroll" role="region" aria-label="Front matter keys" tabindex="0">

| Key | Usage |
|---|---|
| `date` | Sorts collections. Takes `YYYY-MM-DD` or an RFC 3339 time |
| `layout` | The layout the page goes into |
| `templates` | `false` sends the page as written, without running templates |
| `visibility` | `visible`, `hidden` or `none`. See [visibility](#visibility) |
| `type` | `text` sends the page as a [text page]({{relativeURL "/docs/text-pages"}}#sending-a-page-as-text) |
| `paginate` | Turns on [pagination](#pagination) |
| `per_page` | Items on each page number, 10 if not set |

</div>

qwc also sets two values itself, so don't use these names in front matter:

- `.url` is the page's URL.
- `.pagination` is set on [paginated](#pagination) pages.

## Visibility {#visibility}

<div class="table-scroll" role="region" aria-label="Visibility values" tabindex="0">

| Value | URL | Listed | Use |
|---|---|---|---|
| `visible` | Works | Yes | Normal pages. This is the default |
| `hidden` | Works | No | Pages you only link to, like a thank you page |
| `none` | Missing | No | Drafts |

</div>

Listed means the page shows up in `collection`, `pages` and pagination. A `none` page is treated as if its file is not there.

A hidden page is still built and anyone with its URL can open it. To keep search engines away too, add this to the layout's head:

<pre tabindex="0"><code>{{`{{if eq .visibility "hidden"}}<meta name="robots" content="noindex">{{end}}`}}</code></pre>

A `none` 404 or 200 page is skipped too, so qwc keeps walking up to the next one.

## Layouts

`layout: post` puts the page into `layouts/post.html`. The layout shows the page body where it has a `content` block:

<pre tabindex="0"><copy-code></copy-code><code>{{`<!doctype html>
<title>{{.title}}</title>
{{block "head" .}}{{end}}
<main>
  {{block "content" .}}{{end}}
</main>`}}</code></pre>

A page can fill the layout's other blocks with <code>{{`{{define "head"}}...{{end}}`}}</code>. Everything outside those goes into `content`. A layout that has no `content` block does not show the page body.

`layout:` is the only way to use a layout. Each page only sees its own layout, so two layouts can use the same block names without mixing up their defaults.

## Layouts that build on others {#layout-inheritance}

A layout can name its own layout in front matter, the same way a page does. Here `post` builds on `base`:

<pre tabindex="0"><copy-code></copy-code><code>{{`<!--
layout: base
-->
{{define "head"}}<link rel="stylesheet" href="/post.css">{{end}}
<article>{{block "content" .}}{{end}}</article>`}}</code></pre>

A page with `layout: post` goes into the `<article>`, and the `<article>` goes into the `content` block of `base`. When the same block is filled more than once, the page wins over `post`, and `post` wins over `base`.

Only `layout:` is read from a layout's front matter. A layout that is missing, or a chain that loops back on itself, ends the chain there.

For boilerplate like the head and footer, includes are often simpler than a base layout.

## Includes and components

Templates are named after their path, so `includes/header.html` is `includes/header`. An include is shared markup, and a component is one that takes props. Pass them with `dict`:

<pre tabindex="0"><code>{{`{{template "includes/header" .}}
{{template "components/note" (dict "type" "info" "text" "Saved.")}}`}}</code></pre>

Includes and components share one set of names. If two of them define the same name, the last one loaded is used. Name defines after their file, like `card.title` in `components/card.html`, to keep them apart.

A page can replace an include for itself by defining its name, like <code>{{`{{define "includes/header"}}...{{end}}`}}</code>.

## When something goes wrong {#errors}

A mistake in a template does not stop the site. qwc shows the page up to the problem, puts <code>[qwc: error]</code> where it happened, and logs the full message. If nothing could be shown, the page is only the marker. The same happens in a build.

Start qwc with `-verbose` to see the full message on the page too. Keep it off for a live site, since messages can show details about the site.

Front matter mistakes fall back to a safe default, and the log says so:

<div class="table-scroll" role="region" aria-label="Front matter mistakes" tabindex="0">

| Mistake | Treated as |
|---|---|
| Unknown `visibility` | `none`, so a typo never publishes a draft |
| Unknown `type` | `page` |
| `layout` that does not exist | No layout |
| Bad `date` | No date |
| Bad `per_page` | 10 |

</div>

## Functions

<div class="table-scroll" role="region" aria-label="Template functions" tabindex="0">

| Function | Returns |
|---|---|
| `site` | Values from `site.json`, plus `site.data` from the `data` folder |
| `collection "blog/**"` | Pages matching a pattern, newest first |
| `pages "**"` | The same, with index pages too, for a [sitemap]({{relativeURL "/docs/text-pages"}}#example-a-sitemap) |
| `slice list 0 4` | Part of a list or text. It never fails when the list is shorter |
| `dict "key" value` | Props for a component |
| `relativeURL "/blog/"` | A link inside the site |
| `absoluteURL "/blog/"` | The same link with the domain, for canonical links and feeds |

</div>

Write every link inside the site with `relativeURL`, so the site still works when it is built for a subfolder.

## Site data

`site.json` holds values for every page, like the site title. Each JSON file in `data` is added under `site.data`, so `data/team.json` is `site.data.team` and `data/shop/items.json` is `site.data.shop.items`.

## Collections

`collection` lists pages by path. `*` matches one part of a path and `**` matches any number of folders. Index pages are left out, and the list is sorted by `date`, newest first.

<pre tabindex="0"><code>{{`{{range slice (collection "blog/**") 0 4}}
  <a href="{{relativeURL .url}}">{{.title}}</a>
{{end}}`}}</code></pre>

## Pagination

Add `paginate` to a page's front matter, and qwc serves `/page/2`, `/page/3` and so on after it:

<pre tabindex="0"><code>{{`<!--
title: Blog
paginate: blog/**
per_page: 10
-->
{{range .pagination.Items}}
  <a href="{{relativeURL .url}}">{{.title}}</a>
{{end}}
{{with .pagination.NextURL}}<a href="{{relativeURL .}}">Older</a>{{end}}`}}</code></pre>

`.pagination` also has `Page`, `TotalPages`, `Total` and `PrevURL`. Page numbers past the end are a 404, and `/page/1` redirects to the list itself.

## Missing pages and apps {#missing-pages}

When nothing matches a URL, qwc walks up from its folder to the root. The first folder with one of these pages decides:

<div class="table-scroll" role="region" aria-label="Missing page files" tabindex="0">

| Page | Serves | Status |
|---|---|---|
| `404.html` | That page, for anything missing in the folder | 404 |
| `200.html` | That page, for the folder itself and every missing URL in it | 200 |

</div>

A `200.html` makes its folder a single page app: `/app/`, `/app/settings` and `/app/users/42` all get `pages/app/200.html`, and the app's own router takes it from there. Missing files with an extension, like `/app/main.js`, still get a real 404.

A `200.html` only counts in a folder without an index or 404 page, so a folder is either pages or an app. Both can also be Markdown, as `404.md` and `200.md`. If nothing matches all the way up, qwc sends a 404 status with no page. Neither file is a page of its own, and neither shows up in collections.
