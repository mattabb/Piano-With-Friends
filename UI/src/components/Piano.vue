<template>
  <div class="keyboard" :style="style">
    <v-btn
      depressed
      ref="recordButton"
      color="primary"
      @click="
        loader = 'loading';
        recordStart();
      "
    >
      Start Recording
    </v-btn>
    <v-btn
      depressed
      ref="recordButton"
      color="primary"
      @click="
        loader = 'loading';
        recordStop();
      "
    >
      Stop Recording
    </v-btn>
    <v-btn
      depressed
      ref="recordButton"
      color="primary"
      @click="
        loader = 'loading';
        recordPlay();
      "
    >
      Play Recorded Music
    </v-btn>
    <Keypress key-event="keyup" @any="keyUpMonitor" />
    <Keypress key-event="keydown" @any="keyDownMonitor" />
    <ul>
      <li
        v-for="(key, index) in keys"
        :key="index"
        :style="key.style"
        @mousedown="toggleTrue(key.name)"
        @mouseup="toggleFalse(key.name)"
        :class="[...key.class]"
      >
        <span>{{ key.name }}</span>
      </li>
    </ul>
  </div>
</template>

<script>
// import pianoState from "../library/piano-state";
import { addKeyCodeToKeys } from "../library/piano-mappings";
import { Howl } from "howler";

const WHITE_KEYS = ["C", "D", "E", "F", "G", "A", "B"];
const BLACK_KEYS = ["C#", "D#", null, "F#", "G#", "A#", null];
const MIN_OCTAVE = 0;
const MAX_OCTAVE = 7;
const MIN_NOTE = 0;
const MAX_NOTE = 6;
const WHITE_KEYS_PER_OCT = 7;
const BLACK_KEYS_PER_OCT = 5;

export default {
  name: "Piano",

  components: {
    Keypress: () => import("vue-keypress")
  },

  // Props are basically parameters for vue components
  props: {
    octaveStart: {
      type: Number,
      validator(value) {
        return value >= MIN_OCTAVE && value <= MAX_OCTAVE;
      },
      default() {
        return MIN_OCTAVE;
      }
    },

    octaveEnd: {
      type: Number,
      validator(value) {
        return value >= MIN_OCTAVE && value <= MAX_OCTAVE;
      },
      default() {
        return MAX_OCTAVE;
      }
    },

    noteStart: {
      type: [Number, String],
      validator(value) {
        if (typeof value === "string") {
          return WHITE_KEYS.includes(value);
        } else {
          return value >= MIN_NOTE && value <= MAX_NOTE;
        }
      },
      default() {
        return WHITE_KEYS.indexOf("A");
      }
    },

    noteEnd: {
      type: [Number, String],
      validator(value) {
        if (typeof value === "string") {
          return WHITE_KEYS.includes(value);
        } else {
          return value >= MIN_NOTE && value <= MAX_NOTE;
        }
      },
      default() {
        return WHITE_KEYS.indexOf("C");
      }
    },
    // websocket connection
    connection: {
      type: Object
    }
  },
  // Our data variables... think of them as this component's global variables (ONLY FOR THIS COMPONENT)
  data: () => ({
    offsets: {
      octaveStart: 0,
      octaveEnd: 3,
      noteStart: 0,
      noteEnd: 0
    },
    keysData: [],
    conn: {
      ws: null,
      username: ""
    },
    keyPressedName: ""
    // ,pianoState: []
  }),

  // "Created" Vue Lifecycle Hook
  // for more information, see: https://www.digitalocean.com/community/tutorials/vuejs-component-lifecycle
  created() {
    // Set all of our data variables according to props

    if (typeof this.noteStart === "string") {
      this.offsets.noteStart = WHITE_KEYS.indexOf(this.noteStart);
    } else {
      this.offsets.noteStart = this.noteStart;
    }

    if (typeof this.noteEnd === "string") {
      this.offsets.noteEnd = WHITE_KEYS.indexOf(this.noteEnd);
    } else {
      this.offsets.noteEnd = this.noteEnd;
    }

    this.offsets.octaveStart = this.octaveStart;
    this.offsets.octaveEnd = this.octaveEnd + 1;

    if (
      this.offsets.octaveStart > this.offsets.octaveEnd ||
      (this.offsets.octaveStart === this.offsets.octaveEnd &&
        this.offsets.noteStart > this.offsets.noteEnd)
    ) {
      throw new Error(
        "The start octave must be lower than or equal to the end octave and the start note must be lower than the end note."
      );
    }

    this.conn = this.connection;
    console.log(this.conn);

    this.setWhiteKeys(this.keysData);
    this.setBlackKeys(this.keysData);
    addKeyCodeToKeys(this.keysData);

    // this.pianoState = pianoState
  },

  mounted() {
    this.setWebsocketMessageListener();
  },

  computed: {
    offsetStart() {
      return this.clamp(this.offsets.noteStart, MIN_NOTE, MAX_NOTE);
    },

    offsetEnd() {
      return this.clamp(this.offsets.noteEnd, MIN_NOTE, MAX_NOTE);
    },

    totalWhiteKeys() {
      return Math.min(
        Infinity,
        this.numOctaves * WHITE_KEYS_PER_OCT -
          this.offsetStart -
          (WHITE_KEYS_PER_OCT - this.offsetEnd + 1) +
          2
      );
    },

    totalBlackKeys() {
      return Math.min(
        Infinity,
        this.numOctaves * BLACK_KEYS_PER_OCT -
          this.offsetStart -
          (BLACK_KEYS_PER_OCT - this.offsetEnd + 1)
      );
    },

    totalKeys() {
      return this.totalWhiteKeys + this.totalBlackKeys;
    },

    numOctaves() {
      return (
        1 +
        (Math.min(MAX_OCTAVE, this.offsets.octaveEnd) -
          Math.max(MIN_OCTAVE, this.offsets.octaveStart))
      );
    },

    style() {
      return {
        "--keys": this.totalWhiteKeys,
        "--octaves": this.numOctaves
      };
    },

    keys() {
      const keys = [];

      this.setWhiteKeys(keys);
      this.setBlackKeys(keys);

      addKeyCodeToKeys(keys);

      return keys;
    }
  },
  methods: {
    recordStart() {
      var time = this.getCurrentTime();
      var socketPayload = {
        EventName: "recordStart",
        EventPayload: {
          username: this.conn.username,
          message: "",
          time: time
        }
      };
      this.sendWebsocketMessage(socketPayload);
    },

    recordStop() {
      var time = this.getCurrentTime();
      var socketPayload = {
        EventName: "recordStop",
        EventPayload: {
          username: this.conn.username,
          message: "",
          time: time
        }
      };
      this.sendWebsocketMessage(socketPayload);
    },

    recordPlay() {
      var time = this.getCurrentTime();
      var socketPayload = {
        EventName: "recordPlay",
        EventPayload: {
          username: this.conn.username,
          message: "",
          time: time
        }
      };
      this.sendWebsocketMessage(socketPayload);
    },

    // Clamps a number to a range
    clamp(num, min, max) {
      return Math.max(min, Math.min(max, num));
    },

    calculateOctave(n) {
      return (
        Math.floor(n / WHITE_KEYS_PER_OCT) +
        Math.max(MIN_OCTAVE, this.offsets.octaveStart)
      );
    },

    // Probably should abstract this since it's used in piano.vue and app.vue
    sendWebsocketMessage(socketPayload) {
      this.connection.ws.send(
        JSON.stringify({
          eventName: socketPayload.EventName,
          EventPayload: {
            username: this.connection.username,
            message: socketPayload.EventPayload.message,
            time: socketPayload.EventPayload.time
          }
        })
      );
      console.log(
        "json on frontend is",
        JSON.stringify({
          eventName: socketPayload.EventName,
          EventPayload: {
            username: this.connection.username,
            message: socketPayload.EventPayload.message,
            time: socketPayload.EventPayload.time
          }
        })
      );
      console.log("eventname on frontend is", socketPayload.EventName);
    },

    toggleTrue(note) {
      this.addActiveClass(note);
      var option = {
        playback: "",
        sentBy: ""
      };
      this.playSound(option);
      var time = this.getCurrentTime();
      var socketPayload = {
        EventName: "keyboardPress",
        EventPayload: {
          username: this.conn.username,
          message: String(note),
          time: time
        }
      };
      this.sendWebsocketMessage(socketPayload);
    },

    toggleFalse(note) {
      this.removeActiveClass(note);
    },

    setWhiteKeys(keys) {
      for (let i = this.offsetStart, j = 0; j < this.totalWhiteKeys; i++, j++) {
        const octave = this.calculateOctave(i);
        const keyName = WHITE_KEYS[i % 7];

        const key = {
          name: `${keyName}${octave}`,
          class: ["offwhite", keyName, `${keyName}${octave}`],
          style: {
            "grid-column": `${j === 0 ? 1 : 4 + (j - 1) * 3} / span 3`
          }
        };

        keys.push(key);
      }
    },

    setBlackKeys(keys) {
      for (let i = this.offsetStart, j = 0; j < this.totalWhiteKeys; i++, j++) {
        const octave = this.calculateOctave(i);
        const keyName = BLACK_KEYS[i % 7];

        if (!keyName || octave >= 8) {
          continue;
        }

        const keyNameClass = keyName.replace("#", "s");

        const key = {
          name: `${keyName}${octave}`,
          class: ["offblack", keyNameClass, `${keyNameClass}${octave}`],
          style: {
            "grid-column": `${j === 0 ? 3 : 6 + (j - 1) * 3} / span 2`
          }
        };

        keys.push(key);
      }
    },

    keyDownMonitor(response) {
      var keyPressed = response.event.keyCode;
      this.addActiveClass(keyPressed);
      var time = this.getCurrentTime();
      var socketPayload = {
        EventName: "keyboardPress",
        EventPayload: {
          username: this.conn.username,
          message: String(keyPressed),
          time: time
        }
      };
      this.sendWebsocketMessage(socketPayload);
      console.log(this.keyPressedName);
      console.log(`${this.keyPressedName}.mp3`);
      var option = {
        playback: "",
        sentBy: ""
      };
      this.playSound(option);
    },

    keyUpMonitor(response) {
      var keyLifted = response.event.keyCode;
      this.removeActiveClass(keyLifted);
    },

    setWebsocketMessageListener() {
      this.conn.ws.onmessage = messageEvent => {
        const socketPayload = JSON.parse(messageEvent.data);
        console.log("YO we out here...", socketPayload);
        console.log("YO we out here...", socketPayload.eventName);

        switch (socketPayload.eventName) {
          case "keyboardPress": {
            console.log("YO we in here...", socketPayload);

            const messageContent = socketPayload.EventPayload;
            const sentBy = messageContent.username;
            const actualMessage = messageContent.message;
            console.log("YO EventPayload...", socketPayload.EventPayload);
            console.log({
              messageFrom: sentBy,
              message: actualMessage
            });

            var keyPressed = actualMessage;
            this.addActiveClass(keyPressed);
            console.log(this.keyPressedName);
            console.log(`${this.keyPressedName}.mp3`);
            var option = {
              playback: "listen",
              sentBy: messageContent.username
            };
            this.playSound(option);

            setTimeout(this.removeActiveClass(keyPressed), 1000);
            break;
          }
          default: {
            break;
          }
        }
      };
    },

    removeActiveClass(keyPressed) {
      let keys = this.keysData;
      for (var key of keys) {
        if (key.keyCode == keyPressed) {
          let classString = String(
            key.class[0] + " " + key.class[1] + " " + key.class[2]
          );

          this.keyPressedName = key.class[2];

          document
            .getElementsByClassName(classString)[0]
            .classList.remove("active");
        }
      }
    },

    addActiveClass(keyPressed) {
      let keys = this.keysData;
      for (var key of keys) {
        if (key.keyCode == keyPressed) {
          let classString = String(
            key.class[0] + " " + key.class[1] + " " + key.class[2]
          );

          this.keyPressedName = key.class[2];

          document
            .getElementsByClassName(classString)[0]
            .classList.add("active");
        }
      }
    },

    getCurrentTime() {
      var d = new Date();
      var time = d.getTime();
      return time;
    },

    playSound(option) {
      var keyPressedName = this.keyPressedName;
      var sound;
      if (option.playback == "listen") {
        sound = new Howl({
          src: [`${keyPressedName}.mp3`],
          html5: true,
          autoplay: true,
          volume: 1.0,
          format: "mp3",
          onload: function() {
            console.log("song loaded!");
          }
        });
        if (option.sentBy != this.conn.username) {
          sound.play();
        }
      } else {
        sound = new Howl({
          src: [`${keyPressedName}.mp3`],
          html5: true,
          autoplay: true,
          volume: 1.0,
          format: "mp3",
          onload: function() {
            console.log("song loaded!");
          }
        });
        sound.play();
      }
      console.log(sound.state());
    }
  }
};
</script>

<style scoped>
.keyboard {
  margin: 10%;
  width: 75vw;
  height: calc(260px - calc(var(--octaves) * 10px));
}

.keyboard ul {
  height: 100%;
  width: 100%;
  list-style-type: none;
  display: grid;
  grid-template-columns: repeat(calc(var(--keys) * 3), 1fr);
  grid-template-rows: repeat(3, 1fr);
}

li {
  text-align: center;
  background-color: white;
  border: 1px solid black;
  display: flex;
  justify-content: center;
  align-items: flex-end;
  padding-bottom: 0.25rem;
  font-weight: bold;
}

li.black span {
  transform: rotate(90deg);
  transform-origin: center 50%;
  margin-bottom: 0.75rem;
}

.offwhite {
  grid-row: 1 / span 3;
  z-index: 2;
  border-color: black !important;
}

.offblack {
  grid-row: 1 / span 2;
  background-color: black;
  color: white;
  z-index: 3;
}

.blank {
  border-width: 0;
  grid-row: 1 / span 2;
}

li {
  transition: background-color 0.2s;
}

/* One key off so just gonna hide it :) */
.Cs6 {
  visibility: hidden !important;
}

.active {
  background-color: rgb(255, 0, 0) !important;
}
</style>
