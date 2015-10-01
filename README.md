# DelayBuffer: a timeshifting buffer

[GoDocs with example](https://godoc.org/github.com/jweir/delay)

DelayBuffer is an io.ReadWriter which prevents reading of any written bytes, until after a specified period of time.


was created in order to timeshift streaming media from the East coast to the West as if it were playing live. 

It can be used for any []byte stream.
