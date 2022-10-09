native
======

WebExtension native process library. A WebExtension can spawn native processes
and communicate over STDIO/STDOUT with it. This package offers a simple interface
to the messaging inspired by the net/http package.

Writing the native end of a WebExtension becomes as easy as implementing
the `ServeNative` interface and calling `ListenAndServe`.

#### Examples

In the `cmd` directory you can find ready to test examples.
