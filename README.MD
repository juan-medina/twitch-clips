# Twitch Clips

List of Twitch clips on a given channel

## Run via makefile

### Run for a giving client id & channel
```sh
make run CLIENT_ID="xxxx" CHANNEL="yyyy"
```
## Run the compiled binary

First build the binary with

```sh
make build
```

### Run with an specified client id and channel
```sh
./build/twitch-clips -client_id xxxxx -channel yyyy
```

## Environment variables

Instead of passing the client id and channel you could define the variables `TWITCH_CLIENT_ID` and `TWITCH_CHANNEL`

## Additional make commands

- clean : clean the build
- format : format the code
- vet : basic checks
- update : update dependencies & remove unused

## VS Code

Include a `launch.json` that allow to debug in VS Code.

*Note: You may need either to define environment variables in your system or add the arguments in the launch.json file*

## LICENSE
```
Copyright (c) 2023 Juan Antonio Medina IglesiasMIT License

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
