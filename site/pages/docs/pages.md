<!--
section: docs
title: Pages and Templates
description: Front matter, layouts, includes, components, site data, collections and pagination in qwc.
layout: doc
sort_key: 2
tags: guides
-->
# Pages and Templates

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
| `sort_key` | Sorts collections when there is no `date`. Pages without it come last |
| `layout` | The layout the page goes into |
| `templates` | `false` sends the page as written, without running templates |
| `visibility` | `visible`, `hidden` or `none`. See [visibility](#visibility) |
| `type` | `text` sends the page as a [text page]({{relativeURL "/docs/text-pages"}}#sending-a-page-as-text) |
| `paginate` | Turns on [pagination](#pagination) |
| `page_size` | Items on each page. Defaults to 10 |
| `tags` | Comma-separated tags, like `go, web`. See [tags](#tags) |

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

Listed pages show in `collection`, `pages` and pagination. A `none` page acts as if its file is missing.

Hidden pages are still built, and anyone with the URL can open them. To keep search engines away, add this to the layout's head:

<pre tabindex="0"><code>{{`{{if eq .visibility "hidden"}}<meta name="robots" content="noindex">{{end}}`}}</code></pre>

A `none` 404 or 200 page is skipped too.

## Layouts

`layout: post` puts the page into `layouts/post.html`. The layout shows the page body where it has a `content` block:

<pre tabindex="0"><copy-code></copy-code><code>{{`<!doctype html>
<title>{{.title}}</title>
{{block "head" .}}{{end}}
<main>
  {{block "content" .}}{{end}}
</main>`}}</code></pre>

A page can fill the layout's other blocks with <code>{{`{{define "head"}}...{{end}}`}}</code>. Everything outside those goes into `content`. A layout that has no `content` block does not show the page body.

Use `layout:` to pick a layout. Each layout's blocks are separate, so two layouts can use the same block names.

## Layouts that build on others {#layout-inheritance}

A layout can name its own layout in front matter, the same way a page does. Here `post` builds on `base`:

<pre tabindex="0"><copy-code></copy-code><code>{{`<!--
layout: base
-->
{{define "head"}}<link rel="stylesheet" href="/post.css">{{end}}
<article>{{block "content" .}}{{end}}</article>`}}</code></pre>

A page with `layout: post` goes into the `<article>`, and the `<article>` goes into the `content` block of `base`. When the same block is filled more than once, the page wins over `post`, and `post` wins over `base`.

A layout only reads `layout:` from its front matter. A missing layout, or a loop, ends the chain.

For boilerplate like the head and footer, includes are often simpler than a base layout.

## Includes and components

Templates are named after their path, so `includes/header.html` is `includes/header`. An include is shared markup, and a component is one that takes props. Pass them with `dict`:

<pre tabindex="0"><code>{{`{{template "includes/header" .}}
{{template "components/note" (dict "type" "info" "text" "Saved.")}}`}}</code></pre>

Includes and components share names. Give `define` names a prefix, like `card.title`, so they don't clash with other files.

A page can replace an include for itself by defining its name, like <code>{{`{{define "includes/header"}}...{{end}}`}}</code>.

## When something goes wrong {#errors}

A template mistake doesn't stop the site. qwc shows the page up to the mistake, puts <code>[qwc: error]</code> there, and logs the full message. Builds do the same.

Use `-verbose` to show the full message on the page too. Don't use it on a live site.

Bad front matter falls back to a default, and the log says so:

<div class="table-scroll" role="region" aria-label="Front matter mistakes" tabindex="0">

| Mistake | Treated as |
|---|---|
| Unknown `visibility` | `none`, so a typo never publishes a draft |
| Unknown `type` | `page` |
| `layout` that does not exist | No layout |
| Bad `date` | No date |
| Bad `page_size` | 10 |

</div>

## Functions

<div class="table-scroll" role="region" aria-label="Template functions" tabindex="0">

| Function | Returns |
|---|---|
| `site` | Values from `site.json`, plus `site.data` from the `data` folder |
| `collection "blog/**"` | Pages matching a pattern, newest first |
| `tagged "go"` | Visible pages with that tag, newest first |
| `tags` | Every tag in use, sorted, with `.Name`, `.Count` and `.URL` |
| `pages "**"` | The same, with index pages too, for a [sitemap]({{relativeURL "/docs/text-pages"}}#example-a-sitemap) |
| `slice list 0 4` | Part of a list or text. It never fails when the list is shorter |
| `dict "key" value` | Props for a component |
| `relativeURL "/blog/"` | A link inside the site |
| `absoluteURL "/blog/"` | The same link with the domain, for canonical links and feeds |
| `plugin "db"` | A value from your own Go code, see [extending qwc]({{relativeURL "/docs/extending-qwc"}}) |
| `upperCase "sample_text"` | `{{upperCase "sample_text"}}` |
| `lowerCase "SAMPLE_TEXT"` | `{{lowerCase "SAMPLE_TEXT"}}` |
| `titleCase "sample text"` | `{{titleCase "sample text"}}` |
| `replace "_" " " "sample_text"` | `{{replace "_" " " "sample_text"}}` |
| `split "_" "sample_text"` | A list of parts: `{{$result := split "_" "sample_text"}}{{index $result 0}}, {{index $result 1}}` |
| `join " " list` | The items in the list as one text, with the separator between them: `{{join " " $result}}` |

</div>

Use `relativeURL` for links inside the site, so it works from a subfolder.

## Site data

`site.json` holds values for every page, like the site title. Each JSON file in `data` is added under `site.data`, so `data/team.json` is `site.data.team` and `data/shop/items.json` is `site.data.shop.items`.

## Collections

`collection` lists pages by path. `*` matches one part of a path and `**` matches any number of folders. Index pages are left out. Lists are sorted by `date`, newest first, then by `sort_key`.

<pre tabindex="0"><code>{{`{{range slice (collection "blog/**") 0 4}}
  <a href="{{relativeURL .url}}">{{.title}}</a>
{{end}}`}}</code></pre>

Use `else` for an empty list:

<pre tabindex="0"><code>{{`{{range collection "blog/**"}}
  <a href="{{relativeURL .url}}">{{.title}}</a>
{{else}}
  <p>No posts yet.</p>
{{end}}`}}</code></pre>

## Tags {#tags}

Add tags to a page's front matter as one comma-separated line:

<pre tabindex="0"><code>{{`<!--
title: Hello
tags: go, web
-->`}}</code></pre>

Spaces around each tag are trimmed, and repeated tags are dropped. A tag with a `/` in it is ignored, since tags become part of the URL.

`tagged "go"` returns the visible pages with that tag, newest first. `tags` returns every tag, sorted by name:

<pre tabindex="0"><code>{{`{{range tags}}
  <a href="{{relativeURL .URL}}">{{.Name}}</a> ({{.Count}})
{{end}}`}}</code></pre>

qwc makes tag pages from two layouts. Missing layouts are skipped:

<div class="table-scroll" role="region" aria-label="Tag layouts" tabindex="0">

| Layout | Makes | Values |
|---|---|---|
| `layouts/term.html` | One page for each tag. For a tag named `go`, the page is at `/tags/go/` | `.tag` is the tag name |
| `layouts/taxonomy.html` | One page listing every tag, at `/tags/` | `.pagination.Items` holds the tags on this page |

</div>

This `term.html` works for every tag. It lists the pages with the current tag:

<pre tabindex="0"><code>{{`<h1>Posts tagged {{.tag}}</h1>
{{range tagged .tag}}
  <a href="{{relativeURL .url}}">{{.title}}</a>
{{end}}`}}</code></pre>

A page in `pages/` with the same URL is shown instead of the generated page. For a tag named `go`, `pages/tags/go.html` replaces `/tags/go/`. Only `visible` pages add their tags.

Tag pages always paginate, so their layouts don't need a `paginate` line. Use `.pagination.Items`, `.pagination.NextURL` and `.pagination.PrevURL` as on other pages. Tag pages show 10 items per page by default. To change it, set `page_size` in the layout's front matter, or in `site.json` for every page that doesn't set its own.

To list all tags, or all pages with a tag, without paging, use `tags` or `tagged`.

The tag folder is `tags` by default. To use another folder, set `tags_slug` in `site.json`, like `"tags_slug": "topics"`. That moves the list to `/topics/`, and a tag named `go` moves to `/topics/go/`.

## Pagination

Add `paginate` with a pattern for the pages to list, and qwc serves `/page/2`, `/page/3` and so on. Patterns are relative to the page's folder. On `pages/blog/index.html`, `*` lists the pages directly in `blog/` and `**` lists all pages under it. Start with `/` to count from `pages` instead, so `/blog/**` works from any folder:

<pre tabindex="0"><code>{{`<!--
title: Blog
paginate: blog/**
page_size: 10
-->
{{range .pagination.Items}}
  <a href="{{relativeURL .url}}">{{.title}}</a>
{{end}}
{{with .pagination.NextURL}}<a href="{{relativeURL .}}">Next</a>{{end}}`}}</code></pre>

`.pagination` also has `Page`, `TotalPages`, `Total` and `PrevURL`. Page numbers past the end are a 404. `/page/1` redirects to the list.

## Missing pages and apps {#missing-pages}

If nothing matches a URL, qwc checks its folder, then each parent folder up to the root, for these pages:

<div class="table-scroll" role="region" aria-label="Missing page files" tabindex="0">

| Page | Serves | Status |
|---|---|---|
| `404.html` | That page, for anything missing in the folder | 404 |
| `200.html` | That page, for the folder itself and every missing URL in it | 200 |

</div>

A `200.html` makes its folder a single page app. `/app/`, `/app/settings` and `/app/users/42` all get `pages/app/200.html`, and the app's router takes it from there. Missing files with an extension, like `/app/main.js`, still get a 404.

A `200.html` only works in a folder with no index or 404 page. Both can also be Markdown, as `404.md` and `200.md`. If nothing matches, qwc sends a 404 status with no page. These files don't show up in collections.