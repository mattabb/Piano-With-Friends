<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <div class="d-flex align-center">
        <v-img
          alt="Vuetify Logo"
          class="shrink mr-2"
          contain
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-logo-dark.png"
          transition="scale-transition"
          width="40"
        />

        <v-img
          alt="Vuetify Name"
          class="shrink mt-1 hidden-sm-and-down"
          contain
          min-width="100"
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-name-dark.png"
          width="100"
        />
      </div>

      <v-spacer></v-spacer>

      <v-btn
        href="https://github.com/vuetifyjs/vuetify/releases/latest"
        target="_blank"
        text
      >
        <span class="mr-2">Latest Release</span>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>
    </v-app-bar>

    <v-main>
      <HelloWorld />
    </v-main>
  </v-app>
</template>

<script>
import HelloWorld from "./components/HelloWorld";

export default {
  name: "App",

  components: {
    HelloWorld
  },
  data: () => ({
    ws: null,
    serverUrl: "localhost:8000"
  }),
  mounted: function() {
    this.setWebsocketConnection()
  },
  methods: {
    setWebsocketConnection() {
      // Ask for username and connect to websocket with it
      const username = prompt("Enter username")
      if (window["WebSocket"]) {
        const socketConnection = new WebSocket("ws://" + this.serverUrl + "/ws/" + username)
        this.ws = socketConnection
      }

      this.ws.addEventListener("open", event => {this.onWebsocketOpen(event)})
    },
    onWebsocketOpen(event) {
      console.log(event, "connected to websocket!")
    }
  }
};
</script>
