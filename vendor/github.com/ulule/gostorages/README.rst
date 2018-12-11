gostorages
==========

A unified interface to manipulate storage engine for Go.

gostorages is used in `picfit <https://github.com/thoas/picfit>`_ to allow us
switching over storage engine.

Currently, it supports the following storages:

* Amazon S3
* File system

Installation
============

Just run:

::

    $ go get github.com/ulule/gostorages

Usage
=====

It offers you a single API to manipulate your files on multiple storages.

If you are migrating from a File system storage to an Amazon S3, you don't need
to migrate all your methods anymore!

Be lazy again!

File system
-----------

To use the ``FileSystemStorage`` you must have a location to save your files.

.. code-block:: go

    package main

    import (
        "fmt"
        "github.com/ulule/gostorages"
        "os"
    )

    func main() {
        tmp := os.TempDir()

        storage := gostorages.NewFileSystemStorage(tmp, "http://img.example.com")

        // Saving a file named test
        storage.Save("test", gostorages.NewContentFile([]byte("(╯°□°）╯︵ ┻━┻")))

        fmt.Println(storage.URL("test")) // => http://img.example.com/test

        // Deleting the new file on the storage
        storage.Delete("test")
    }


Amazon S3
---------

To use the ``S3Storage`` you must have:

* An access key id
* A secret access key
* A bucket name
* Give the region of your bucket
* Give the ACL you want to use

You can find your credentials in `Security credentials <https://console.aws.amazon.com/iam/home?nc2=h_m_sc#security_credential>`_.

In the following example, I'm assuming my bucket is located in european region.

.. code-block:: go

    package main

    import (
        "fmt"
        "github.com/ulule/gostorages"
        "github.com/mitchellh/goamz/aws"
        "github.com/mitchellh/goamz/s3"
        "os"
    )

    func main() {
        baseURL := "http://s3-eu-west-1.amazonaws.com/my-bucket"

        storage := gostorages.NewS3Storage(os.Getenv("ACCESS_KEY_ID"), os.Getenv("SECRET_ACCESS_KEY"), "my-bucket", "", aws.Regions["eu-west-1"], s3.PublicReadWrite, baseURL)

        // Saving a file named test
        storage.Save("test", gostorages.NewContentFile([]byte("(>_<)")))

        fmt.Println(storage.URL("test")) // => http://s3-eu-west-1.amazonaws.com/my-bucket/test

        // Deleting the new file on the storage
        storage.Delete("test")
    }

Roadmap
=======

see `issues <https://github.com/ulule/gostorages/issues>`_

Don't hesitate to send patch or improvements.
