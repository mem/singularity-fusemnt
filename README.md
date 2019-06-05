Singularity FUSE mount plugin
=============================

This is an example plugin implementing a singularity plugin to mount
FUSE filesystems.

Building
--------

In order to build the plugin you need a copy of code matching the
version of singularity that you wish to use. You can find the commit
matching the singularity binary by running:

    $ singularity version
    3.1.1-723.g7998470e7

this means this version of singularity is _post_ 3.1.1 (but before the
next version after that one). The suffix .gXXXXXXXXX indicates the exact
commit in github.com/sylabs/singularity used to build this binary
(7998470e7 in this example).

Obtain a copy of the source code by running:

    git clone https://github.com/sylabs/singularity.git
    cd singularity
    git checkout 7998470e7
    git clone https://github.com/mem/singularity-fusemnt singularity-fusemnt

Still from within the singularity directory, run:

	singularity plugin compile -o singularity-fusemnt.sif ./singularity-fusemnt

This will produce a file called ./singularity-fusemnt.sif

Installing
----------

Once you have compiled the plugin into a SIF file, you can install it
into the correct singularity directory using the command:

	$ singularity plugin install ./singularity-fusemnt.sif

Singularity will automatically load the plugin code from now on.

Using
-----

This plugin adds a flag `--fusemnt` which takes the mount point _inside_
the container as an argument. You _must_ provide the `fuse-example`
program. You can use the [hello
program](https://github.com/libfuse/libfuse/blob/master/example/hello.c)
provided with FUSE in the examples directory.
