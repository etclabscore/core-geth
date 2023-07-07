---
hide:
  - toc        # Hide table of contents
title: Publishing a Release
---

# Developers: How to Make a Release

- [ ] Decide what the new version should be. In this example, __`v1.11.16[-stable]`__ will be used.
- [ ] `git checkout master`
- [ ] `make lint` and `make test` are passing on master. :white_check_mark:
  > This is important because the artifacts to be included with the release will be generated
  by the CI workflows. If linting or tests fail, the workflows will be interrupted
  and artifacts will not be generated.
- [ ] `git checkout release/v1.11.16`
- [ ] Edit `params/version.go` making the necessary changes to version information. (To `-stable` version.) _Gotcha:_ make sure this passes linting, too.
- [ ] `git commit -S -s -m "bump version from v1.11.16-unstable to v1.11.16-stable"`
- [ ] `git tag -S -a v1.11.16`
- [ ] `git push etclabscore v1.11.16`
  > Push the tag to the remote. I like to do it this way because it triggers the tagged version on CI before the branch/PR version,
  expediting artifact delivery.
- [ ] Edit `params/version.go` making the necessary changes to version information. (To `-unstable` version.)
- [ ] `git commit -S -s -m "bump version from v1.11.16-stable to v1.11.17-unstable"`
- [ ] `git push etclabscore`
  > Push the branch. This will get PR'd, eg. https://github.com/etclabscore/core-geth/pull/197
- [ ] Draft a new release, following the existing patterns for naming and notes. https://github.com/etclabscore/core-geth/releases/new
    - Define the tag the release should be associated with (eg `v1.11.16`).
    - Linux, OSX, and Windows artifacts will be uploaded automatically to this release draft by the CI jobs. There should be CI-generated 34 assets total.

        !!! Note

            If the release is not drafted manually, it will be automatically drafted by the CI.

- [ ] Await a complete set of uploaded artifacts. If artifacts fail to upload due to issue with the CI jobs, review
  those jobs to determine if their failure(s) is OK, restarting them if so.
- [ ] Once artifacts have been uploaded and the release draft reviewed by one other person for the following, it's time to publish!
    + proofreading
    + artifact fingerprint verification
    + notes content approval
- [ ] Once the release is published, merge the associated PR bumping versions.
