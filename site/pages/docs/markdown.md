<!--
section: docs
title: Markdown
description: Write qwc pages in Markdown, mixed with HTML, includes and components.
layout: doc
-->
# Markdown

A page can be a `.md` file instead of `.html`. `pages/blog/hello.md` is served at `/blog/hello`, the same as `pages/blog/hello.html` would be. If a folder has both, the `.html` one is used.

Layouts, includes and components stay HTML.

## Writing a page

Front matter is the same HTML comment as in HTML pages. Give the page a `layout`, and the Markdown goes into its `content` block:

<pre tabindex="0"><copy-code></copy-code><code>{{`<!--
title: Hello
date: 2026-11-06
layout: post
-->
# Hello

This is **Markdown**, with a [link](/about) and a list:

- one
- two`}}</code></pre>

Without a layout, qwc sends the HTML made from the Markdown as it is, with no `<head>`. That is fine for a quick test, but real pages should use a layout.

To send the Markdown itself instead of HTML, use [`type: text`]({{relativeURL "/docs/text-pages"}}#sending-a-page-as-text).

## What works

<div class="table-scroll" role="region" aria-label="Markdown features" tabindex="0">

| Feature | Written as |
|---|---|
| Headings, lists, links, images, quotes, code | Standard Markdown |
| Tables | <code>&#124; a &#124; b &#124;</code> rows |
| Strikethrough | `~~old~~` |
| Task lists | `- [x] done` |
| Links from bare addresses | `https://example.com` |
| Footnotes | `text[^1]` and `[^1]: note` |
| Definition lists | A term, then `: definition` on the next line |
| Heading IDs | Made from the heading text, or set with `{#id .class}` |
| Smart quotes and dashes | `"quotes"`, `--` and `...` |
| HTML | Written as it is |

</div>

Code blocks get a class for their language, like `class="language-go"`, so a highlighting script or stylesheet can color them.

## Mixing in HTML and templates

Templates run first, then the Markdown is turned into HTML. So a Markdown page can use everything an HTML page can: site values, collections, includes and components.

<pre tabindex="0"><copy-code></copy-code><code>{{`## Latest posts

{{range slice (collection "blog/*") 0 3}}
- [{{.title}}]({{relativeURL .url}})
{{- end}}

{{template "components/note" (dict "text" "Saved.")}}

<div class="note">

Markdown *inside* HTML works with blank lines around it.

</div>`}}</code></pre>

## Things to watch

- **Markdown inside HTML** needs a blank line after the opening tag and before the closing tag, as in the example above. Without them, the Markdown is left as plain text.
- **Includes and components** should not have blank lines in their HTML. After a blank line, a line indented by 4 spaces or more turns into a code block.
- **Template code in examples** runs, like any other template code. To show it, print it as a raw string inside an HTML `<pre tabindex="0"><code>` block, as below. A Markdown code block would show the escaped quotes as written.

<pre tabindex="0"><code>{{"<pre><code>{{`{{.title}}`}}</code></pre>"}}</code></pre>

## Turning templates off

`templates: false` in the front matter sends the page without running templates. The Markdown is still turned into HTML, and the layout is still used. This is handy for pages full of code examples, as long as their links do not need `relativeURL`.

<pre tabindex="0"><code>{{`<!--
title: Template cheat sheet
layout: doc
templates: false
-->
Write {{.title}} to show the title.`}}</code></pre>
