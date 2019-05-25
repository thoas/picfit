picfit
======

.. image:: https://secure.travis-ci.org/thoas/picfit.png?branch=master
    :alt: Build Status
    :target: http://travis-ci.org/thoas/picfit

.. image:: https://d262ilb51hltx0.cloudfront.net/max/800/1*oR04S6Ie7s1JctwjsDsN0w.png

picfit is a reusable Go server to manipulate images (resize, thumbnail, etc.).

It will act as a proxy on your storage engine and will be
served ideally behind an HTTP cache system like varnish_.

It supports multiple `storage backends <https://github.com/ulule/gostorages>`_
and multiple `key/value stores <https://github.com/ulule/gokvstores>`_.

Installation
============

Build it
--------

1. Make sure you have a Go language compiler and git installed.
2. Make sure you have the following go system dependencies in your $PATH: bzr, svn, hg, git
3. Ensure your GOPATH_ is properly set.
4. Download it:

::

    git clone https://github.com/thoas/picfit.git

4. Run ``make build``

You have now a binary version of picfit in the ``bin`` directory which
fits perfectly with your architecture.

picfit has also a Docker integration to built a unix binary without having to install it

::

    make docker-build

Debian and Ubuntu
-----------------

We will provide Debian package when we will be completely stable ;)

Configuration
=============

Configuration should be stored in a readable file and in JSON format.

``config.json``

.. code-block:: json

    {
      "kvstore": {
        "type": "[KVSTORE]"
      },
      "storage": {
        "src": {
          "type": "[STORAGE]"
        }
      }
    }

``[KVSTORE]`` can be:

- **redis** - generated keys stored in Redis_, see `below <#store-images-on-amazon-s3-keys-in-redis-and-shard-filename>`_ how you can customize connection parameters
- **cache** - generated keys stored in an in-memory cache
- **redis-cluster** - generated keys stored in `Redis cluster <https://redis.io/topics/cluster-tutorial>`_

``[STORAGE]`` can be:

- **fs** - generated images stored in your File system
- **http+fs** - generated images stored in your File system and loaded using HTTP protocol
- **s3** - generated images stored in Amazon S3
- **gcs** - generated images stored in Google Cloud Storage
- **http+s3** - generated images stored in Amazon S3 and loaded using HTTP protocol

Basic
-----

* no key/value store
* no image storage
* images are given in absolute url

``config.json``

.. code-block:: json

    {
      "port": 3001
    }

Images are generated on the fly at each request.

Store images on file system and keys in an in-memory cache
----------------------------------------------------------

* key/value in-memory store
* file system storage

An image is generated from your source storage (``src``) and uploaded
asynchronously to this storage.

A unique key is generated and stored in a in-memory key/value store to process
a request only once.

``config.json``

.. code-block:: json

    {
      "port": 3001,
      "storage": {
        "src": {
          "type": "fs",
          "location": "/path/to/directory/"
        }
      },
      "kvstore": {
        "type": "cache"
      },
    }

Store images on Amazon S3, keys in Redis and shard filename
-----------------------------------------------------------

* key/value store provided by Redis
* Amazon S3 storage
* shard filename

``config.json``

.. code-block:: json

    {
      "kvstore": {
        "type": "redis",
        "redis": {
          "host": "127.0.0.1",
          "port": 6379,
          "password": "",
          "db": 0
        }
      },
      "port": 3001,
      "storage": {
        "src": {
          "type": "s3",
          "access_key_id": "[ACCESS_KEY_ID]",
          "secret_access_key": "[SECRET_ACCESS_KEY]",
          "bucket_name": "[BUCKET_NAME]",
          "acl": "[ACL]",
          "region": "[REGION_NAME]",
          "location": "path/to/directory"
        }
      },
      "shard": {
        "width": 1,
        "depth": 2,
        "restonly": true
      }
    }

Keys will be stored on Redis_, (you better need to setup persistence_).

Image files will be loaded and stored on Amazon S3 at the location ``path/to/directory``
in the bucket ``[BUCKET_NAME]``.

``[ACL]`` can be:

- private
- public-read
- public-read-write
- authenticated-read
- bucket-owner-read
- bucket-owner-full-control

``[REGION_NAME]`` can be:

- us-gov-west-1
- us-east-1
- us-west-1
- us-west-2
- eu-west-1
- eu-central-1
- ap-southeast-1
- ap-southeast-2
- ap-northeast-1
- sa-east-1
- cn-north-1

**Filename** will be sharded:

- ``depth`` - 2 directories
- ``width`` - 1 letter for each directory
- ``restonly`` - true, filename won't contain characters in sharding path

Example:

``06102586671300cd02ae90f1faa16897.png`` will become ``0/6/102586671300cd02ae90f1faa16897.jpg``

with restonly=false it would become ``0/6/06102586671300cd02ae90f1faa16897.jpg``

It would be useful if you are using the file system storage backend.

Load images from file system and store them in Amazon S3, keys on Redis cluster
-------------------------------------------------------------------------------

* key/value store provided by Redis cluster
* File system to load images
* Amazon S3 storage to process images

``config.json``

.. code-block:: json

    {
      "kvstore": {
        "type": "redis-cluster",
        "redis": {
          "addrs": [
            "127.0.0.1:6379"
          ],
          "password": "",
        }
      },
      "port": 3001,
      "storage": {
        "src": {
          "type": "fs",
          "location": "path/to/directory"
        },
        "dst": {
          "type": "s3",
          "access_key_id": "[ACCESS_KEY_ID]",
          "secret_access_key": "[SECRET_ACCESS_KEY]",
          "bucket_name": "[BUCKET_NAME]",
          "acl": "[ACL]",
          "region": "[REGION_NAME]",
          "location": "path/to/directory"
        }
      }
    }

You will be able to load and store your images from different storages backend.

In this example, images will be loaded from the file system storage
and generated to the Amazon S3 storage.

Load images from storage backend base url, store them in Amazon S3, keys prefixed on Redis
------------------------------------------------------------------------------------------

* key/value store provided by Redis
* File system to load images using HTTP method
* Amazon S3 storage to process images

``config.json``

.. code-block:: json

    {
      "kvstore": {
        "type": "redis",
        "redis": {
          "host": "127.0.0.1",
          "port": 6379,
          "password": "",
          "db": 0
        },
        "prefix": "dummy:"
      },
      "port": 3001,
      "storage": {
        "src": {
          "type": "http+fs",
          "base_url": "http://media.example.com",
          "location": "path/to/directory"
        },
        "dst": {
          "type": "s3",
          "access_key_id": "[ACCESS_KEY_ID]",
          "secret_access_key": "[SECRET_ACCESS_KEY]",
          "bucket_name": "[BUCKET_NAME]",
          "acl": "[ACL]",
          "region": "[REGION_NAME]",
          "location": "path/to/directory"
        }
      }
    }

In this example, images will be loaded from the file system storage
using HTTP with ``base_url`` option and generated to the Amazon S3 storage.

Keys will be stored on Redis_ using the prefix ``dummy:``.

Running
=======

To run the application, issue the following command:

::

    $ picfit -c config.json

By default, this will run the application on port 3001 and
can be accessed by visiting:

::

    http://localhost:3001

The port number can be configured with ``port`` option in your config file.

To see a list of all available options, run:

::

    $ picfit --help

Usage
=====

General parameters
------------------

Parameters to call the picfit service are:

::

    <img src="http://localhost:3001/{method}?url={url}&path={path}&w={width}&h={height}&upscale={upscale}&sig={sig}&op={operation}&fmt={format}&q={quality}&deg={degree}&pos={position}"

- **path** - The filepath to load the image using your source storage
- **operation** - The operation to perform, see Operations_
- **sig** - The signature key which is the representation of your query string and your secret key, see Security_
- **method** - The method to perform, see Methods_
- **url** - The url of the image to generate (not required if ``path`` provided)
- **width** - The desired width of the image, if ``0`` is provided the service will calculate the ratio with ``height``
- **height** - The desired height of the image, if ``0`` is provided the service will calculate the ratio with ``width``
- **upscale** - If your image is smaller than your desired dimensions, the service will upscale it by default to fit your dimensions, you can disable this behavior by providing ``0``
- **format** - The output format to save the image, by default the format will be the source format (a ``GIF`` image source will be saved as ``GIF``),  see Formats_
- **quality** - The quality to save the image, by default the quality will be the highest possible, it will be only applied on ``JPEG`` format
- **degree** - The degree (``90``, ``180``, ``270``) to rotate the image
- **position** - The position to flip the image

To use this service, include the service url as replacement
for your images, for example:

::

    <img src="https://www.google.fr/images/srpr/logo11w.png" />

will become:

::

    <img src="http://localhost:3001/display?url=https%3A%2F%2Fwww.google.fr%2Fimages%2Fsrpr%2Flogo11w.png&w=100&h=100&op=resize&upscale=0"

This will retrieve the image used in the ``url`` parameter and resize it
to 100x100.

Using source storage
--------------------

If an image is stored in your source storage at the location ``path/to/file.png``,
then you can call the service to load this file:

::

    <img src="http://localhost:3001/display?w=100&h=100&path=path/to/file.png&op=resize"

    or

    <img src="http://localhost:3001/display/resize/100x100/path/to/file.png"

Formats
=======

picfit currently supports the following image formats:

- ``image/jpeg`` with the keyword ``jpg`` or ``jpeg``
- ``image/png`` with the keyword ``png``
- ``image/gif`` with the keyword ``gif``
- ``image/bmp`` with the keyword ``bmp``

Operations
==========

Resize
------

This operation will able you to resize the image to the specified width and height.

If width or height value is 0, the image aspect ratio is preserved.

-  **w** - The desired image's width
-  **h** - The desired image's height

You have to pass the ``resize`` value to the ``op`` parameter to use this operation.

Thumbnail
---------

Thumbnail scales the image up or down using the specified resample filter,
crops it to the specified width and height and returns the transformed image.

-  **w** - The desired width of the image
-  **h** - The desired height of the image

You have to pass the ``thumbnail`` value to the ``op`` parameter
to use this operation.

Flip
----

Flip flips the image vertically (from top to bottom) or
horizontally (from left to right) and returns the transformed image.

-  **pos** - The desired position to flip the image, ``h`` will flip the image horizontally, ``v`` will flip the image vertically

You have to pass the ``flip`` value to the ``op`` parameter
to use this operation.

Rotate
------

Rotate rotates the image to the desired degree and returns the transformed image.

-  **deg** - The desired degree to rotate the image

You have to pass the ``rotate`` value to the ``op`` parameter
to use this operation.

Flat
----

Flat draws a given image on the image resulted by the previous operation.
Flat can be used only with the [multiple operation system].

- **path** - the foreground image path
- **color** - the foreground color in Hex (without ``#``), default is transparent
- **pos** - the destination rectange

In order to undersand the Flat operation, please read the following `docs <https://github.com/thoas/picfit/blob/superpose-images/docs/flat.md>`_.

Methods
=======

Display
-------

Display the image, useful when you are using an ``img`` tag.

The generated image will be stored asynchronously on your
destination storage backend.

A couple of headers (``Content-Type``, ``If-Modified-Since``) will be set
to allow you to use an http cache system.


Redirect
--------

Redirect to an image.

Your file will be generated synchronously then the redirection
will be performed.

The first query will be slower but next ones will be faster because the name
of the generated file will be stored in your key/value store.

Get
---

Retrieve information about an image.

Your file will be generated synchronously then you will get the following information:

* **filename** - Filename of your generated file
* **path** - Path of your generated file
* **url** - Absolute url of your generated file (only if ``base_url`` is available on your destination storage)

The first query will be slower but next ones will be faster because the name
of the generated file will be stored in your key/value store.

Expect the following result:

.. code-block:: json

    {
        "filename":"a661f8d197a42d21d0190d33e629e4.png",
        "path":"cache/6/7/a661f8d197a42d21d0190d33e629e4.png",
        "url":"https://ds9xhxfkunhky.cloudfront.net/cache/6/7/a661f8d197a42d21d0190d33e629e4.png"
    }

Upload
------

Upload is disabled by default for security reason.
Before enabling it, you must understand you have to secure yourself
this endpoint like only allowing the /upload route in your nginx
or apache webserver for the local network.

Exposing the **/upload** endpoint without a security mechanism is not **SAFE**.

You can enable it by adding the option and a source
storage to your configuration file.

``config.json``

.. code-block:: json

    {
      "storage": {
        "src": {
          "type": "[STORAGE]"
        }
      },
      "options": {
        "enable_upload": true
      }
    }

To work properly, the input field must be named "data"

Test it with the excellent httpie_:

::

    http -f POST localhost:3000/upload data@myupload

You will retrieve the uploaded image information in ``JSON`` format.

Multiple operations
===================

Multiple operations can be done on the same image following a given order.

First operation must be described as above then other operation are described in parameters ``op``.
The order of ``op`` parameters is the order used.

Each options of the operation must be described with subparameters separed by
``:`` with the operation name as argument to ``op``.

Example of a resize followed by a rotation:

::

    <img src="http://localhost:3001/display?w=100&h=100&path=path/to/file.png&op=resize&op=op:rotate+deg:180"

Security
========

Request signing
---------------

In order to secure requests and avoid unknown third parties to
use the service, the application can require a request to provide a signature.
To enable this feature, set the ``secret_key`` option in your config file.

The signature is an hexadecimal digest generated from the client
key and the query string using the HMAC-SHA1 message authentication code
(MAC) algorithm.

The below python code provides an implementation example::

    import hashlib
    import hmac
    import six
    import urllib

    def sign(key, *args, **kwargs):
        m = hmac.new(key, None, hashlib.sha1)

        for arg in args:
            if isinstance(arg, dict):
                m.update(urllib.urlencode(arg))
            elif isinstance(arg, six.string_types):
                m.update(arg)

        return m.hexdigest()

The implemention has to sort and encode query string to generate a proper signature.

The signature is passed to the application by appending the ``sig``
parameter to the query string; e.g.
``w=100&h=100&sig=c9516346abf62876b6345817dba2f9a0c797ef26``.

Note, the application does not include the leading question mark when verifying
the supplied signature. To verify your signature implementation, see the
``signature`` command described in the `Tools`_ section.

Limiting allowed sizes
----------------------

Depending on your use case it may be more appropriate to simply restrict the
image sizes picfit is allowed to generate. See the `Allowed sizes`_ section for
more information on this configuration.

Tools
=====

To verify that your client application is generating correct signatures,
use the command::

    $ picfit signature --key=abcdef "w=100&h=100&op=resize"
    Query String: w=100&h=100&op=resize
    Signature: 6f7a667559990dee9c30fb459b88c23776fad25e
    Signed Query String: w=100&h=100&op=resize&sig=6f7a667559990dee9c30fb459b88c23776fad2

Error reporting
===============

picfit logs events by default in ``stderr`` and ``stdout``. You can implement sentry_
to log errors using raven_.

To enable this feature, set ``sentry`` option in your config file.

``config.json``

.. code-block:: json

    {
      "sentry": {
        "dsn": "[YOUR_SENTRY_DSN]",
        "tags": {
          "foo": "bar"
        }
      }
    }

Debug
=====

Debug is disabled by default.

To enable this feature set ``debug`` option to ``true`` in your config file:

``config.json``

.. code-block:: json

    {
      "debug": true
    }

CORS
====

picfit supports CORS headers customization in your config file.

To enable this feature, set ``allowed_origins``, ``allowed_headers`` and ``allowed_methods``,
for example:

``config.json``

.. code-block:: json

    {
      "allowed_headers": ["Content-Type", "Authorization", "Accept", "Accept-Encoding", "Accept-Language"],
      "allowed_origins": ["*.ulule.com"],
      "allowed_methods": ["GET", "HEAD"]
    }

Image engine
============

Quality
-------

The quality rendering of the image engine can be controlled
globally without adding it at each request:

``config.json``

.. code-block:: json

    {
      "engine": {
        "quality": 70
      }
    }

With this option, each image will be saved in ``70`` quality.

By default the quality is the highest possible: ``95``

Format
------

The format can be forced globally without adding it at each request:

``config.json``

.. code-block:: json

    {
      "engine": {
        "format": "png"
      }
    }

With this option, each image will be forced to be saved in ``.png``.

By default the format will be chosen in this order:

* The ``fmt`` parameter if exists in query string
* The original image format
* The default format provided in the `application <https://github.com/thoas/picfit/blob/master/application/constants.go#L6>`_

Options
=======

Deletion
--------

Deletion is disabled by default for security reason, you can enable
it in your config:

``config.json``

.. code-block:: json

    {
      "options": {
        "enable_delete": true
      }
    }

You will be able to delete root image and its children, for example if you upload an image with
the file path `/foo/bar.png`, you can delete the main image on stockage by sending the following HTTP request:


::

   DELETE https://localhost:3001/foo/bar.png

or delete a child:

::

   DELETE https://localhost:3001/display/thumbnail/100x100/foo/bar.png

If you want to delete the main image and cascade its children, you can enable it in your config:

``config.json``

.. code-block:: json

    {
      "options": {
        "enable_delete": true,
        "enable_cascade_delete": true
      }
    }

when a new image will be processed, it will be linked to the main image and stored in the kvstore.

Upload
------

Upload is disabled by default for security reason, you can enable
it in your config:

``config.json``

.. code-block:: json

    {
      "options": {
        "enable_upload": true
      }
    }

Stats
-----

Stats are disabled by default, you can enable them in your config.

``config.json``

.. code-block:: json

    {
      "options": {
        "enable_stats": true
      }
    }

It will store various information about your web application (response time, status code count, etc.).

To access these information, you can visit: http://localhost:3001/sys/stats

Health
------

Health is disabled by default, you can enable it in your config.

``config.json``

.. code-block:: json

    {
      "options": {
        "enable_stats": true
      }
    }

It will show various internal information about the Go runtime (memory, allocations, etc.).

To access these information, you can visit: http://localhost:3001/sys/health

Profiler
--------

Profiler is disabled by default, you can enable it in your config.

``config.json``

.. code-block:: json

    {
      "options": {
        "enable_pprof": true
      }
    }

It will start pprof_ then use the pprof tool to look at the heap profile:

::

   go tool pprof http://localhost:3001/debug/pprof/heap

Or to look at a 30-second CPU profile:

::

   go tool pprof http://localhost:3001/debug/pprof/profile

Or to look at the goroutine blocking profile, after calling runtime.SetBlockProfileRate in your program:

::

   go tool pprof http://localhost:3001/debug/pprof/block

Or to collect a 5-second execution trace:

::

   wget http://localhost:3001/debug/pprof/trace?seconds=5

Logging
-------

By default the logger level is `debug`, you can change it in your config:

``config.json``

.. code-block:: json

    {
      "logger": {
        "level": "info"
      }
    }

Levels available are:

* debug
* info
* error
* warning
* fatal

Allowed sizes
-------------

To restrict the sizes picfit is allowed to generate you may specify the
``allowed_sizes`` option as an array of sizes. Note that if you omit a width or
height from a size it will allow requests that exclude height or width to
preserve aspect ratio.

``config.json``

.. code-block:: json

    {
      "options": {
        "allowed_sizes": [
          {"width": 1920, "height": 1080},
          {"width": 720, "height": 480},
          {"width": 480}
        ]
      }
    }

IP Address restriction
----------------------

You can restrict access to upload, stats, health, delete and pprof endpoints by enabling
restriction in your config:

``config.json``

.. code-block:: json

    {
      "options": {
        "allowed_ip_addresses": [
          "127.0.0.1"
        ]
      }
    }

Deployment
==========

It's recommended that the application run behind a CDN for larger applications
or behind varnish for smaller ones.

Provisioning is handled by Ansible_, you will find files in
the `repository <https://github.com/thoas/picfit/tree/master/provisioning>`_.

You must have Ansible_ installed on your laptop, basically if you have python
already installed you can do ::

    $ pip install ansible

Roadmap
=======

see `issues <https://github.com/thoas/picfit/issues>`_

Don't hesitate to send patch or improvements.


Clients
=======

Client libraries will help you generate picfit urls with your secret key.

* `picfit-go <https://github.com/ulule/picfit-go>`_: a Go client library

In production
=============

- Ulule_: an european crowdfunding platform

Inspirations
============

* pilbox_
* `thumbor <https://github.com/thumbor/thumbor>`_
* `trousseau <https://github.com/oleiade/trousseau>`_

Thanks to these beautiful projects.

.. _GOPATH: http://golang.org/doc/code.html#GOPATH
.. _Redis: http://redis.io/
.. _Redis cluster: https://redis.io/topics/cluster-tutorial
.. _pilbox: https://github.com/agschwender/pilbox
.. _varnish: https://www.varnish-cache.org/
.. _persistence: http://redis.io/topics/persistence
.. _Ansible: http://www.ansible.com/home
.. _Ulule: http://www.ulule.com
.. _sentry: https://github.com/getsentry/sentry
.. _raven: https://github.com/getsentry/raven-go
.. _httpie: https://github.com/jakubroztocil/httpie
.. _pprof: https://blog.golang.org/profiling-go-programs
