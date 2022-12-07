module.exports = {
  apps: [
    {
      name: "ezvj",
      script: "go run .",
    },
    {
      name: "vlc",
      script: "vlc --http-host=localhost",
    },
  ],
};
