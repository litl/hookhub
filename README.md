# HookHub - A simple webhook handler for Github Releases

## Installation

Make sure you have your Go environment setup, in particular that GOPATH is set.
Then run:

  $ go install github.com/litl/hookhub

## Configuration

You will need a hookhub.toml configuration file in the server's working
directory. An example is provided.

You will also need to configure your Github webhooks to point at
http://\<host\>:\<port\>/github_webhook, and then modify them to listen
for Releases.

First, find the id for your hook:

    curl -H "Authorization: token YOURTOKEN" https://api.github.com/repos/:owner/:repo/hooks

Next, PATCH in support for Release events:

	curl -X PATCH -H "Authorization: token YOURTOKEN" -H "Content-Type: application/json" -d '{"add_events":["release"]}' https://api.github.com/repos/:owner/:repo/hooks/:id

## Copyright and License

HookHub is Copyright (c) 2013 litl, LLC and licensed under the MIT license.
See the LICENSE file for full details.

Heavily derived from https://github.com/litl/hookyapp :-)
