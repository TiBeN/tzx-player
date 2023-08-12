TZX Player
==========

A player for 8bit computers tapes files in TZX format.

This project is **WIP**: it only supports Amstrad CPC "CDT" files for now with WAV file format output.

Expected goals:

- Integration of an autonomous player (pulse audio/pipewire...)
  - Play/stop/rewind etc. controls through keyboard shortcuts
  - Counter support (reset, goto etc..)
  - Play/Stop control for Amstrad CPC 464 through GPIO port connected (soldered) to the integrated tape "datacorder" relay
- Support for others machines tapes that are supported by the TZX format
- Support more audio formats for audio file output (.au, .voc, .mp3 etc.)