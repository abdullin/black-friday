This is a black friday experiment code: implementing a subset of warehouse management system with event sourcing.

This is an implementation of [Trustbit Black Friday Kata](https://github.com/trustbit/bfkata).

This experiment focuses:

- event sourcing;
- domain-driven design;
- event-driven tests (specs).

Usage:


```bash

# rebuild the binary
make clean bin/bf

# run tests
bin/bf test

# use tests to benchmark performance
bin/bf perf

# run spec N47 and keep database state intact for exploration
bin/bf explore --spec 47


# run this implementation against bfkata test suite
bin/bf subject
```


You can read more in:

- [Report 1](https://abdullin.com/post/black-friday/experiment/)
- [Better performance with DOD](https://abdullin.com/post/black-friday/data-oriented-design/)


The latest branch doesn't contain DOD optimisations. These are in the [locs branch](https://github.com/abdullin/black-friday/tree/locs/inventory/mem)


