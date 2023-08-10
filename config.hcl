####################################################################################################
# Do not delete this file or any of the specified fields. This results in the program crashing.    #
# 08.08.2023 - @voxelin                                                                            #
####################################################################################################

server "Simple Cache - @voxelin" "Simple Cache - @voxelin" {                                           # Set `app_name` and `server_name` to your desired values. Has no performance or runtime difference.
  port = 8080                                                               # Port
  database_url = "postgres://user:password@localhost:5432/db"               # TimescaleDB Connection URL
  seed = 3719                                                               # Seed for ID generation [type: random number > 0]
}

logger {
  force_colors = false                                                    # Force colored output
  full_timestamp = true
}

limits {
  max_size = 10                                                           # Max file size in MB
  block_tor = true                                                        # Block TOR exit nodes
  ip_blacklist = []                                                       # ["168.0.2.3", "1.2.3.4"...]
}