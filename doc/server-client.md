Server/Client Modes
===================

cfmedias is built to be mainly used as a standalone client program, using it as a media player and management software.
However, one can also use it with a server instance to synchronize data between multiple instances or _thin_ instances, like a mobile phone.

Client
------

The normal mode of operation, most components are loaded. cfmedias acts as a complete media management suite. If configured, data is regularly synced to a server instance.
Also, a client instance can proxy data through to a thin client instance or raw file storage, while also synchronizing the music itself.

Thin Client
-----------

This can be an android application or similiar. It only stores a limited amount of database history to minimize file storage use. Additionally, it is optimized to conserve battery by only doing as little work as necessary; only when the device is currently awake  -- waking up a mobile device for timed tasks kills battery, especially if the work can be done once the user actually uses it.
When synchronizing data to a server or client instance, as much work as possible is handed off to the other instance for performance reasons.

Server
------

cfmedias can act as a standalone server to run without any GUI clients and playback capabilities. Such an instance can be used to maintain a library on a remote server (or a cloud server). A cfmedias server is the main hub for client devices to send listening statistics to (be it a cfmedias client, a mobile scrobbler or a third-party music player with a plugin).
Additionally, the database on a server can be synchronized with any other cfmedias instance.

Cloud
-----

A cloud instance is really just a server instance, but hosted for many users and set up with a focus on security concerns.
