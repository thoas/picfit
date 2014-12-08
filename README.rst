picfit
======

picfit is a reusable Go server to manipulate (resizing, croping, etc.) images built
on top of `negroni <https://github.com/codegangsta/negroni>`_ and `gorilla mux <https://github.com/gorilla/mux>`_.

Installation
============

Build it
--------

1. Make sure you have a Go language compiler >= 1.3 (mandatory) and git installed.
2. Make sure you have the following go system dependencies in your $PATH: bzr, svn, hg, git
3. Ensure your GOPATH_ is properly set.
4. Run ``make``

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
        "type": "fs",
        "location": "/path/to/directory/"
      },
      "kvstore": {
        "type": "cache"
      },
    }

Store images on Amazon AWS S3, keys in Redis and shard filename
---------------------------------------------------------------

* key/value store provided by Redis
* Amazon AWS S3 storage
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
        "source": {
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

With the following config, we will store keys on Redis_.

Images will be stored on Amazon AWS S3 at the location ``/path/to/directory``.

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

Running
=======

To run the application, issue the following command::

    $ picfit config.json

By default, this will run the application on port 8888 and can be accessed by visiting:::

    http://localhost:3001

To see a list of all available options, run::

    $ picfit --help

Calling
=======

...

Security
========

...

Tools
=====

...

Deployment
==========

...

Inspirations
============

* `pilbox <https://github.com/agschwender/pilbox>`_
* `thumbor <https://github.com/thumbor/thumbor>`_
* `trousseau <https://github.com/oleiade/trousseau>`_

.. _GOPATH: http://golang.org/doc/code.html#GOPATH
.. _Redis: http://redis.io/
