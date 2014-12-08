picfit
======

picfit is a reusable Go server to manipulate (resizing, croping, etc.) images built
on top of `negroni <https://github.com/codegangsta/negroni>`_ and `gorilla mux <https://github.com/gorilla/mux>`_.

It will act as a proxy on top of your storage engine and served ideally behind an http cache system like varnish_.

Installation
============

Build it
--------

1. Make sure you have a Go language compiler >= 1.3 (mandatory) and git installed.
2. Make sure you have the following go system dependencies in your $PATH: bzr, svn, hg, git
3. Ensure your GOPATH_ is properly set.
4. Download picfit::

    git clone https://github.com/thoas/picfit.git

4. Run ``make build``

You have now a binary version of picfit in the ``bin`` directly which fits perfectly with your architecture.

Debian and Ubuntu
-----------------

We will provide Debian package when we will be completely stable ;)

Configuration
=============

picfit only accepts configuration in JSON format, this configuration should be stored in a file.

Basic
-----

* no key/value store
* no image storage
* images are given in absolute url

Images are processed on the fly at each requests

``config.json``

.. code-block:: json

    {
      "port": 3001,
    }

Store images on file system and keys in an in-memory cache
----------------------------------------------------------

* key/value in-memory store
* file system storage

An image is processed and uploaded asynchronously to the storage.

An unique key is generated and stored in a in-memory key/value store to process
a dedicated request one time.

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
---------------------------------------------------------------

* key/value store provided by Redis
* Amazon S3 storage
* shard filename

.. code-block:: json

    {
      "kvstore": {
        "type": "redis",
        "host": "127.0.0.1",
        "port": "6379",
        "password": "",
        "db": 0
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
        "depth": 2
      }
    }

With this config, we will store keys on Redis_.

Images will be stored on Amazon S3 at the location ``/path/to/directory``.

``[ACL]`` can be:

* private
* public-read
* public-read-write
* authenticated-read
* bucket-owner-read
* bucket-owner-full-control

``[REGION_NAME]`` can be:

* us-gov-west-1
* us-east-1
* us-west-1
* us-west-2
* eu-west-1
* eu-central-1
* ap-southeast-1
* ap-southeast-2
* ap-northeast-1
* sa-east-1
* cn-north-1

**Filename** will be sharded:

* ``depth``: 2 directories
* ``width``: 1 letter for each directory

Example:

``06102586671300cd02ae90f1faa16897.png`` will become ``0/6/102586671300cd02ae90f1faa16897.jpg``

Load images from file system and store them in Amazon S3, keys on Redis
=======================================================================

* key/value store provided by Redis
* File system to load images
* Amazon S3 storage to process images

.. code-block:: json

    {
      "kvstore": {
        "type": "redis",
        "host": "127.0.0.1",
        "port": "6379",
        "password": "",
        "db": 0
      },
      "port": 3001,
      "storage": {
        "src": {
          "type": "fs",
          "location": "path/to/directory"
        },
        "dest": {
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

With this config, you can load and store your images from different storage backends.

Running
=======

To run the application, issue the following command::

    $ picfit config.json

By default, this will run the application on port 8888 and can be accessed by visiting:::

    http://localhost:3001

To see a list of all available options, run::

    $ picfit --help

Usage
=====

Format
------

The format to call the service is ::

    <img src="http://localhost:3001/{method}?url={url}&path={path}&w={width}&h={height}&upscale={upscale}&sig={sig}&op={operation}"

- *path*: The filepath to load the image using your source storage
- *operation*: The method to perform (``resize``, ``thumbnail``)
- *sig*: The signature key which is the representation of your query string and your secret key
- *method*: The operation to perform (``get``, ``display``)
- *url*: The url of the image to be processed (not required if **filepath** provided)
- *width*: The desired width of the image, if ``0`` is provided the service will calculate the ratio with **height**
- *height*: The desired height of the image, if ``0`` is provided the service will calculate the ratio with **width**
- *upscale*: If your image is smaller than your desired dimensions, the service will upscale by default to fit your dimensions, you can disable this behavior by providing ``0``.

To use this service, include the service url as replacement for your images, for example:::

    <img src="https://www.google.fr/images/srpr/logo11w.png" />

will become::

    <img src="http://localhost:3001/display?url=https%3A%2F%2Fwww.google.fr%2Fimages%2Fsrpr%2Flogo11w.png&w=100&h=100&op=resize"

This will request the image served at the supplied url and resize it to 100x100 using the **resize** method.

Using source storage
--------------------

If an image is stored in your source storage at the location ``path/to/file.png``, then you can call the service
to load this file::

    <img src="http://localhost:3001/display?w=100&h=100&path=path/to/file.png&op=resize"


Security
========

In order to secure requests so that unknown third parties cannot easily
use the resize service, the application can require that requests
provide a signature. To enable this feature, set the ``secret_key``
option in your config file.

The signature is a hexadecimal digest generated from the client
key and the query string using the HMAC-SHA1 message authentication code
(MAC) algorithm. The below python code provides an example
implementation.

::

    import hashlib
    import hmac
    import json
    import six

    def sign(key, *args, **kwargs):
        m = hmac.new(key, None, hashlib.sha1)

        for arg in args:
            if arg is None:
                continue
            elif isinstance(arg, dict):
                m.update(json.dumps(arg))
            elif isinstance(arg, six.string_types):
                m.update(arg)

        return m.hexdigest()

The signature is passed to the application by appending the ``sig``
parameter to the query string; e.g.
``w=100&h=100&sig=c9516346abf62876b6345817dba2f9a0c797ef26``.

Note, the application does not include the leading question mark when verifying
the supplied signature. To verify your signature implementation, see the
``signature`` command described in the `Tools`_ section.

Tools
=====

To verify that your client application is generating correct signatures, use the signature command.

::
    $ picfit signature --key=abcdef "w=100&h=100&op=resize"
    Query String: w=100&h=100&op=resize
    Signature: 6f7a667559990dee9c30fb459b88c23776fad25e
    Signed Query String: w=100&h=100&op=resize&sig=6f7a667559990dee9c30fb459b88c23776fad2

Deployment
==========

...

Inspirations
============

* `pilbox <https://github.com/agschwender/pilbox>`_
* `thumbor <https://github.com/thumbor/thumbor>`_
* `trousseau <https://github.com/oleiade/trousseau>`_

Thanks to them, beautiful projects.

.. _GOPATH: http://golang.org/doc/code.html#GOPATH
.. _Redis: http://redis.io/
.. _varnish: https://www.varnish-cache.org/
