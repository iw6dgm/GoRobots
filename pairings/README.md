# Pairings

Script to randomly generate tournament rounds for all ever written robots (macro, midi, micro).

It attempts to create rounds with no conflicts (robots with different path but same base name).

It prints out the list of rounds, the YAML configuration to be used with GoRobots and the SQL code to initialise a sqlite3 database (must be already setup with the Crobots schema).