Scrobbling
==========

To maximize the data output we can get from services like last.fm, cfmedias uses a modified version of the Simple Last.fm Scrobbler, which not only scrobbles listened tracks to last.fm but also submits more detailed info (_the raw data_) to the users cfmedias database.
These infos include any instance of a played track, with the addition of how far in the track was listened to and if it was skipped or the playlist stopped afterwards. We also include whether the track was scrobbled to last.fm afterwards to eliminate double playcounts.

These messages are aggregated at a cfmedias server, where they are input into the database and/or stored for a client to fetch.
