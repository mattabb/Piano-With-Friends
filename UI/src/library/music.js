import { note, transpose } from "@tonaljs/tonal";
import { enharmonic } from "@tonaljs/note";

// Create the range of notes we will be using using tonalJS
export function createRange(from, to) {
  let fromNote = note(from);
  let toNote = note(to);

  if (fromNote.height >= toNote.height) {
    throw new Error("Reverse ranges are not implemented at this time.");
  }

  // if the accidentals are "b" change them to "#"
  if (fromNote.acc === "b") {
    fromNote = note(enharmonic(fromNote));
  }

  let range = [];

  // Push the notes onto our range array
  for (
    let i = 0, l = toNote.height - fromNote.height, currNote = fromNote;
    i < l;
    i++
  ) {
    range.push(currNote);
    currNote = note(enharmonic(transpose(currNote, "m2")));
  }

  return range;
}
