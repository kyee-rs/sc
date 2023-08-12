####################################################################################################
# Do not delete this file or any of the specified fields. This results in the program crashing.    #
# 08.08.2023 - @voxelin                                                                            #
####################################################################################################

server {
  serverName = "Simple Cache v1.2.4"                                         # Sets the `Server` header for each response
  appName = "Simple Cache"                                                   # Sets the app name to display in terminal
  port = 8080                                                                # Server port
  databaseUrl = "postgres://user:password@localhost:5432/timescaledb"        # TimescaleDB Connection URL
  seed = 3719                                                                # Seed for ID generation [type: random number > 0]
}

logger {
  forceColors = false                                                        # Force colored output
  fullTimestamp = true
}

limits {
  maxSize = 10                                                               # Max uploaded body size in MB
  blockTor = false                                                           # Block TOR exit nodes
  ipBlacklist = []                                                           # ["168.0.2.3", "1.2.3.4"...]
}