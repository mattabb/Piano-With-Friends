# Piano With Friends
> Piano with Friends is a project which allows multiple people to connect and play piano together using websockets. 


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
Below is an image of the home screen that you are greeted to when entering the website
![Home Page](./screenshots/HomePage.png)

After you are connected, your webpage should change to show this:
![Piano Page](./screenshots/PianoPage.png)

## Technologies
### Deployment

#### Docker

We chose Docker because it allowed for a more streamlined deployment onto the server. It creates a lightweight, self sufficient container that can be transferred between environments. Docker is also one of the most popular container services so we thought learning how to use it for a project could be helpful.

The docker containers are then hosted in the Amazon Web Service, Elastic Container Registry (ECR). From the repositories, the front and back end are hosted in two separate instances of the Elastic Container Service (ECS)  and deployed with FARGATE. ECR is the simplest way we found to host docker containers on AWS and had seamless integration with ECS. FARGATE takes care of the provisioning and allocation of resources for the server. It helps us remain in the free tier and not pay for unnecessary servers and computation.
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

### TonalJS

Tonal is a music theory library. It contains functions to manipulate tonal elements of music (note, intervals, chords, scales, modes, keys). It deals with abstractions (not actual music or sound).  We use this mainly in our piano-mapping.js file

You can visit Tonal's github here:

https://github.com/tonaljs/tonal

### HowlerJS

Howler is a popular JavaScript audio library which we utilize to play the sound for our notes.  One of the key features that Howler includes is it's auto-caching.  This means that loaded sounds are automatically cached and re-used on subsequent calls which allows for better performance and bandwidth.  We use this in our Piano.vue file on key press and keep our audio files in the public folder in our UI.

You can visit Howler's homepage here:

https://howlerjs.com/


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
* Vue-CLI Front-end utilizing Vuetify, HowlerJS, TonalJS, websockets and more!
* Go back-end utilizing gorilla/mux and gorilla/websocket

To-do list:
* Improve the look of our website using Vue

## Status
Project is: _in progress_
