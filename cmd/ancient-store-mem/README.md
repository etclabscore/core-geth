# Ancient Store Memory Mapped Test Util

This application is intended for testing purposed only. Ancient data is stored ephemerally.
The program expects first and only argument to be an IPC path, or, the directory
in which a default 'mock-freezer.ipc' path should be created.
This memory mapped ancient store can also be used as a library.
Package 'lib' logic may be imported and used in testing contexts as well.

## Usage
```
ancient-store-mem your-ipc-path 
```