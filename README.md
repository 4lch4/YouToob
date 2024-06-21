# YouToob

This repo is home to a small API, written in Go, that provides a simple way to retrieve the following data about a YouTube channel:

- The latest video published.
  - `/:channelName/video`
- The latest short published.
  - `/:channelName/short`
- The latest VOD published.
  - `/:channelName/vod`
- The next upcoming livestream.
  - `/:channelName/live`

## Tech Stack

- [Go][0]
- [Gin][1]
- [Docker][2]
- [Axiom Logging][3]

## To Do

- [ ] Add a `Dockerfile` to the project.
  - [ ] While I love building a single executable, I'd much rather have a container that I can deploy instead.
- [ ] Add [support for Turso][4] to store some response data.
  - [ ] Should help prevent quota/rate limit issues (if we ever run in to those) by limiting calls to the YouTube API.
  - [ ] For example, the following bits of data should prove useful:
    - [ ] The channel name -> channel ID mapping.
    - [ ] The channel name/ID -> playlist ID(s) mapping.
- [ ] Add tests.
  - [ ] Mostly for the experience of how to do this because I _really_ hate writing tests...

[0]: https://go.dev
[1]: https://gin-gonic.com
[2]: https://www.docker.com
[3]: https://axiom.co
[4]: https://docs.turso.tech/sdk/go/quickstart
