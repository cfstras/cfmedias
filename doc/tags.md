Tags
====

In order to organize tags in cfmedias, we use SQL foreign keys and deletion triggers.
A table named `tags` contains every tag. A table named `items_tags` contains links between items and tags. The entries in this table are linked to the items itself with `ON DELETE CASCADE`. Searching for tags is done by joining the `items` table with `items_tags` table and filtering for specific tags.

Occasionally -- or if there is considerable performance drop with tag-related queries -- we delete unreferenced tags which do not have any corresponding entries in `item_tags`.

We further enhance this capability to full-text-like search by creating tags for artists, albums, track names and genres. We also feed tags from last.fm and/or libre.fm.
