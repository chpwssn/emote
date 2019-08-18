# emote

A federated emoticon service for collecting all of your custom Slack, Discord, etc. emoticons and discovering new ones.

This project is currently just a sandbox to play around with Go. It's not intended to be deployed as a production app and already has fun security vulnerabilities. Exploits are left as an exercise for the reader.

# emote server

The [emote server](server/README.md) stores an instance's emoticons and proxies requests to other emote servers.
