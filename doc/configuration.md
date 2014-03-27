Configuration
=============

cfmedias is configured via a JSON file located in the working directory. If none is found, it is created.

Options
-------

General options:

- _DbFile_: Path to primary SQLite database file.  
  default: `db.sqlite`

- _MediaPath_: Path to the main music folder.  
  default: `~/Music`

- _WebPort_: Port to the HTTP interface.  
  default: 38888

Media-specific options:

- _ListenedLowerThreshold_: The minimum amount of a track to be listened before it starts to get score for the occurrence.  
  default: 0.3

- _ListenedUpperThreshold_: The maximum amount of a track to be listened before it gets the full score for the occurence (considered to be _fully listened to_).  
  default: 0.7

  