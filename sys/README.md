### In Memory File System

File System
-----------
The filesystem type is a bit specific, as it is a map with key type string and value type *zip.File, I'd like to expand this later, but for now this is all I need.

Manager
-------
File manager is a thin wrapper over aws sdk's s3manager uploader and downloader, meant to be used with mem.File

Mem
-----
mem.File is an in memory file that satisfies the io.WriterAt and io.ReaderAt interfaces.
The s3manager uploader and downloader usually read directly from a file pointer and this was created when there is very little disk space to work with (ex: AWS Lambda functions)
