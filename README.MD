# Twitch Clips

Create a list of Twitch clips on a given channel that can be imported in Excel or Google sheets.

## Example output

```csv
Title,Game,Date,Views,URL
super clip one,One Game,2023-12-24T21:52:23,5,https://clips.twitch.tv/url1-2zd1OiY4ZvfeP44Q
Clip Two,One Game,2023-12-22T21:42:23,7.https://clips.twitch.tv/url2-2zd1OiY4ZvfeP44Q
More clip,Other Game,2023-21-21T31:52:23,4,https://clips.twitch.tv/url3-2zd1OiY4ZvfeP44Q
```

## Releases 

You can download the last available release [here](https://github.com/juan-medina/twitch-clips/releases/latest)

We provide binary versions for Windows, Linux & MacOS.

Follow the [instructions](INSTRUCTIONS.txt) contained within the zip for your platform.

*Note: Your antivirus software may delete the application downloaded from this package, and is because executables need to be signed by a certification company to be trusted, unfortunately this as a hobby project that does not allow me to pay for a certification, but you have access to the full source code, you can run it yourself and even look at what actually does.*

## Run via makefile

### Run for a giving client id & channel
```sh
make run CLIENT_ID="xxxx" SECRET="yyyy" CHANNEL="zzzz" OUTPUT="clips.csv" DATE_FROM="2023-11-01" DATE_TO="2023-11-15" SORT_BY="DATE"
```
## Run the compiled binary

First build the binary with:

```sh
make build
```

### Run with an specified client id and channel

```sh
./build/twitch-clips -client_id xxxx -secret yyyy -channel zzzz -output clips.csv -date_from="2023-11-01" -date_to="2023-11-15" -sort_by=DATE
```

## Date From & To

These arguments are optional and they could define a window for the clips to be retrieved, however we can not bypass the way the twitch API handle these arguments.

If you specify a date from, but not a date to, it will return the clips from that date to two weeks after.

If you specify a date to, but not a date from, it will return all the clips without considering any date.

## Environment variables

Instead of passing the client id, secret and channel you could define the variables `TWITCH_CLIENT_ID`, `TWITCH_SECRET` and `TWITCH_CHANNEL`.

## Additional make commands

- clean : clean the build
- format : format the code
- vet : basic checks
- update : update dependencies & remove unused
- publish : creates a new release and publish it

## VS Code

Include a `launch.json` that allow to debug in VS Code.

*Note: You may need either to define environment variables in your system or add the arguments in the [launch.json](.vscode/launch.json) file*

## This software uses

- [Helix](https://github.com/nicklaw5/helix), a Twitch Helix API client written in Go (Golang).
- [Zero Allocation JSON Logger](https://github.com/rs/zerolog) a fast and simple logger dedicated to JSON output.

## LICENSE
```
MIT License

Copyright (c) 2023 Juan Medina

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
