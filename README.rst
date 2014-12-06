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

Config:

* no key/value store
* no image storage
* images are given in absolute url

Images are processed on the fly at each requests

``config.json``

.. code-block:: json

    {
      "port": 3001,
    }

Shard files on file system and store ids on Redis
-------------------------------------------------

Config:

* key/value store
* file system storage
* shard filename

An image is processed and uploaded asynchronously to the storage.

An unique key is generated and stored in our key/value store to process
a dedicated request one time.

The image's filename is also sharded into multiple directories.

.. code-block:: json

    {
      "redis": {
        "host": "127.0.0.1",
        "port": "6379",
        "password": "",
        "db": 0
      },
      "port": 3001,
      "storage": "fs",
      "fs": {
        "location": "/path/to/directory/"
      }
      "shard": {
        "width": 1,
        "depth": 2
      }
    }

With the following config, we will store ids on Redis_ and store the image file
on the file system at the location ``/path/to/directory``.

Filename will be sharded in 2 directories (``depth``) with 1 letter for each (``width``):

``06102586671300cd02ae90f1faa16897.png`` will become ``0/6/102586671300cd02ae90f1faa16897.jpg``.

Inspirations
============

* `pilbox <https://github.com/agschwender/pilbox>`_
* `thumbor <https://github.com/thumbor/thumbor>`_

.. _GOPATH: http://golang.org/doc/code.html#GOPATH
.. _Redis: http://redis.io/
