# Multiplayer Piano

Multiplayer piano is a web application which uses a VueJS front-end and 
serves it using the Gorilla mux package in GoLang.

The VueJS is created using vue-cli-service, meaning there are some boiler-plates built in, such as ESLint, Jest testing (unit-testing which will be used later) and Babel (in case we have someone using IE lol).  While all of these files seem confusing in our UI folder, the only ones that matter right now are the public and src which contain our html and JS.
## Installation

Use the package manager npm to install multiplayer piano.

```bash
npm install
```

## Structure

In order to deploy our front-end we're able to either serve or build
```bash
npm run serve
npm run build
```

Server.go contains our main function in Golang which will be ran to run our server.

However, i'm looking to implement our websockets and back-end using go, and connect this to our front-end (I think this might be flawed?)


### Server.go
```

func main() {
  ... <- just log.println
  // Create our router instance using gorilla mux
  router := mux.NewRouter()
  // Call AddAppRoutes in routes.go
  AddAppRoutes(router)

  log.Fatal(http.ListenAndServe(":8080", router))
}

```

### routes.go
```
func AddAppRoutes() {
    // This is supposed to serve our JS files that are separated from the index.html
    // I think this is where my problem is somehow
    setStaticFolder(route)

	// Implement websockets and handlers (I don't think this is changing anything and separate
    // from the problem i'm having
	pool := handlers.NewPool()
	go pool.Run()
	log.Print("pool ran")


	// Set the default route to the index.html
	route.HandleFunc("/", handlers.RenderHome)
	
    // This is irrelevant and will be worked on after I get the front-end served
	// // Websocket handling
	// route.HandleFunc("/ws/{username}", func(responseWriter http.ResponseWriter, request *http.Request) {
		
	// 	var upgrader = websocket.Upgrader{
	// 		ReadBufferSize: 	1024,
	// 		WriteBufferSize: 	1024,
	// 	}

	// 	username := mux.Vars(request)["username"]

	// 	connection, err := upgrader.Upgrade(responseWriter, request, nil)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}

	// 	handlers.CreateNewSocketUser(pool, connection, username)
	// })

	log.Println("Routes loaded.")
}
```

### handlers/routes-handlers
This is where we call RenderHome from => This is just supposed to render the index.html file... I believe the problem is that the JS is separate from this file and so it's serving the index.html but not the JS, therefore the VueJS isn't getting mounted onto the app element... resulting in the blank page.  However i'm not sure how to fix this.

```
// RenderHome renders the home page
func RenderHome(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "UI/public/index.html")
}

```




