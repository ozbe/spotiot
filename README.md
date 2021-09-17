# spotiot

Spotify shuffle AWS IoT Button.

## Setup

### Environment Variables
```bash
$ cat > .env << EOF
SPOTIFY_ID=<CLIENT_ID>
SPOTIFY_SECRET=<CLIENT_SECRET>
ASSERT_USER_ID=<OPTIONAL_SPOTIFY_USER_ID>
EOF
```

### AWS

- [Setup AWS CLI V2](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
- [Setup Go Lambda](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html)
- [Setup AWS IoT Button](https://aws.amazon.com/iotbutton/)

## Build

```bash
$ make build
```

## Deploy

```bash
$ FUNC_NAME=<AWS_LAMBDA_FUNCTION_NAME> make deploy
```

After deploying, configure your AWS IoT Button to trigger `<AWS_LAMBDA_FUNCTION_NAME>`

## Usage

Start playing Spotify on a device and then use the remote to control Spotify:
- Single click to play/pause
- Double click to skip to the next track
- Long click to toggle shuffle