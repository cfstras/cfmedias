cfmedias [![Build Status](https://travis-ci.org/cfstras/cfmedias.svg?branch=master)](https://travis-ci.org/cfstras/cfmedias)
========
The media player and manager you always wanted.

Overview
--------
cfmedias aims to let you take control of your media library.  

Priorities:

- fast
- support all common media formats, both lossless and lossy
- iPod, Android and mp4-player sync
- incremental import & export of iTunes, Foobar2k & Google Play Music databases
- metadata, rating, play count & skip count import from Last.fm and all of the above
- highly portable, full and incremental backups
- play behaviour analysis and automatic playlist generator
- both simple and advanced or scriptable filter playlists

Status
------
_I'll try to keep this up to date._

What works:

- single-binary with web server & assets
- media scanner for music
- iPod sync: titles, no playlists
- HTTP and commandline API
- crude web-interface with API fiddle
- sqlite database (not final)
- audioscrobbler server endpoint

Compiling
---------
You will need the Go package, git, mercurial, make and npm.  
Also, you need development headers for libportaudio, libsqlite3, taglib and
libgpod.

    git clone --recursive https://github.com/cfstras/cfmedias
    cd cfmedias
    make run


License
-------
This software is released under the 2-clause BSD-license. For details, see LICENSE.md
Also, the author would love pull requests and reported issues.

Disclaimer
----------
*This program is far from finished.*
It may destruct your computer, your whole network or induce World War III.
I am not responsible for anything this code does to your computer, your family or your cat.
You are completely on your own. Any problems you encounter are probably new to me. (But go ahead and create [an issue][issues]).

[issues]: https://github.com/cfstras/cfmedias/issues
