Basic Management
================

cfmedias keeps a database of all media.
It creates *packages* of things which get sha-256 sums for identification.
Machines in the local network can then share those packages in bittorrent-manner.


packages
--------
packages could be:

- music:
  - a single song
  - an album
  - a whole artist
  - a genre
- videos:
  - an episode of a series
  - a season of a series
  - a whole series
  - a movie
  - a movie and its sequels
- software:
  - a game package
  - a software package
  - a series of games

database
--------

software:

- sqlite
   - pro: fast, small.
   - con: only single-instance. hard to multi-thread.
- mysql
   - pro: easy to multithread, multi-user
   - con: bigger, not that fast
- postgres
   - pro: heard it's faster
   - con: complex setup
- hmm.

tables:

 + music, videos, etc. each its own table because they have different attributes
 + one table for packages
 + one for tags (how do we save taggings? how do we search for them fast?)
 + every user gets a seperate db (table/file?) containing his playcounts, favorites and lists

performance:
do we

 + keep the library requests to a minimum by caching data in memory or
 + keep only in memory what's on-screen and release everything else
 + let the user specify the amount of ram used?

