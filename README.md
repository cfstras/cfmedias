cfmedias
========
Making it easy to organize your media.

Overview
--------
cfmedias aims to let you take control of your media library.
It lets you manage all of your media in a central place.

What it should be one day:

- lightweight
- fast
- supporting every media format you could imagine
- networking-enabled to let you share media within your home network
  and keeping it in sync whilst letting you stream it freely
- able to import libraries from iTunes, WinAmp and others
- able to automatically add tags and pictures to your media
- able to make full or incremental backups
- highly portable

etc, etc...

Compiling
---------
You will need the Go package, git, mercurial, libportaudio, libsqlite3 and taglib.  
On OS X, you need at least Go 1.2rc5.

    git clone --recursive https://github.com/cfstras/cfmedias
    cd cfmedias
    make deps
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
You are completely on your own, I probably don't yet know myself how to solve any problems you might encounter
(If I did, I would have fixed it).
