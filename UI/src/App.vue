<template>
    <!-- 
      This is where we export our components, main.js imports this app component 
      and "mounts" it onto our DOM in public/index.html
    -->
  <div id="app">
    <img alt="Vue logo" src="./assets/logo.png" />
    <HelloWorld msg="Welcome to Your Vue.js App" />
    <!-- <PianoContainer/> -->

  </div>
</template>

<script>
import HelloWorld from "./components/HelloWorld.vue";
// import PianoContainer from "./components/PianoContainer";

export default {
  name: "App",
  components: {
    HelloWorld,
    // PianoContainer <- This is where we export "components" (think of them as custom-div elements) to use in the app template
  },
  data () {
    return {
      ws: null,
      serverUrl: "ws://localhost:8080/ws" // Idk what to put here lol
    }
  },
  // mounted is one of Vue's Virtual DOM's Lifecycles (similar to onpageload for regular js)
  mounted: function() {
    this.connectToWebsocket()
  },
  methods: {
    connectToWebsocket() {
      // connect to websocket at server url
      this.ws = new WebSocket( this.serverUrl );
      this.ws.addEventListener('open', (event) => { this.onWebsocketOpen(event)}
      )
    },
    onWebsocketOpen() {
      console.log("connected to ws!")
    }
  }

};
</script>

<style>
/* 
  Here is where we have our global styling, (fonts/main colors, etc)
  components will get their own specific styling in their file 
*/
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
