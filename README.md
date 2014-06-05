Golang Client
=============
stackforge/golang-client is yet another implementation of [OpenStack]
(http://www.openstack.org/) API client in [Go language](http://golang.org).
The code follows OpenStack licensing and borrows its infrastructure for code
hosting.  It currently implements [Identity Service v2] 
(http://docs.openstack.org/api/openstack-identity-service/2.0/content/) 
and [Object Storage v1] 
(http://docs.openstack.org/api/openstack-object-storage/1.0/content/).
Some API calls are not implemented initially, but the intention is to expand
the lib over time (where pragmatic).

Code maturity is considered experimental.

Installation
------------
Use `go get git.openstack.org/stackforge/golang-client.git`.  Or alternatively,
download or clone the repository.

The lib was developed and tested on go 1.2. No external dependencies, so far.

Examples
--------
The examples directory contains examples for using the SDK using
real world working code. Each example starts with a two digit number followed
by a name (e.g., `00-authentication.go`). If you have a `config.json` file in the
examples directory following the format of `config.json.dist` the example can be
executed using `go run [example name] setup.go`. Or, all the examples can be
executed running the script `run-all.sh` from the examples directory.

Testing
-------
There are two types of test files.  The `*_test.go` are standard
golang unit test files.  The examples can be run as integration tests.

The tests were written against the [OpenStack API specifications]
(http://docs.openstack.org/api/api-specs.html).
The integration test were successful against the following:

- [HP Helion Public Cloud](http://docs.hpcloud.com/api/)

If you use another provider and successfully completed the tests, please email
the maintainer(s) so your service can be mentioned here.  Alternatively, if you
are a service provider and can arrange a free (temporary) account, a quick test
can be arranged.

License
-------
Apache v2.

Contributing
------------
The code repository borrows OpenStack StackForge infrastructure.
Please use the [recommended workflow]
(https://wiki.openstack.org/wiki/GerritWorkflow).  If you are not a member yet,
please consider joining as an [OpenStack contributor]
(https://wiki.openstack.org/wiki/HowToContribute).  If you have questions or
comments, you can email the maintainer(s).

Coding Style
------------
The source code is automatically formatted to follow `go fmt` by the [IDE]
(https://code.google.com/p/liteide/).  And where pragmatic, the source code
follows this general [coding style]
(http://slamet.neocities.org/coding-style.html).