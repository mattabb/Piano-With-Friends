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
      <div v-if="!isConnected" class="d-flex flex-column align-center">
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
              <v-form
                ref="usernameLogin"
                @submit="submit"
                onSubmit="return false;"
              >
                <v-text-field
                  v-model="connection.username"
                  label="Username"
                ></v-text-field>
                <v-btn
                  depressed
                  ref="submitButton"
                  color="primary"
                  :loading="loading"
                  :disabled="loading"
                  @click="
                    loader = 'loading';
                    setUsername();
                  "
                >
                  Submit
                </v-btn>
              </v-form>
            </v-col>
          </v-row>
        </v-card>
        <v-alert v-if="connectionErr" dense outlined type="error">
          Error when connecting to the server, please try again later.
        </v-alert>
        <v-alert v-if="invalidUsername" dense outlined type="error">
          Please enter a valid username, empty usernames are not accepted.
        </v-alert>
      </div>
      <div v-else>
        <Piano :octave-start="1" :octave-end="6" />
      </div>
    </v-main>
  </v-app>
</template>

<script>
//import PianoPage from "./components/PianoPage";
import Piano from "./components/Piano";

export default {
  name: "App",

  components: {
    //PianoPage,
    Piano
  },

  data: () => ({
    connection: {
      ws: null,
      username: ""
    },
    serverUrl: "localhost:8000",
    isConnected: false,
    connectionError: false,
    loader: null,
    loading: false
  }),

  watch: {
    loader() {
      const l = this.loader;
      this[l] = !this[l];
      this.loader = null;
      setTimeout(() => (this[l] = false), 3000);
    }
  },

  computed: {
    invalidUsername() {
      return this.connection.username == "";
    },

    connectionErr() {
      return this.connectionError && this.connection.username != "";
    }
  },

  methods: {
    // This is just so when we hit "enter" on the form, the page doesn't reload
    submit() {
      return false;
    },

    // this is called on mount => sets connection and then is connected so we know to move to piano page
    setUsername() {
      if (this.connection.username != "") {
        if (this.connectionError) {
          this.setConnectionError();
        }
        this.setWebsocketConnection();
      }
    },

    setConnectionError() {
      this.connectionError = !this.connectionError;
    },

    setIsConnected() {
      this.isConnected = !this.isConnected;
    },

    createCompleteUsername(username) {
      if (!this.verifyValidUser(username)) {
        for (var i = 0; i < 5; i++) {
          var randomNum = Math.floor(Math.random() * 10 + 1);
          username = username + randomNum;
        }
        this.connection.username = username;
      }
    },

    verifyValidUser(username) {
      var lastFive = username.substr(username.length - 5);
      var lastFiveInt = Number(lastFive);
      console.log(lastFiveInt);
      if (!Number.isInteger(lastFiveInt)) {
        return false;
      } else {
        return true;
      }
    },

    /* All Websocket Methods are below: */

    setWebsocketConnection() {
      this.createCompleteUsername(this.connection.username);

      if (window["WebSocket"]) {
        const socketConnection = new WebSocket(
          "ws://" + this.serverUrl + "/ws/" + this.connection.username
        );
        this.connection.ws = socketConnection;
      }

      this.setWebsocketErrorListener();

      this.setWebsocketOpenListener();
    },

    setWebsocketErrorListener() {
      this.connection.ws.addEventListener("error", event => {
        console.log("Error connection:", event);
        this.setConnectionError();
        if (this.isConnected) {
          this.setIsConnected();
        }
      });
    },

    setWebsocketOpenListener() {
      this.connection.ws.addEventListener("open", event => {
        this.onWebsocketOpen(event);
      });
    },

    setWebsocketCloseListener() {
      this.connection.ws.onclose = err => {
        console.log("Your connection is closed");
        console.log(err);
      };
    },

    setWebsocketMessageListener() {
      this.connection.ws.onmessage = messageEvent => {
        const socketPayload = JSON.parse(messageEvent.data);
        console.log("received message from backend...", socketPayload);

        switch (socketPayload.EventName) {
          // Join case
          case "join": {
            this.checkIfValidPayload(socketPayload);

            console.log(socketPayload.eventPayload, "has joined the chat");
            break;
          }
          case "disconnect": {
            this.checkIfValidPayload(socketPayload);

            console.log(socketPayload.eventPayload, "has left the chat");
            break;
          }
          case "keyboardPress": {
            this.checkIfValidPayload(socketPayload);

            const messageContent = socketPayload.EventPayload;
            const sentBy = messageContent.username;
            const actualMessage = messageContent.message;

            console.log({
              messageFrom: sentBy,
              message: actualMessage
            });
            break;
          }
          default: {
            break;
          }
        }
      };
    },

    onWebsocketOpen(event) {
      console.log(event, "connected to websocket!");

      this.setIsConnected();

      this.setWebsocketCloseListener();

      this.listenToWebsocketMessage();

      var mockMessage = {
        eventName: "keyboardPress",
        message: "abcdefgh"
      };
      this.sendWebsocketMessage(mockMessage);
    },

    // this is how we send messages to the backend
    sendWebsocketMessage(socketPayload) {
      console.log("message being sent", socketPayload);
      this.connection.ws.send(
        JSON.stringify({
          EventName: socketPayload.eventName,
          EventPayload: {
            username: this.connection.username,
            message: socketPayload.message
          }
        })
      );
    },

    // Make sure payload is not empty
    checkIfValidPayload(socketPayload) {
      if (!socketPayload.EventPayload) {
        return;
      }
    },

    // listen to response from the websocket
    listenToWebsocketMessage() {
      // If we have no connection, we can't listen
      if (this.connection.ws === null) {
        console.log("hit error in listen to websocket message");
        return;
      }

      this.setWebsocketMessageListener();
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
