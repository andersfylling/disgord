# Description
<!--
Please include a summary of the PR. Do not create a PR into branch:develop unless there exist no branch for the next release.

eg. If the current release is v0.10, then you should create a PR for branch:release/v0.11. If the next release branch is not out/created yet, create an issue or make a draft PR that goes into develop and change it to the release branch later on. Don't worry, I'll try my best to make this easy for everyone.
-->

## Type of change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)

## Benchmarks
<!--
If this PR requires benchmarks (say it is an very dependent component or takes a lot of resources/use, use pprof if you need to) then the benchmarks are provided before and after such that we can make logical decisions.
Note that if you add a benchmark and find your solution to run slower, the code might still be valuable so your results are welcomed anyways!
If no benchmarks are needed, feel free to delete til paragraph.
-->

# Checklist:

- [ ] I ran `go generate`
- [ ] I ran `go fmt ./...`
- [ ] I have performed a self-review of my own code
- [ ] Commented complex situations or referenced the discord documentation
- [ ] Updated documentation
- [ ] Added/Updated unit tests
- [ ] Added/Updated benchmarks (if this is a performance critical component)
