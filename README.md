# Dwarf Wars

This will be a Dwarf fortress / Rimworld enspired Multiplayer game. I am hoping to make it have the same sandbox story driven feel as it's inspiration but with friends. or strangers, however you want to play.

## TCP Server

Need to figure out how to close a socket connection when the client sends a FIN packet to the server. the problem is looking like some clients will send a RESET connection request if they have data in thier buffer or not.

