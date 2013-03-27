hyperloglog
===========

Package hyperloglog implements the HyperLogLog algorithm for
cardinality estimation. In English: it counts things. It counts things
using very small amounts of memory compared to the number of objects
it is counting.

For a full description of the algorithm, see the paper HyperLogLog:
the analysis of a near-optimal cardinality estimation algorithm by
Flajolet, et. al. at http://algo.inria.fr/flajolet/Publications/FlFuGaMe07.pdf

For documentation see http://godoc.org/github.com/eclesh/hyperloglog

Quick start
===========

	$ go get github.com/eclesh/hyperloglog
	$ go test -test.v

Tests take about a minute and don't have a failure case. Sorry.

License
=======

hyperloglog is licensed under the MIT license.
