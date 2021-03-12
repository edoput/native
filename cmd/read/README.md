example read from web-extension
===============================

This command implements a reader that deserialize JSON messages
coming from the web-extension. The `native` package implements
JSON decoding over the stdout/stdin pair with syncronisation
for multiple goroutines.
