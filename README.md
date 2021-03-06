# ![Trinket Logo](https://github.com/thee-engineer/trinket/blob/master/trinket.png?raw=true)

[![Go Report Card](https://goreportcard.com/badge/alexandru.cc/go/trinket)](https://goreportcard.com/report/alexandru.cc/go/trinket)

Is a server mimic for COMP28112 exercise 2, which gives you the ability to run it locally and configure it as you wish. No more annoying errors while debugging, no more undefined/unexpected responses. You now have control over the system and can test all edge cases of your program.

## Why ?

* Are you a second year CS student at The University of Manchester?
* You looked over your choices and thought to yourself, "Wow, distributed computing sounds awesome!"?
* You started working on your second lab for COMP28112?
* Did you get tired with the server responding with **404**, **402**, **503**, **400**, Error 482: Somebody shot the server with a 12-gauge. Please contact your administrator?
* You find the documentation/labscript old?

If the answer to any of the above is YES, then you are in the right place!

While the second lab for COMP28112 is one of my favorite labs (because it attempts to recreate a real-world distributed system environment with all of it's flaws and annoyances), but having to test your bot on a server which can respond with an error code 10+ times in a row just because it wants to can get infuriating really fast. During my year, it wasn't that bad, I caught on early and developed my bot to handle *anything* you throw at it!

But seeing the current second year students struggle with the server (to the point where they doubt the correctness of their programs) gave me the great idea of "reverse engineering"/mimicking the "protocol" as used by the course BUT with a few improvements.

## Features

* Run the server locally (no need to be in Kilburn or to use the VPN)
* Easy to configure
  * Disable errors
  * Set/Clear/Block slots
  * View requests as they go in the backend

## Install

* If you have [Go](https://golang.org) installed use the following command in your terminal. You can download the latest version of [Go from here](https://golang.org/dl/).
  * `go get alexandru.cc/go/trinket`
* The other option is to download the binary from the [GitHub Releases](https://github.com/thee-engineer/trinket/releases) page
  * [macOS](https://github.com/thee-engineer/trinket/releases/download/v0.1/trinket-macOS)
  * [Linux x86 (32 bit)](https://github.com/thee-engineer/trinket/releases/download/v0.1/trinket-linux-x86)
  * [Linux x64 (64 bit)](https://github.com/thee-engineer/trinket/releases/download/v0.1/trinket-linux-x64)
* Obtain the source code and compile it using Go
  * `git clone https://github.com/thee-engineer/trinket.git`
  * `cd trinket`
  * `go build`

## Usage

Note that in order to simulate the COMP28112 server for Exercise 2 you will
have to start two instances of `trinket`, one for simulating the hotel and
one for the band.

```shell
# Usage
trinket <PORT> <SLOTS> <USER LIST> <PASS LIST> <MAX BOOKINGS> <FAIL CHANCE %>

# Examples

# Create a server, listening on port 3000, with 200 slots a single user john
# with the password doe and a maximum number of 2 booked slots per user
# with a 43% chance fail for each PUT/GET request on queue
trinket 3010 200 "john" "doe" 2 43

# Create a server with multiple users
trinket 3010 200 "john mike paul" "doe pass test" 2 43
```

## TODO

- [ ] Create Wiki detailing the protocol, expected & mimic behavior
- [x] Finish stateless server mimic
  - [x] Receive requests PUT
  - [x] Enqueue requests
  - [x] Parse listing queue
  - [x] Parse individual request
  - [x] Handle URI GETs
  - [x] Return responses
- [x] Create user views
- [ ] Create binary releases (macOS, Linux)
- [ ] Write contribution guide

## Extra

I will not be working on these, if you feel generous you are more than welcome
to contribute to the project.

- [ ] Fancy user views with Bootstrap
- [ ] Implement stateful server, some database (SQL)
- [ ] Codebase cleanup and organization
