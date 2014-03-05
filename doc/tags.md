Tags
====

In order to organize tags in cfmedias, we make heavy use of SQL features.
A table named _tags_ contains every tag there is. A table named _items_tags_ contains links between items and tags. The entries in this table are linked to the items itself with `ON DELETE CASCADE`. Searching for tags is done by joining the _items_ table with _items_tags_ table and filtering for specific tags.

Occasionally -- or if there is considerable performance drop with tag-related queries -- we delete unreferenced tags which do not have any corresponding entries in _item_tags_.

We further enhance this capability to full-text-like search by cerating tags for artists, albums, track names and genres. We also feed tags from last.fm and/or libre.fm.
