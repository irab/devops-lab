// Auto-generated by the Load Impact converter

import "./libs/shim/core.js";

export let options = { maxRedirects: 2, iterations: "1000" };

const Request = Symbol.for("request");
postman[Symbol.for("initial")]({
  options
});

let file_path = '/workspace/cloudrun_address'
const URL = open(file_path) + "/version"

export default function() {
  postman[Request]({
    name: "devops-lab1",
    id: "015cb1b5-c0da-4c8f-9845-29952dda213b",
    method: "GET",
    address: URL,
    headers: {
      "Content-Type": "application/x-www-form-urlencoded"
    }
  });
}
