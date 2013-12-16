Scrobbling
==========

To maximize the data output we can get from services like last.fm, cfmedias uses a modified version of the Simple Last.fm Scrobbler, which not only scrobbles listened tracks to last.fm but also submits more detailed info (_the raw data_) to the users cfmedias database.
These infos include any instance of a played track, with the addition of how far in the track was listened to and if it was skipped or the playlist stopped afterwards. We also include whether the track was scrobbled to last.fm afterwards to eliminate double playcounts.

These messages are aggregated at a cfmedias server, where they are input into the database and/or stored for a client to fetch.

Server
------

cfmedias can act as a standalone server to run without any GUI clients and playback capabilities. Such an instance can be used to maintain a library on a remote server (or a cloud server). A cfmedias server is the main hub for client devices to send listening statistics to (be it a cfmedias client, a mobile scrobbler or a third-party music player with a plugin).
Additionally, the messages on a server can be synchronized with any cfmedias instance to guarantee no false listening data is produced.
