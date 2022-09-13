import * as short from "@pulumi/short-io";

new short.Link("link", {
  domain: "rawkode.link",
  short: "hello",
  long: "google.com",
}, {});
