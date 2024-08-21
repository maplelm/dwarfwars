# Dwarf Wars

This will be a Dwarf fortress / Rimworld enspired Multiplayer game. I am hoping to make it have the same sandbox story driven feel as it's inspiration but with friends. or strangers, however you want to play.

## TCP Server

## Game Scope

### Backend

The backend is written in go and will relay on sqlite files to store information. I want this project to be easy to spin up on anyones system. This is supposed to be a self host game. I am not sure if there are any better technologies out there that would allow for better saving of data then sqlite.

Need to figure out how to close a socket connection when the client sends a FIN packet to the server. the problem is looking like some clients will send a RESET connection request if they have data in thier buffer or not.

#### Save Data Options

- Protocol Buffers (ProtoBuf)
- FlatBuffers
- MessagePack
- Cap'n Proto
- Avro
- CBOR

I could also just use a database but I don't know how well that will scale, I probably don't need it to scale that much anyways as this will not be a very big product. I don't plan on it being anyways, we will get there when we get there if we get there. I am thinking that just using a log binary blob in a nosql database would work the best.

### Frontend

