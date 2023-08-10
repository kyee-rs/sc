####################################################################################################
# Do not delete this file or any of the specified fields. This results in the program crashing.    #
# 08.08.2023 - @voxelin                                                                            #
####################################################################################################

port = 8080                                                               # Port
database_url = "postgres://user:password@localhost:5432/db"               # TimescaleDB DSN
seed = 3719                                                               # Seed for ID generation [type: random number > 0]

logger {
  force_colors = false                                                    # Force colored output
  full_timestamp = true
}

limits {
  max_size = 10                                                           # Max file size in MB
  block_tor = true                                                        # Block TOR exit nodes
  ip_blacklist = []                                                       # ["168.0.2.3", "2.3.4.5"...]
}