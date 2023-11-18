# Simple Messenger Go

Go server for [Simple Messenger](https://github.com/anthonydip/flutter-messenger-go) Flutter application.

## Introduction

Simple Messenger Go is a Go-based server designed to power the backend for the "Simple Messenger" Flutter application. In addition to its role as the backend API server to handle HTTP requests, it also implements WebSockets to enable real-time messaging between users of the Flutter application.

## Features
- **User Authentication**: Supports both traditional email-password authentication, along with Google sign-in
- **WebSockets for Real-time Messaging**: Using Gorilla WebSocket, it enables direct messaging between users, offering a real-time, bi-directional communication channel
- **Token Verification**: Implements token-based verification using JWT to ensure security and integrity of user sessions and information
- **Firestore Database Integration**: Integrates with Firestore, a flexible and scalable NoSQL cloud database. Utilizing Firestore, it allows storage of user data, tokens, and other relevant information.

## Directory Structure
```bash
├───app
│   └───storefront-api
│       ├───middleware
│       ├───routes
│       │   ├───auth
│       │   │   ├───signin
│       │   │   └───tokens
│       │   │       └───access
│       │   └───users
│       │       └───friends
│       ├───utils
│       ├───webserver
│       │   └───mock
│       └───ws
├───internal
│   └───storefront
│       └───mock
├───keys
└───pkg
    ├───authentication
    │   └───mock
    └───dtos
```

- **app/**: Entry point for the Go server and is where the main.go lives. This is were all the HTTP pipeline is built, along with its implementation details. WebSocket management through the client and hub is also done here.
- **internal/**: This is where all the domain logic goes, along with any Firestore data queries.
- **keys/**: Holds the various private and public keys used to sign, verify and issue JSON Web Tokens.
- **pkg/**: Holds data transfer objects, which allows structs to be designed for sharing data between packages and encoding/trasmitting over the wire as JSON. Any authentication functions and protocols are handled here as well.