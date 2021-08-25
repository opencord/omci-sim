# DEPRECATED

[![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech/)

`VOLTHA 2.6` was the last official set of releases that  used the `omci-sim` repository in the openonu-adapter in python.
From `2.5` onwards VOLTHA uses component a Golang version of the openonu adapter, thus the repo has been deprecated, the go equivalent is [omci-lib-go](https://github.com/opencord/omci-lib-go). This codebase is going to be removed after the VOLTHA 2.8 release LTS support ends in December 2022.

# Omci-sim

This library is a dump of the OMCI messages reported by an ALPHA device
and it is currently used by [BBSim](github.com/opencord/bbsim) to 
emulate OMCI responses during [VOLTHA](docs.voltha.org) scale tests.
