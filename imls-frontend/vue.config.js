import process from "process";

module.exports = {
  publicPath:
    process.env.NODE_ENV === "production"
      ? process.env.BASEURL
      : "/",
};
