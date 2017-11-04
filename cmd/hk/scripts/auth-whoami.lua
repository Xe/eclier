script.verb = "auth:whoami"
script.help = "prints the username currently authenticated with heroku."
script.author = "Xe"
script.version = "0.1"
script.usage = ""

function run()
  print(netrc:machine("api.heroku.com"):get("login"))
end
