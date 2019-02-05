# Go TurboPFor teaching implementation

This teaching implementation illustrates the accompanying article at
https://michael.stapelberg.de/posts/2019-02-05-turbopfor-analysis/

To confirm my understanding of the details of the format, I implemented a
pure-Go TurboPFor256 decoder. Note that it is intentionally *not optimized* as
its main goal is to use simple code to teach the TurboPFor256 on-disk format.

If you’re looking to use TurboPFor from Go, I recommend using cgo. cgo’s
function call overhead is about 51ns [as of Go
1.8](https://go-review.googlesource.com/c/go/+/30080), which will easily be
offset by TurboPFor’s carefully optimized, vectorized (SSE/AVX) code.
