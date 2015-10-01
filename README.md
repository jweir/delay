# DelayBuffer: a timeshifting buffer

[GoDocs with example](https://godoc.org/github.com/jweir/delay)

DelayBuffer, which implements io.ReadWriter, prevents reading of any written bytes until after a specified period of time.


It was created in order to timeshift a live MP3 stream from the East coast to the West as if it were playing live. 
