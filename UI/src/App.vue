<template>
  <v-app>
    <v-app-bar app color="primary" dark>
      <div class="d-flex align-center">
        <v-img
          alt="Logo"
          class="shrink mr-2"
          contain
          :src="require('./assets/piano-logo.png')"
          transition="scale-transition"
          width="75"
        />

        <div
          alt="Name"
          class="title-font shrink mt-1 hidden-sm-and-down"
          min-width="100"
          width="100"
        >
        Piano With Friends
        </div>
      </div>

      <v-spacer></v-spacer>

      <!-- <v-btn
        href="https://github.com/vuetifyjs/vuetify/releases/latest"
        target="_blank"
        text
      >
        <span class="mr-2">Latest Release</span>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn> -->
    </v-app-bar>

    <v-main>
      <div v-if="isConnected == false" class="d-flex flex-column align-center">
        <v-card width="400">
          <v-row class="text-center">
            <v-col cols="12">
              <v-img
                :src="require('./assets/piano-logo.png')"
                class="my-3"
                contain
                height="150"
              />
            </v-col>

            <v-col class="mb-4">
              <h1 class="display-2 font-weight-bold mb-3">
                Piano With Friends
              </h1>
              <p>Enter your username</p>
              <v-form ref="usernameLogin" @submit="submit" onSubmit="return false;">
                <v-text-field
                  v-model="connection.username"
                  label="Username"
                ></v-text-field>
                <v-btn 
                  depressed 
                  ref="submitButton"
                  color="primary"
                  :loading = "loading"
                  :disabled = "loading"
                  @click = "loader = 'loading'"
                > 
                  Submit 
                </v-btn>
              </v-form>
            </v-col>
          </v-row>
        </v-card>
      </div>
      <div v-else>
        <PianoPage />
      </div>
    </v-main>
  </v-app>
</template>

<script>
import PianoPage from "./components/PianoPage";

export default {
  name: "App",

  components: {
    PianoPage
  },
  data: () => ({
    connection: {
      ws: null,
      username: ""
    },
    serverUrl: "localhost:8000",
    isConnected: false, 
    loader: null,
    loading: false
  }),
  mounted: function() {
    this.setUsername()
  },
  watch: {
    loader () {
      const l = this.loader
      this[l] = !this[l]
      this.setUsername()
      setTimeout(() => (this[l] = false), 3000)
    }
  },
  methods: {
    setUsername() {
        if (this.connection.username != "") {
          this.setWebsocketConnection()
          this.setIsConnected()
        }
    },
    setIsConnected() {
      if (this.isConnected == false) {
        this.isConnected = !this.isConnected
        console.log(this.isConnected)
      }
    },
    setWebsocketConnection() {
      // Ask for username and connect to websocket with it
      if (window["WebSocket"]) {
        const socketConnection = new WebSocket(
          "ws://" + this.serverUrl + "/ws/" + this.connection.username
        );
        console.log(socketConnection)
        this.connection.ws = socketConnection;
      }

      this.connection.ws.addEventListener("open", event => {
        this.onWebsocketOpen(event);
      });
    },
    onWebsocketOpen(event) {
      console.log(event, "connected to websocket!");
    }
  }
};
</script>

<style scoped>
title-font {
  font-family: "Lobster", cursive !important;
  font: "Lobster";
}
</style>
