# Description

Please include a summary of the change and which issue is fixed.

Use smart commits here to manipulate issues (eg. Fixes #issue)

## Type of change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)

## Benchmarks
If this PR requires benchmarks (say it is an very dependent component or takes a lot of resources/use, use pprof if you need to) then the benchmarks are provided before and after such that we can make logical decisions.
Note that if you add a benchmark and find your solution to run slower, the code might still be valuable so your results are welcomed anyways!
If no benchmarks are needed, feel free to delete til paragraph.

# Checklist:

- [ ] I ran `go generate`
- [ ] I have performed a self-review of my own code (remember to run `go fmt ./...`)
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] Any dependent changes have been merged and published in downstream modules
- [ ] Added benchmarks if this is a performant required component (potential bottlenecks)
