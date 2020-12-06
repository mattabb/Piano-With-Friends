const pianoKeyCodes = [
  {
    //key: "Z",
    keyCode: 90
  },
  {
    //key: "C",
    keyCode: 67
  },
  {
    //key: "V",
    keyCode: 86
  },
  {
    //key: "G",
    keyCode: 71
  },
  {
    //key: "D",
    keyCode: 68
  },
  {
    //key: "S",
    keyCode: 83
  },
  {
    //key: "Q",
    keyCode: 81
  },
  {
    //key: "E",
    keyCode: 69
  },
  {
    //key: "T",
    keyCode: 84
  },
  {
    //key: "5",
    keyCode: 53
  },
  {
    //key: "3",
    keyCode: 51
  },
  {
    //key: "1",
    keyCode: 49
  },
  {
    //key: "6",
    keyCode: 54
  },
  {
    //key: "8",
    keyCode: 56
  },
  {
    //key: "0",
    keyCode: 48
  },
  {
    //key: "O",
    keyCode: 79
  },
  {
    //key: "I",
    keyCode: 73
  },
  {
    //key: "Y",
    keyCode: 89
  },
  {
    //key: "J",
    keyCode: 74
  },
  {
    //key: "K",
    keyCode: 75
  },
  {
    //key: ";",
    keyCode: 186
  },
  {
    //key: ".",
    keyCode: 190
  },
  {
    //key: "M",
    keyCode: 77
  },
  {
    //key: "N",
    keyCode: 78
  },
  {
    //key: "X",
    keyCode: 88
  },
  {
    //key: "B",
    keyCode: 66
  },
  {
    //key: "F",
    keyCode: 70
  },
  {
    //key: "A",
    keyCode: 65
  },
  {
    //key: "W",
    keyCode: 87
  },
  {
    //key: "R",
    keyCode: 82
  },
  {
    //key: "4",
    keyCode: 52
  },
  {
    //key: "2",
    keyCode: 50
  },
  {
    //key: "7",
    keyCode: 55
  },
  {
    //key: "9",
    keyCode: 57
  },
  {
    //key: "P",
    keyCode: 80
  },
  {
    //key: "U",
    keyCode: 85
  },
  {
    //key: "H",
    keyCode: 72
  },
  {
    //key: "L",
    keyCode: 76
  },
  {
    //key: "/",
    keyCode: 191
  },
  {
    //key: ",",
    keyCode: 188
  },
  {
    //key: "",
    keyCode: null
  }
];

export function addKeyCodeToKeys(keys) {
  for (let i = 0; i < keys.length; i++) {
    keys[i] = { ...keys[i], ...pianoKeyCodes[i] };
  }
}
