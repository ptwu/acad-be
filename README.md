# A Chengyu a Day Backend

[Website Link](https://www.achengyuaday.com)

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/peterwu)

Daily Chinese Chengyu idiom learning tool. Provides users with a custom link
intended as a homepage to learn and review idioms every day! Supports both
simplified and traditional characters.

This is the repo for the backend middleware. Made with Go (`gorilla-mux`), and
continuously deployed on Heroku. The database is a PostgreSQL server hosted on
AWS RDS. [Click here](https://github.com/ptwu/acad-fe) to see the frontend code.

## Getting Started

Clone the repository. Next, run `go run main.go` to start the server locally.
You will need a `.env` file containing a `POSTGRES_URL=` field of your own
PostgreSQL server.

## Contributing

Contribute potential issues / suggestions in the Issues section of this repo.
If you have any potential fixes, feel free to make a pull request.
