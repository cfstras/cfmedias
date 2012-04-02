Components of the program
====================

- core:
  + manages memory
  + holds information about current media and following media
  + loads components
  + (offers plugin-api?)
  
- cli:
  + is just a set of methods
  + with a parser
  
- gui:
  + displays information
  + builds upon cli methods
  + is fast, responsive, clean and easy-to-use
  + can be disabled at any time
  
- db:
  + delivers content to core
  + manages lists and users
  + keeps media hashes
  + is replaceable (network, sqlite, mysql, postgres, ...)
  
- db-importer:
  + uses db to add media to the library
  + reads files, itunes, winamp, m3u, upnp, ftp and categorizes them
  + watches folders/libraries
  
- player:
  + plays media using ffmpeg/libvlc/mplayer
  + integrates video to gui, if necessary
  + reports progress to core

- info:
  + reads tags of media / software
  + used by db-importer
  + works headless and without player
  
- share:
  + shares media in bittorrent-manner
  + keeps list of friends
  + also headless