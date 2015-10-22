Preparing the environment for libnetwork
========================================

The libnetwork Calico demonstration is run on two Linux servers that
have a number of installation requirements.

This document describes how to get a working environment for the
demonstration.

Follow instructions using either the `Automated
setup <#automated-setup>`__ for fast setup in a virtualized environment,
or `Manual Setup <#manual-setup>`__ if, for example, you are installing
on bare metal servers or would like to understand the process in more
detail.

Once you have the environment set up, you can run through the
`demonstration <Demonstration.md>`__.

Automated setup
---------------

The worked examples below allow you to quickly set up a virtualized
environment using Vagrant or a cloud service - be sure to follow the
appropriate instructions for using libnetwork.

-  `Vagrant install with CoreOS <../VagrantCoreOS.md>`__
-  `Vagrant install with Ubuntu <../VagrantUbuntu.md>`__
-  `Amazon Web Services (AWS) <../AWS.md>`__
-  `Google Compute Engine (GCE) <../GCE.md>`__
-  `DigitalOcean <../DigitalOcean.md>`__

Manual Setup
============

Summary
-------

The two Linux servers require a recent version of Linux with a few
network network related kernel modules to be loaded. The easiest way to
ensure this and that the servers meet the requirements is to run
``calicoctl checksystem``

The servers also need: - A specific Docker release to be running - since
the Calico agent is packaged as a Docker container, and the libnetwork
features required are currently only available in an experimental
release. - A consul server used for clustering Docker - An Etcd cluster
- which Calico uses for coordinating state between the nodes. - The
``calicoctl`` binary to be placed in the system ``$PATH``.

Requirements
------------

For the demonstration, you just need 2 servers (bare metal or VMs) with
a modern 64-bit Linux OS and IP connectivity between them.

We recommend configuring the hosts with the hostname ``calico-01`` and
``calico-02``. The demonstration will refer to these hostnames.

They must have the following software installed: - `Docker 1.9 or
greater <#Docker>`__ - etcd installed and available on each node: `etcd
documentation <https://coreos.com/etcd/docs/latest/>`__ - ``ipset``,
``iptables``, and ``ip6tables`` kernel modules.

Docker
~~~~~~

Follow the instructions for installing
`Docker <https://docs.docker.com/installation/>`__.

A version of 1.9 or greater is required. At the current time, the 1.9
release is finishing development, however you can download the binaries
from the overnight builds of the master branch (1.9.dev) from
https://master.dockerproject.org/linux/amd64

Docker permissions
~~~~~~~~~~~~~~~~~~

Running Docker is much easier if your user has permissions to run Docker
commands. If your distro didn't set this permissions as part of the
install, you can usually enable this by adding your user to the
``docker`` group and restarting your terminal.

::

    sudo usermod -aG docker <your_username>

If you prefer not to do this you can still run the demo but remember to
run ``docker`` using ``sudo``.

Getting calicoctl Binary
~~~~~~~~~~~~~~~~~~~~~~~~

Get the calicoctl binary onto each host.

::

    wget http://www.projectcalico.org/latest/calicoctl
    chmod +x calicoctl

This binary should be placed in your ``$PATH`` so it can be run from any
directory.

.. raw:: html

   <!--- master only -->

    Note that projectcalico.org is not an HA repository, so using this
    download URL is not recommended for any automated production
    installation process. Also, the latest calicoctl refers to builds
    from the master branch which is under active development.

    For integration and real deployment, we recommend you follow the
    instructions for a specific release.

    Follow the documents for the latest release. Alternatively, view the
    documentation for the relevant release tag.

.. raw:: html

   <!--- end of master only -->

Preload the Calico docker image (optional)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

You can optionally preload this image to avoid the delay when you run
``calicoctl node --libnetwork`` the first time. For example, to pull the
latest released version, run

::

    docker pull calico/node:latest
    docker pull calico/node-libnetwork:latest

Final checks
------------

Verify the hostnames. If they don't match the recommended names above
then you'll need to adjust the demonstration instructions accordingly.

Check that the hosts have IP addresses assigned, and that your hosts can
ping one another.

Check that you are running with a suitable version of Docker.

::

    docker version

It should indicate a version of 1.9 or greater.

You should also verify each host can access etcd. The following will
return the current etcd version if etcd is available.

::

    curl -L http://127.0.0.1:4001/version

