# Branch Cleanup Action

A GitHub action to automatically label all PRs according to file patterns.

**Table of Contents**

<!-- toc -->

- [Usage](#usage)

<!-- tocstop -->

## Usage

Action is meant to be run as periodic job. This is needed to workaround issues regarding
[lack of write access when executed from fork](https://help.github.com/en/actions/automating-your-workflow-with-github-actions/authenticating-with-the-github_token#permissions-for-the-github_token)
which is a common problem when using https://github.com/actions/labeler.

```
---
name: Pull request labeler
on:
  schedule:
    - cron: '*/5 * * * *'
jobs:
  labeler:
    runs-on: ubuntu-latest
    steps:
      - uses: paulfantom/periodic-labeler@master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          GITHUB_REPOSITORY: ${{ github.repository }}
          LABEL_MAPPINGS_FILE: .github/labeler.yml
```

By default action uses `.github/labeler.yml` located in repository from `GITHUB_REPOSITORY` as a source of pattern matchers.
This file uses the same schema as in https://github.com/actions/labeler
