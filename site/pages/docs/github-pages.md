<!--
section: docs
title: GitHub Pages
description: Build and publish a qwc site to GitHub Pages with the qwc workflow.
layout: doc
-->
# GitHub Pages

qwc has a reusable workflow that builds a site and publishes it to GitHub Pages.

## Add the workflow

In your site's repository, add `.github/workflows/deploy.yml`:

<pre tabindex="0"><copy-code></copy-code><code>name: Deploy site

on:
  push:
    branches:
      - main
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: false

jobs:
  pages:
    uses: qlyco/qwc/.github/workflows/deploy.yml@v0.0.1
    with:
      version: v0.0.1
      root: public</code></pre>

<div class="table-scroll" role="region" aria-label="Workflow inputs" tabindex="0">

| Input | Description |
|---|---|
| `version` | The qwc release to build with |
| `root` | The site folder in your repository |

</div>

Keep the version after `@` the same as `version`.

## Turn on Pages

In the repository, go to **Settings > Pages** and set **Source** to **GitHub Actions**. After that, every push to `main` builds and publishes the site. You can also run it by hand from the **Actions** tab.

The workflow passes the right `-url` for you, so project sites at `user.github.io/project/` work too.