# stream-vbyte-go
A port of Stream VByte to Go

Stream VByte is a variable-length unsigned int encoding designed to make SIMD processing more efficient.

See https://lemire.me/blog/2017/09/27/stream-vbyte-breaking-new-speed-records-for-integer-compression/ and https://arxiv.org/pdf/1709.08990.pdf for details on the format.

The reference C implementation is https://github.com/lemire/streamvbyte.

There is also a Rust implementation https://bitbucket.org/marshallpierce/stream-vbyte-rust.
