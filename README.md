# Overview
"discord-wordle-bot" is a simple, self-hosted bot for playing **Wordle** in your server!
- Random answers for every puzzle
- (Planned) support for custom word length
- Written in [Go](https://go.dev/) with [DiscordGo](https://github.com/bwmarrin/discordgo)

# Setup
If you haven't already installed the Go language, you may find resources for doing so [here](ihttps://go.dev/doc/install).

To host this bot yourself, clone the repository and create an application on the Discord Developer Portal. You'll need the new bot's token, which you should paste into a file called `.env`.
```.env
export BOT_TOKEN=EXAMPLEOFA.TOKEN987654323
```

# Acknowledgements
discord-wordle-bot relies on a few external packages. Go check them out!
- [DiscordGo](https://github.com/bwmarrin/discordgo)
- [GoDotEnv](https://github.com/joho/godotenv)
- [emoji](https://github.com/joho/godotenv)

# TODOs
- Slash commands
- Custom word length

# License
This project is freely available under the MIT license.
