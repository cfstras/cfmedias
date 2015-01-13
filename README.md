cfmedias [![Build Status](https://travis-ci.org/cfstras/cfmedias.svg?branch=master)](https://travis-ci.org/cfstras/cfmedias)
========
The media player and manager you need.

Overview
--------
cfmedias aims to let you take control of your media library.  

Priorities:

- fast
- support all common media formats, both lossless and lossy
- iPod, Android and mp4-player sync
- incremental import, export databases from iTunes, Foobar2k, Google Play Music
- metadata, rating and play/skip-count import from all of the above, plus Last.fm
- highly portable, full or incremental backups
- play behaviour analysis and automatic playlisting
- both easy and advanced filter playlists

Status
------
_I'll try and keep this up to date._

What works:

- single-binary with web server & assets
- media scanner for music
- iPod sync, at least for titles
- HTTP and commandline API
- crude web-interface with no features
- sqlite database (subject to a lot of change)
- audioscrobbler server endpoint

Compiling
---------
You will need the Go package, git, mercurial, make and npm.  
Also, you will need development headers for libportaudio, libsqlite3, taglib and
libgpod.  
On OS X, you need at least Go 1.2rc5.

    git clone --recursive https://github.com/cfstras/cfmedias
    cd cfmedias
    make run


License
-------
This software is released under the 2-clause BSD-license. For details, see LICENSE.md
Also, the author would love it if you send him any useful modifications you made to this code.


Disclaimer
----------
*This program is far from finished.*
It may destruct your computer, your whole network or induce World War III.
I am not responsible for anything this code does on your computer, your family or your cat.
You are completely on your own, I probably don't yet know myself how to solve any problems you might encounter (If I did, I would have fixed them).
