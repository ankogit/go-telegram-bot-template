# Golang Telegram bot template

This project serves as a template for quickly creating chatbots in telegram based [go-telegram-bot-api/telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api/).

## Installation

1) Install go mods
2) Create config.ini file in folder ./configs
3) Set your Telegram API KEY in telegramkey
4) Build and run cmd/main.go

## Usage

The template out of the box allows you to implement the following tasks:
- Create commands handles for telegram bot
- Create inline command handler
- Repository structure (ChatRepository implemented) base on boltdb
- Notification by CRON
- Http server for your API
- JWT authorization for api

## Change log

Please see the [changelog](changelog.md) for more information on what has changed recently.


## Contributing

ToDo:
- Tests coverage for repositories
- Tests coverage for handles
- Add new telegram-bot api functions

## Security

If you discover any security related issues, please email author email instead of using the issue tracker.

## Credits

- [All Contributors](link-contributors)

## License

license. Please see the [license file](LICENSE) for more information.

