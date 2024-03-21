# grain

grain is a tiny gin-alike minimal middleware provider for `net/http` created
for educational purposes.

This package implements an absolute bare minimum (~50 lines of code) middleware
system that is meant to be used natively with `net/http`. It purposefully does
not include "features" like error catching, request logging, etc. It is not a
web server! It just generates an http.Handler for you to register when creating
your server.

Instead of offering these features, it offers simplicity and flexibilty: it's
easy enough for you to write your own error handler and logging middleware
functions and either add them to your middleware chain or call them from your
request's registered middleware directly.

See the `example/` directory for usage.
