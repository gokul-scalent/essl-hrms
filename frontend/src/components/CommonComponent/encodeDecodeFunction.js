const characters = "abcdefghijklmnopqrstuvwxyz1234567890"; // 36 chars
const BASE = characters.length;
const SALT = 9999999; //for more security

export const urlEncodedId = (id) => {
  //id=5
  let num = Number(id) + SALT; // 5+9999999 = 10000004
  let str = "";
  while (num > 0) {
    str = characters[num % BASE] + str;
    // iteration -> 10000004 % 36 = "add some random string n every iterate"
    num = Math.floor(num / BASE);
  }
  return str;
};

export const urlDecodedId = (str) => {
  // str = "f8ynq"
  let num = 0;
  for (let i = 0; i < str.length; i++) {
    num = num * BASE + characters.indexOf(str[i]);
    // iterate according to length of string num = 278221 * 36 + 16 = 10000004
  }
  return num - SALT; // return 10000004 - 9999999 = 5
};