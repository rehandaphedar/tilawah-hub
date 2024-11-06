# Introduction

`tilawah-hub` is a forge for Qurʾān recitations. It aims to streamline collection of Qurʾān recitations, and thus, development of related software.

Unfortunately, there is no online demo at this point.

# Features

- CRUD on Users, Recitations, Recitation Files, Recitation Timings
- Automatically generate word level timings using [lafzize](https://sr.ht/~rehandaphedar/lafzize)

# Limitations/Upcoming Features

- Emulation of the [quran.com API](https://api-docs.quran.com/docs/category/quran.com-api) for audio endpoints
- Ability to make a recitation private
- Administration panel
- Docker + Compose deployment
- A setting to automatically lafzize on upload
- Automatic database migration
- Better input validation, for example, for verse keys

# Ecosystem

This project aims to ease the development of a range of software. Some examples include:

- Qurʾān applications can use the audio and timing files to enable anyone to listen to their own recitations with word by word highlighting.
- A program can be written to automate generation of Qurʾān videos. This would likely need a bit more configuration/modifications to the timings in case of long verses.
- [Tarteel AI](https://tarteel.ai) or similar software can be used to automatically analyse recitations for mistakes. Or, in reverse, user sessions could be exported from Tarteel AI to tilawah-hub so that all the benefits of tilawah-hub's standardisation can be reaped.
- A GUI timings editor would be greatly beneficial. [zonetecde/QuranCaption-2](https://github.com/zonetecde/QuranCaption-2) is very similar, however, a standardisation of the timings format is needed. Note that the word splitting is as per the quran.com API's `text_uthmani` field, and preserving it is important for reliability across programs.

Other ideas are welcome!

# API Reference

Check `docs/swagger.yaml`. Do note that authenticated requests require a `session_token` cookie. `swag` doesn't seem to support documenting this yet.

The audio files are present at `/uploads/{username}/{slug}/{verse_key}.mp3`.

The timings files are present at `/uploads/{username}/{slug}/{verse_key}.json`.

# Install Instructions

## Development Dependencies

- go
- [`sqlc`](https://sqlc.dev/)
- [`swag`](https://github.com/swaggo/swag)
- [`migrate`](https://github.com/golang-migrate/migrate)

## Runtime Dependencies

- `ffmpeg`
- [`migrate`](https://github.com/golang-migrate/migrate)

## Compiling

Clone the source code

``` shell
git clone https://git.sr.ht/~rehandaphedar/tilawah-hub
cd tilawah-hub
```

Generate SQL and build

``` shell
sqlc generate
go build .
```

Run database migrations (Make sure `data/` exists)

```shell
migrate -path internal/db/migrations -database sqlite3://data/db.sqlite up
```

Generate OpenAPI specificatipn

```shell
swag f
swag i
```

## Deploying

Copy `tilawah-hub`, `data/config.toml`, and `internal/db/migrations`.

Run database migrations (Make sure `data/` exists)

```shell
migrate -path internal/db/migrations -database sqlite3://data/db.sqlite up
```

Edit `config.toml` and run:

``` shell
./tilawah-hub
```
