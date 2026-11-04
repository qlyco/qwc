<!--
section: docs
title: Docs
description: Guides, setup and deployment for qwc.
layout: default
paginate: **
-->
<div class="list-page">

# Docs

Guides, setup and deployment for qwc.

<ul>
{{range .pagination.Items}}<li><a href="{{relativeURL .url}}">{{.title}}</a></li>
{{end}}</ul>
<hr class="list-rule">
<nav class="doc-nav list-nav" aria-label="Docs list pages">
  {{with .pagination.PrevURL}}<a href="{{relativeURL .}}" rel="prev">Previous</a>{{else}}<span></span>{{end}}
  <span>{{.pagination.Page}} / {{.pagination.TotalPages}}</span>
  {{with .pagination.NextURL}}<a href="{{relativeURL .}}" rel="next">Next</a>{{end}}
</nav>

</div>