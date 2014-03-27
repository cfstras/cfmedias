# API

Any instance of cfmedias exposes a REST API on the command line as well as via HTTP. If required, the HTTP api can be configured or turned off in the [configuration][config].

[config]: configuration.md

Every command has parameters and returns a JSON object in the following form:

_Sample output for the command `help c=help`_
```json
> help c=help

{
  "status": "OK",
  "result": {
    "help": {
      "MinAuthLevel": 0,
      "Aliases": [
        "help",
        "h",
        "?"
      ],
      "Desc": "Prints help.",
      "Args": {
        "c": "(Optional) the command to get help for"
      }
    }
  }
}
```

## Parameters

On the command line, parameters are given as simple `key=value` statements. Spaces can be inserted by escaping them with a backslash: `> list text=My\ Fair\ Lady`


### Authentication

On the HTTP side, most commands need authentication. This is done by passing an additional parameter `auth_token=xyz` to each request. This token can be acquired with the `login` command.

## HTTP URL

The command is given in the url.

_Sample HTTP output_
```bash
user@box ~ $ curl "http://localhost:38888/api/login?name=cfs&password=testtest"
{
  "status": "OK",
  "result": {
    "auth_token": "vNiTesQhgrvBQikM"
  }
}
```

## Audioscrobbler

Built into the API is an implementation of the [Audioscrobbler protocol v1.2.1][as_api].
It can be used to [scrobble][scrobbling] into the cfmedias database right from your favourite scrobbler application. Simply use `http://cfmedias-server-ip/api/audioscrobbler/?hs=true` as the handshake URL in your client, your normal username as login and the authentication token you got from logging in as the password.

It should be noted that the Audioscrobbler protocol lacks some of the information -- primarily the actual amount of time the user listened to the track before possibly skipping it -- and is not the most secure. However, it enables a variety of clients to submit data to cfmedias in a simple manner.

[scrobbling]: scrobbling.md
[as_api]: http://www.last.fm/api/submissions
