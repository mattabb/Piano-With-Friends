import Vue from "vue";
import { createRange } from "./music";

const notes = createRange("A0", "C8");

const noteMap = notes.reduce((map, note) => {
  map[note.name] = false;
  return map;
}, {});

// This stores the state of the piano  and basically stores what note is clicked and what isnt
/* ex: 	
		C#: false, 
		C: false, 
		A: true,
		etc.
*/
const pianoState = new Vue.observable(noteMap);

export default pianoState;

// This resets the piano and sets every note to false (not being clicked)
export function reset() {
  for (const note of notes) {
    pianoState[note.name] = false;
  }
}
