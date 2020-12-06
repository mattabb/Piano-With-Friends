<template>
  <div
    class="keyboard"
    :style="style"
    @keyup="keyUpMonitor"
    @keydown="keyDownMonitor"
  >
    <ul>
      <li
        v-for="(key, index) in keys"
        :key="index"
        :style="key.style"
        @mousedown="toggleActive(key.name)"
        @mouseup="toggleActive(key.name)"
        :class="[...key.class, { active: noteActive(key.name) }]"
      >
        <span>{{ key.name }} </span>
      </li>
    </ul>
  </div>
</template>

<script>
import pianoState from "../library/piano-state";

const WHITE_KEYS = ["C", "D", "E", "F", "G", "A", "B"];
const BLACK_KEYS = ["C#", "D#", null, "F#", "G#", "A#", null];
// const NOTE_OFFSETS = ["C", "D", "E", "F", "G", "A", "B"];
const MIN_OCTAVE = 0;
const MAX_OCTAVE = 7;
const MIN_NOTE = 0;
const MAX_NOTE = 6;
const WHITE_KEYS_PER_OCT = 7;
const BLACK_KEYS_PER_OCT = 5;

export default {
  name: "Piano",
  // Props are basically parameters for vue components
  props: {
    // Octave start prop => Where the piano's octave will start
    octaveStart: {
      type: Number,
      // Validate that octave start passed in is between our max and min
      validator(value) {
        return value >= MIN_OCTAVE && value <= MAX_OCTAVE;
      },
      // if not we start at min octave
      default() {
        return MIN_OCTAVE;
      }
    },

    // Octave end prop => Where the piano's octave will end
    octaveEnd: {
      type: Number,
      validator(value) {
        return value >= MIN_OCTAVE && value <= MAX_OCTAVE;
      },
      default() {
        return MAX_OCTAVE;
      }
    },

    // noteStart prop => Where the piano's notes will start
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

    // noteEnd prop => Where the piano's notes will end
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
    }
  },
  // Our data variables... think of them as this component's global variables (ONLY FOR THIS COMPONENT)
  data: () => ({
    offsets: {
      octaveStart: 0,
      octaveEnd: 3,
      noteStart: 0,
      noteEnd: 0
    }
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
    this.offsets.octaveEnd = this.octaveEnd;

    if (
      this.offsets.octaveStart > this.offsets.octaveEnd ||
      (this.offsets.octaveStart === this.offsets.octaveEnd &&
        this.offsets.noteStart > this.offsets.noteEnd)
    ) {
      throw new Error(
        "The start octave must be lower than or equal to the end octave and the start note must be lower than the end note."
      );
    }
  },

  // See https://vuejs.org/v2/guide/computed.html for an explanation on computed
  // In the simplest sense, they are ways to cut down on ugly in-line javascript expressions
  computed: {
    // This returns the state of the piano, look at piano-state.js for a better explanation
    pianoState() {
      return pianoState;
    },

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
          (WHITE_KEYS_PER_OCT - this.offsetEnd + 1)
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

      // White keys
      for (let i = this.offsetStart, j = 0; j < this.totalWhiteKeys; i++, j++) {
        const octave = this.calculateOctave(i);
        const keyName = WHITE_KEYS[i % 7];

        const key = {
          name: `${keyName}${octave}`,
          class: ["white", keyName, `${keyName}${octave}`],
          style: {
            "grid-column": `${j === 0 ? 1 : 4 + (j - 1) * 3} / span 3`
          }
        };

        keys.push(key);
      }

      // Black keys
      for (let i = this.offsetStart, j = 0; j < this.totalWhiteKeys; i++, j++) {
        const octave = this.calculateOctave(i);
        const keyName = BLACK_KEYS[i % 7];

        if (!keyName || octave >= 8) {
          continue;
        }

        const keyNameClass = keyName.replace("#", "s");

        const key = {
          name: `${keyName}${octave}`,
          class: ["black", keyNameClass, `${keyNameClass}${octave}`],
          style: {
            "grid-column": `${j === 0 ? 3 : 6 + (j - 1) * 3} / span 2`
          }
        };

        keys.push(key);
      }

      return keys;
    }
  },
  methods: {
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

    toggleActive(note) {
      console.log(pianoState);
      pianoState[note] === true
        ? (pianoState[note] = false)
        : (pianoState[note] = true);
    },

    keyDownMonitor: function(event) {
      var keyPressed = event.keyCode;
      console.log(pianoState);
      switch (keyPressed) {
        case 192:
          document
            .getElementByClassName("white C C2")
            .classList.add(this.active);
          document.toggleActive("C2");
          break;
        case 49:
          document
            .getElementByClassName("black Cs Cs2")
            .classList.add(this.active);
          document.toggleActive("Cs2");
          break;
        case 50:
          document
            .getElementByClassName("white D D2")
            .classList.add(this.active);
          document.toggleActive("D2");
          break;
        case 51:
          document
            .getElementByClassName("black D Ds2")
            .classList.add(this.active);
          document.toggleActive("Ds2");
          break;
        case 52:
          document
            .getElementByClassName("white E E2")
            .classList.add(this.active);
          document.toggleActive("E2");
          break;
        case 53:
          document
            .getElementByClassName("white F F2")
            .classList.add(this.active);
          document.toggleActive("F2");
          break;
        case 54:
          document
            .getElementByClassName("black F Fs2")
            .classList.add(this.active);
          document.toggleActive("Fs2");
          break;
        case 55:
          document
            .getElementByClassName("white G G2")
            .classList.add(this.active);
          document.toggleActive("G2");
          break;
        case 56:
          document
            .getElementByClassName("black Gs Gs2")
            .classList.add(this.active);
          document.toggleActive("Gs2");
          break;
        case 57:
          document
            .getElementByClassName("white A A2")
            .classList.add(this.active);
          document.toggleActive("A2");
          break;
        case 48:
          document
            .getElementByClassName("black As As2")
            .classList.add(this.active);
          document.toggleActive("As2");
          break;
        case 189:
          document
            .getElementByClassName("white B B2")
            .classList.add(this.active);
          document.toggleActive("B2");
          break;
        case 187:
          document
            .getElementByClassName("white C C3")
            .classList.add(this.active);
          document.toggleActive("C3");
          break;
        case 81:
          document
            .getElementByClassName("black Cs Cs3")
            .classList.add(this.active);
          document.toggleActive("Cs3");
          break;
        case 87:
          document
            .getElementByClassName("white D D3")
            .classList.add(this.active);
          document.toggleActive("D3");
          break;
        case 69:
          document
            .getElementByClassName("black Ds Ds3")
            .classList.add(this.active);
          document.toggleActive("Ds3");
          break;
        case 82:
          document
            .getElementByClassName("white E E3")
            .classList.add(this.active);
          document.toggleActive("E3");
          break;
        case 84:
          document
            .getElementByClassName("white F F3")
            .classList.add(this.active);
          document.toggleActive("F3");
          break;
        case 89:
          document
            .getElementByClassName("black Fs Fs3")
            .classList.add(this.active);
          document.toggleActive("Fs3");
          break;
        case 85:
          document
            .getElementByClassName("white G G3")
            .classList.add(this.active);
          document.toggleActive("G3");
          break;
        case 73:
          document
            .getElementByClassName("black Gs Gs3")
            .classList.add(this.active);
          document.toggleActive("Gs3");
          break;
        case 79:
          document
            .getElementByClassName("white A A3")
            .classList.add(this.active);
          document.toggleActive("A3");
          break;
        case 80:
          document
            .getElementByClassName("black As As3")
            .classList.add(this.active);
          document.toggleActive("As3");
          break;
        case 219:
          document
            .getElementByClassName("white B B3")
            .classList.add(this.active);
          document.toggleActive("B3");
          break;
        case 221:
          document
            .getElementByClassName("white C C4")
            .classList.add(this.active);
          document.toggleActive("C4");
          break;
        case 65:
          document
            .getElementByClassName("black Cs Cs4")
            .classList.add(this.active);
          document.toggleActive("Cs4");
          break;
        case 83:
          document
            .getElementByClassName("white D D4")
            .classList.add(this.active);
          document.toggleActive("D4");
          break;
        case 68:
          document
            .getElementByClassName("black Ds Ds4")
            .classList.add(this.active);
          document.toggleActive("Ds4");
          break;
        case 70:
          document
            .getElementByClassName("white E E4")
            .classList.add(this.active);
          document.toggleActive("E4");
          break;
        case 71:
          document
            .getElementByClassName("white F F4")
            .classList.add(this.active);
          document.toggleActive("F4");
          break;
        case 72:
          document
            .getElementByClassName("black Fs Fs4")
            .classList.add(this.active);
          document.toggleActive("Fs4");
          break;
        case 74:
          document
            .getElementByClassName("white G G4")
            .classList.add(this.active);
          document.toggleActive("G4");
          break;
        case 75:
          document
            .getElementByClassName("black Gs Gs4")
            .classList.add(this.active);
          document.toggleActive("Gs4");
          break;
        case 76:
          document
            .getElementByClassName("white A A4")
            .classList.add(this.active);
          document.toggleActive("A4");
          break;
        case 186:
          document
            .getElementByClassName("black As As4")
            .classList.add(this.active);
          document.toggleActive("As4");
          break;
        case 222:
          document
            .getElementByClassName("white B B4")
            .classList.add(this.active);
          document.toggleActive("B4");
          break;
        case 90:
          document
            .getElementByClassName("white C C5")
            .classList.add(this.active);
          document.toggleActive("C5");
          break;
        case 88:
          document
            .getElementByClassName("black Cs Cs5")
            .classList.add(this.active);
          document.toggleActive("Cs5");
          break;
        case 67:
          document
            .getElementByClassName("white D D5")
            .classList.add(this.active);
          document.toggleActive("D5");
          break;
        case 86:
          document
            .getElementByClassName("black Ds Ds5")
            .classList.add(this.active);
          document.toggleActive("Ds5");
          break;
        case 66:
          document
            .getElementByClassName("white E E5")
            .classList.add(this.active);
          document.toggleActive("E5");
          break;
        case 78:
          document
            .getElementByClassName("white F F5")
            .classList.add(this.active);
          document.toggleActive("F5");
          break;
        case 77:
          document
            .getElementByClassName("black Fs Fs5")
            .classList.add(this.active);
          document.toggleActive("Fs5");
          break;
        case 188:
          document
            .getElementByClassName("white G G5")
            .classList.add(this.active);
          document.toggleActive("G5");
          break;
        case 190:
          document
            .getElementByClassName("black Gs Gs5")
            .classList.add(this.active);
          document.toggleActive("Gs5");
          break;
        case 191:
          document
            .getElementByClassName("white A A5")
            .classList.add(this.active);
          document.toggleActive("A5");
          break;
        case 37:
          document
            .getElementByClassName("black As As5")
            .classList.add(this.active);
          document.toggleActive("As5");
          break;
        case 39:
          document
            .getElementByClassName("white B B5")
            .classList.add(this.active);
          document.toggleActive("B5");
          break;
      }
    },

    keyUpMonitor: function(event) {
      var keyPressed = event.keyCode;
      console.log(pianoState);
      switch (keyPressed) {
        case 192:
          document
            .getElementByClassName("white C C2")
            .classList.remove(this.active);
          document.toggleActive("C2");
          break;
        case 49:
          document
            .getElementByClassName("black Cs Cs2")
            .classList.remove(this.active);
          document.toggleActive("Cs2");
          break;
        case 50:
          document
            .getElementByClassName("white D D2")
            .classList.remove(this.active);
          document.toggleActive("D2");
          break;
        case 51:
          document
            .getElementByClassName("black D Ds2")
            .classList.remove(this.active);
          document.toggleActive("Ds2");
          break;
        case 52:
          document
            .getElementByClassName("white E E2")
            .classList.remove(this.active);
          document.toggleActive("E2");
          break;
        case 53:
          document
            .getElementByClassName("white F F2")
            .classList.remove(this.active);
          document.toggleActive("F2");
          break;
        case 54:
          document
            .getElementByClassName("black F Fs2")
            .classList.remove(this.active);
          document.toggleActive("Fs2");
          break;
        case 55:
          document
            .getElementByClassName("white G G2")
            .classList.remove(this.active);
          document.toggleActive("G2");
          break;
        case 56:
          document
            .getElementByClassName("black Gs Gs2")
            .classList.remove(this.active);
          document.toggleActive("Gs2");
          break;
        case 57:
          document
            .getElementByClassName("white A A2")
            .classList.remove(this.active);
          document.toggleActive("A2");
          break;
        case 48:
          document
            .getElementByClassName("black As As2")
            .classList.remove(this.active);
          document.toggleActive("As2");
          break;
        case 189:
          document
            .getElementByClassName("white B B2")
            .classList.remove(this.active);
          document.toggleActive("B2");
          break;
        case 187:
          document
            .getElementByClassName("white C C3")
            .classList.remove(this.active);
          document.toggleActive("C3");
          break;
        case 81:
          document
            .getElementByClassName("black Cs Cs3")
            .classList.remove(this.active);
          document.toggleActive("Cs3");
          break;
        case 87:
          document
            .getElementByClassName("white D D3")
            .classList.remove(this.active);
          document.toggleActive("D3");
          break;
        case 69:
          document
            .getElementByClassName("black Ds Ds3")
            .classList.remove(this.active);
          document.toggleActive("Ds3");
          break;
        case 82:
          document
            .getElementByClassName("white E E3")
            .classList.remove(this.active);
          document.toggleActive("E3");
          break;
        case 84:
          document
            .getElementByClassName("white F F3")
            .classList.remove(this.active);
          document.toggleActive("F3");
          break;
        case 89:
          document
            .getElementByClassName("black Fs Fs3")
            .classList.remove(this.active);
          document.toggleActive("Fs3");
          break;
        case 85:
          document
            .getElementByClassName("white G G3")
            .classList.remove(this.active);
          document.toggleActive("G3");
          break;
        case 73:
          document
            .getElementByClassName("black Gs Gs3")
            .classList.remove(this.active);
          document.toggleActive("Gs3");
          break;
        case 79:
          document
            .getElementByClassName("white A A3")
            .classList.remove(this.active);
          document.toggleActive("A3");
          break;
        case 80:
          document
            .getElementByClassName("black As As3")
            .classList.remove(this.active);
          document.toggleActive("As3");
          break;
        case 219:
          document
            .getElementByClassName("white B B3")
            .classList.remove(this.active);
          document.toggleActive("B3");
          break;
        case 221:
          document
            .getElementByClassName("white C C4")
            .classList.remove(this.active);
          document.toggleActive("C4");
          break;
        case 65:
          document
            .getElementByClassName("black Cs Cs4")
            .classList.remove(this.active);
          document.toggleActive("Cs4");
          break;
        case 83:
          document
            .getElementByClassName("white D D4")
            .classList.remove(this.active);
          document.toggleActive("D4");
          break;
        case 68:
          document
            .getElementByClassName("black Ds Ds4")
            .classList.remove(this.active);
          document.toggleActive("Ds4");
          break;
        case 70:
          document
            .getElementByClassName("white E E4")
            .classList.remove(this.active);
          document.toggleActive("E4");
          break;
        case 71:
          document
            .getElementByClassName("white F F4")
            .classList.remove(this.active);
          document.toggleActive("F4");
          break;
        case 72:
          document
            .getElementByClassName("black Fs Fs4")
            .classList.remove(this.active);
          document.toggleActive("Fs4");
          break;
        case 74:
          document
            .getElementByClassName("white G G4")
            .classList.remove(this.active);
          document.toggleActive("G4");
          break;
        case 75:
          document
            .getElementByClassName("black Gs Gs4")
            .classList.remove(this.active);
          document.toggleActive("Gs4");
          break;
        case 76:
          document
            .getElementByClassName("white A A4")
            .classList.remove(this.active);
          document.toggleActive("A4");
          break;
        case 186:
          document
            .getElementByClassName("black As As4")
            .classList.remove(this.active);
          document.toggleActive("As4");
          break;
        case 222:
          document
            .getElementByClassName("white B B4")
            .classList.remove(this.active);
          document.toggleActive("B4");
          break;
        case 90:
          document
            .getElementByClassName("white C C5")
            .classList.remove(this.active);
          document.toggleActive("C5");
          break;
        case 88:
          document
            .getElementByClassName("black Cs Cs5")
            .classList.remove(this.active);
          document.toggleActive("Cs5");
          break;
        case 67:
          document
            .getElementByClassName("white D D5")
            .classList.remove(this.active);
          document.toggleActive("D5");
          break;
        case 86:
          document
            .getElementByClassName("black Ds Ds5")
            .classList.remove(this.active);
          document.toggleActive("Ds5");
          break;
        case 66:
          document
            .getElementByClassName("white E E5")
            .classList.remove(this.active);
          document.toggleActive("E5");
          break;
        case 78:
          document
            .getElementByClassName("white F F5")
            .classList.remove(this.active);
          document.toggleActive("F5");
          break;
        case 77:
          document
            .getElementByClassName("black Fs Fs5")
            .classList.remove(this.active);
          document.toggleActive("Fs5");
          break;
        case 188:
          document
            .getElementByClassName("white G G5")
            .classList.remove(this.active);
          document.toggleActive("G5");
          break;
        case 190:
          document
            .getElementByClassName("black Gs Gs5")
            .classList.remove(this.active);
          document.toggleActive("Gs5");
          break;
        case 191:
          document
            .getElementByClassName("white A A5")
            .classList.remove(this.active);
          document.toggleActive("A5");
          break;
        case 37:
          document
            .getElementByClassName("black As As5")
            .classList.remove(this.active);
          document.toggleActive("As5");
          break;
        case 39:
          document
            .getElementByClassName("white B B5")
            .classList.remove(this.active);
          document.toggleActive("B5");
          break;
      }
    },

    noteActive(note) {
      return pianoState[note] === true;
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

.white {
  grid-row: 1 / span 3;
  z-index: 2;
}

.black {
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

.active {
  background-color: black !important;
}

.Fs.active {
  background-color: rgb(174, 0, 0);
}

.G.active {
  background-color: rgb(255, 0, 0);
}

.Gs.active {
  background-color: rgb(255, 0, 0);
}

.A.active {
  background-color: rgb(255, 102, 0);
}

.As.active {
  background-color: rgb(255, 239, 0);
}

.B.active {
  background-color: rgb(153, 255, 0);
}

.C.active {
  background-color: rgb(0, 40, 255);
}

.Cs.active {
  background-color: rgb(0, 255, 242);
}

.D.active {
  background-color: rgb(0, 122, 255);
}

.Ds.active {
  background-color: rgb(5, 0, 255);
}

.E.active {
  background-color: rgb(71, 0, 237);
}

.F.active {
  background-color: rgb(99, 0, 178);
}
</style>
