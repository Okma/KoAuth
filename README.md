# KoAuth
KoAuth is a robust, embeddable two-factor authentication(TFA) framework.

## Features
* Server to handle HTTP requests from mobile and web client.
* Web client demo.
* Android client for pairing with account.
* Supports MongoDB.

## Install
To run and/or modify the Go server, you'll need [Go](https://golang.org/).

The server requires a database config [TOML](https://github.com/toml-lang/toml) file in the same directory as the server. This file should contain a string field (named as `DBConnectionURI`) containg the URI that allows the MongoDB driver to connect to a MongoDB database.

The server also contains a config TOML, which allows you to customize:
* Key life
* Generated key size 
* Generated serial size

To run and/or modify the Android client code, you'll need [Android Studio](https://developer.android.com/studio/index.html).

## Purpose
KoAuth is used to provide simple TFA, adding an additional layer of security to any authentication process.

The authentication process effectively binds a user's email to their mobile device. As a result, KoAuth can be utilized across multiple services, assuming the user uses the same email.

## Server 
The server handles incoming HTTP requests, as well as serial and key generation for clients.
Routes include:
* Registering a new user.
* Checking if user exists. (If KoAuth is enabled for user.)
* Serial validation from mobile input.
* Key validation from web input.
* Returning a randomly generating key to mobile client.

## Client
Web client demo is provided. 

The client code is used to facilitate user registration with KoAuth using an identifier (uses email currently), as well as consume and validate key input required for authentication.

## Mobile Client (Android)
When first launching the mobile application, the user will be required to add a serial (requested and displayed by the web client) to their configuration. 

If validated, the device will be bound to the given serial in a database.

The mobile client will routinely query the server at an arbitrary interval for their current key code and display it for the user on the main page.

## Mobile Client (iOS) 
Not currently available.

## License
KoAuth is licensed under the [MIT license](https://en.wikipedia.org/wiki/MIT_License).
