# Note to Code Shelter maintainers

The `multi-geth` project is an extension of upstream
[go-ethereum](https://github.com/ethereum/go-ethereum) that add more
blockchain network support. Some philosophical notes to maintainers:

* We accept all patches to add new network supports, as long as it's
  reasonable. In other words, we don't reject new network supports due
  to political reasons.
* We follow upstream closely.

Besides checking out new PR to merge, please also help out following
upstream and making new releases. Thank you!

## Merge Upstream Changes

Add `https://github.com/ethereum/go-ethereum` to one of your
remote. Then pull the upstream `master` branch (note that it must be
`master` branch!), resolve conflicts and create a new PR.

When you review a merge-upstream PR, please make sure it is merged via
a merge commit.

## Making New Releases

Whenever the upstream makes a new release, multi-geth usually makes a
new release as well. Besides, new releases can also be made when a
notable new feature is added in multi-geth.

To create a new release, branch off `master` into a release branch
`release/vx.x.x`. Change `VersionMeta` to `stable` in
`params/version.go`. Create a **signed commit** and push it to this
repo. The CI will them automatically build all release binaries. After
all binaries are built, create a Github release on that tag.
