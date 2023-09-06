TZX Player
==========

A player for 8bit computers tapes files in TZX format.

This project is **WIP**. 

Features:

- Play/stop/rewind etc. controls through keyboard shortcuts
- Export to WAV file
- Counter support (reset, goto etc..)
- Tape player remote control through GPIO module input soldered in the CPC tape player PCB

WIP state:

- It only supports Amstrad CPC "CDT" files for now.
- Not all the blocks types of the TZX specification are supported
- Tested on Linux only
- Looking for a more standard interface for GPIO remote control feature. My current device, which is a 
  Numato Labs 8 channels, does not support standard kernel libgpiod framework and expose a proprietary 
  command set which requires specific code.

Build: 

    $ ./build.sh

exec: 

    $ ./bin/tzx-player COMMAND
 
or to show usage 

    $ ./bin/tzx-player COMMAND

This project is written in Go. It makes use of PortAudio library for audio output.