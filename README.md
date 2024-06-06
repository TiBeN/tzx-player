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

- Not all blocks types of the TZX specification are supported
  (Most data related though, which make almost every Amstrad CPC cdt file to be read)
- Tested on Linux only
- Looking for a more standard interface for GPIO remote control feature. My current device, which is a 
  Numato Labs 8 channels, does not support standard kernel libgpiod framework and expose a proprietary 
  command set which requires specific code.

Build
-----

    $ ./build.sh

exec 

    $ ./bin/tzx-player COMMAND
 
or to show usage 

    $ ./bin/tzx-player help

Usage
-----

```
Usage: tzx-player COMMAND [cmd opts]

Play 8bits computers data tapes as TZX files

Commands:
  help                Show this help message
  convert             Convert TZX tape to an audio PCM Wav file
    Args:
      tzx-player convert INPUT_TZX_FILE OUTPUT_WAV_FILE
    Options:
      -s int              Sampling rate (default: 44100)
      -b int              Bit depth (default: 8, possibles values: 8 or 16)
      -f float            Speed factor: multiply the speed of the tones (experimental) (default: 1.0)
  info                Output TZX tape informations
    Args:
      tzx-player info INPUT_TZX_FILE
  play                Play a TZX tape
    Args:
      tzx-player play INPUT_TZX_FILE
    Options:
      -s int              Sampling rate (default: 44100)
      -b int              Bit depth (default: 8, possibles values: 8 or 16)
      -g port:baudrate:ionbEnable tape remote control using a GPIO device. Support only Numato labs GPIO Modules for now
      -f float            Speed factor: multiply the speed of the tones (experimental) (default: 1.0)
   Player control keystrokes:
       Space : Toggle play/pause
       p : Pause
       s : Save current tape position
       g : Set tape to last saved position
```

This project is written in Go. It makes use of PortAudio library for audio output.

Resources
---------

The TZX tape file format specification is available [here](https://k1.spdns.de/Develop/Projects/zasm/Info/TZX%20format.html)
