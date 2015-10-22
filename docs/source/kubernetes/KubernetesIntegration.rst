Integrating the Calico Network Plugin with Kubernetes
=====================================================

This guide will describe the configuration required to use the Calico
network plugin in your Kubernetes deployment.

Setting Up Your Environment
---------------------------

The Calico network plugin looks for five environment variables. If one
is not set, it will assume a default value. If you need to override the
defaults, you must set these variables in the environment of the
*kubelet* process.

-  ETCD\_AUTHORITY
   '''''''''''''''

   By default, the Calico network plugin will assume that the etcd
   datastore is located at ``<MASTER_IP>:6666``. Setting the
   ``ETCD_AUTHORITY`` variable in your environment will direct Calico to
   the correct IP if your cluster is set up differently.

-  CALICOCTL\_PATH
   '''''''''''''''

   This plugin requires access to the ``calicoctl`` binary. If your
   binary is not located at ``/usr/bin/calicoctl``, set the
   ``CALICOCTL_PATH`` environment variable to the correct path.

-  KUBE\_API\_ROOT
   '''''''''''''''

   The ``KUBE_API_ROOT`` environment variable specifies the URL for the
   root of the Kubernetes API. The transport must be included. This
   variable defaults to ``http://kubernetes-master:8080/api/v1/``.

-  DEFAULT\_POLICY (added in calico-kubernetes v0.2.0)
   '''''''''''''''''''''''''''''''''''''''''''''''''''

   The ``DEFAULT_POLICY`` environment variable applies `security
   policy <http://docs.projectcalico.org/en/latest/security-model.html>`__
   to a set of pods. The default policy in the Calico network plugin is
   ``allow``, which allows all incoming and outgoing traffic to and from
   a pod. Alternately, you may also indicate ``ns_isolation``, which
   will only allow incoming traffic from pods of the same namespace and
   allow all outgoing traffic.

-  CALICO\_IPAM (added in calico-kubernetes v0.2.0)
   ''''''''''''''''''''''''''''''''''''''''''''''''

   The ``CALICO_IPAM`` environment variable gives the option to utilize
   Calico IP Address Management (IPAM). When set to ``true``, Calico
   will automatically assign pods an IP address that is unique in the
   cluster. By default, ``CALICO_IPAM`` is set to ``false``, and pods
   utilize the IP address assigned by Docker.

-  KUBE\_AUTH\_TOKEN (added in calico-kubernetes v0.3.0)
   '''''''''''''''''''''''''''''''''''''''''''''''''''''

   The ``KUBE_AUTH_TOKEN`` environment variable specifies the token to
   use for https authentication with the Kubernetes apiserver. Each
   Kubernetes Service Account has its own API token. You can create
   Service Accounts by following the instructions in the `Kubernetes
   docs <http://kubernetes.io/v1.0/docs/user-guide/service-accounts.html>`__.

Configuring Nodes
-----------------

Creating a Calico Node with the Network Plugin
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

-  Automatic Install
   '''''''''''''''''

As of Calico v0.5.1, we have included a ``--kubernetes`` flag to the
``calicoctl node`` command that will automatically install the Calico
Network Plugin as you spin up a Calico Node.
``sudo ETCD_AUTHORITY=<ETCD_IP>:<ETCD_PORT> calicoctl node --ip=<NODE_IP> --kubernetes``
>\ *Note in this example, we set the ETCD*\ AUTHORITY environment config
for the duration of the command\_

-  Manual Install
   ''''''''''''''

Alternatively, you can download the `latest
release <https://github.com/projectcalico/calico-docker/releases/latest>`__
of the plugin binary directly from our GitHub Repo.

Configuring Kubelet Services
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

On each of your nodes, you will need to configure the Kubelet to use the
Calico Networking Plugin. This can be done by including the
``--network_plugin=calico`` option when starting the Kubelet. If you are
using systemd to manage your services, you can add this line to the
Kubelet config file (``/etc/systemd/`` by default) and restart your
Kubelets to begin using Calico.

Configuring Policy
~~~~~~~~~~~~~~~~~~

See our doc on `Programming Kubernetes Policy <KubernetesPolicy.md>`__
to start enforcing security policy on Kubernetes pods!
