Changes
=======

0.5.0
-----

...

0.4.0
-----

* Complete rewrite using `Gin <https://github.com/gin-gonic/gin>`_ and net/context
* Migrate to Go 1.7
* Implement ``noop`` operation
* Fix animated gif handling (still some minor issues)
* Migrate to `viper <https://github.com/spf13/viper>`_ for config handling which adds support for environment variables
* Add ``allowed_headers`` new config for CORS

0.3.0
-----

* Delete handler to delete an image on storage
* Upload handler to upload an image to the storage

0.2.0
-----

* Robust test suite
* Animated gif support
* New operations: flip, rotate
* An interface to implement multiple engines
* Add ``q`` parameter to control the rendering quality
* Add ``fmt`` parameter to force the output format


0.1.0
-----

Initial release
