# Piano With Friends
> Piano with Friends is a project which allows multiple people to connect and play piano together using websockets.  You can try it out for yourself at:

> http://thisiswebsiteurl.com

## Table of contents
* [General info](#general-info)
* [Screenshots](#screenshots)
* [Technologies](#technologies)
* [Setup](#setup)
* [Features](#features)
* [Status](#status)

## General info
### Why?
The idea of Piano With Friends came about after COVID-19 led to a prompt halt of playing musical instruments with friends.

### What?
Piano With Friends contains two separate deployments using Docker; the back-end API built using GoLang and the front-end built using Vue-CLI 3.  This allows for more abstraction and separation between client and server.  

The GoLang back-end is structured around the use of Websockets and utilizes the Gorilla Mux router, as well as Gorilla Websocket.

The front-end, built using Vue-CLI 3, uses the Vue library Vuetify and incorporates unit testing, linting, and Babel for backwards browser compatibility.

Both the client and server utilize Docker and are in their own separate containers. 


## Screenshots
![Example screenshot](./img/screenshot.png)

## Technologies
### Deployment

#### Docker

Enter docker description here
### Front-end

Vue is a new, sleek JavaScript library that is closely related to the rising ReactJS, but different in a few ways, with the major difference being the extra abstraction of Vue's components.  The choice to use Vue was made due to its extensive documentation, allowing for a less steep learning curve.  Furthermore, the Vue-CLI 3 tooling was used, allowing for a quick out-of-the-box front-end setup with multiple libraries that are easy to add to.

You can view the Vue-CLI 3 homepage here: https://cli.vuejs.org/


#### Vuetify
Vuetify is a Vue UI library which contains multiple components that make the development much less strenuous.  One of the tools used throughout the project is the Vuetify flexbox style, which mimics the growing flexbox style in the HTML community.  Furthermore, the components that come with Vuetify are used to build all of the front-end.

You can visit Vuetify's homepage here: 

https://vuetifyjs.com/en/


#### ESLint
ESLint, a popular plugin linter is used in this project.  ESLint statically analyzes your code to quickly find problems. Many problems ESLint finds can be automatically fixed. ESLint fixes are syntax-aware so you won't experience errors introduced by traditional find-and-replace algorithms.

You can visit ESLint's homepage here:

https://eslint.org/

#### Jest
Jest is one of the major JavaScript unit-testing libraries that has a keen focus on simplicity and readability.  It allows you to make tests which keep track of large objects with ease. Snapshots live either alongside your tests, or embedded inline.

You can visit Jest's homepage here:

https://jestjs.io/


### Back-end 

The backend for this project is built using the language: GoLang.  GoLang is an open source programming language built by a team of Google developers that were focused on readability, reliability and efficiency.  Go contains a couple of cool features such as concurrency, interfaces and a plethora of built-in packages.

You can visit Go's homepage here:

https://golang.org/


#### gorilla/mux

gorila/mux is a powerful HTTP router and URL matcher for building Go web servers with 🦍.  The name mux stands for "HTTP request multiplexer". Like the standard http.ServeMux, mux.Router matches incoming requests against a list of registered routes and calls a handler for the route that matches the URL or other conditions. 

You can visit gorilla/mux's github page here:

https://github.com/gorilla/mux

#### gorilla/websocket
Gorilla WebSocket is a Go implementation of the WebSocket protocol.  WebSockets are a next-generation bidirectional communication technology for web applications which operate over a single socket and are exposed via a JavaScript interface in HTML 5 compliant browsers. The socket starts out as a HTTP connection and then "Upgrades" to a TCP socket after a HTTP handshake. After the handshake, either side can send data.

You can read more about websockets here:


https://www.tutorialspoint.com/html5/html5_websocket.htm


You can visit gorilla/websocket's github page here:

https://github.com/gorilla/websocket/

## Setup


### Docker



### Independent Vue Front-End

Navigate to the UI directory


First download all the npm dependencies using:
```bash
$npm install
```

To run the development server
```bash 
$ npm run serve
```

To build the project into production
```bash 
$ npm run build
```
To run a unit-test using Jest
```bash 
$ npm run test:unit 
```

To run the ESLinter 
```bash 
$ npm run lint
```

### Independent Go Back-End
Navigate to the backend directory

To Run the server
```bash 
$ go run ./
```

To build for production
```bash
$ go build ./
```

Run the exe that is generated from the above command.


## Code Examples
Show examples of usage:
`put-your-code-here`

## Features
* Vue-CLI Front-end utilizing Vuetify and websockets
* Go back-end utilizing gorilla/mux and gorilla/websocket

To-do list:
* Wow improvement to be done 1
* Wow improvement to be done 2

## Status
Project is: _in progress_
