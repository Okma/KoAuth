# Two Factor Authentication

## General Idea
Server generates key for user. Key is sent to user. Key expires after X amount of time.

### Idea 1 - Text/SMS API
Server generates and texts user a validation code.

### Idea 2 - Mobile App
Server generates code; sends to user via HTTP to client-side app.
Client-side app pairs with account via some serialization code bound to account.

## Server 
Written in Go.

## Client
Hybrid mobile application written with Ionic/Framework7. Packaged via Phonegap/Cordova?
