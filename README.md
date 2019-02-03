
# mrcli-go - GitLab pending merge requests overview [![Build Status](https://img.shields.io/travis/danielpieper/mrcli-go.svg?branch=master&style=flat-square)](https://travis-ci.org/danielpieper/mrcli-go?branch=master)

[![Scrutinizer Coverage](https://img.shields.io/scrutinizer/coverage/g/danielpieper/mrcli-go.svg?branch=master&style=flat-square)](https://scrutinizer-ci.com/g/danielpieper/mrcli-go/?branch=master)
[![Scrutinizer Code Quality](https://img.shields.io/scrutinizer/g/danielpieper/mrcli-go.svg?branch=master&style=flat-square)](https://scrutinizer-ci.com/g/danielpieper/mrcli-go/?branch=master)
[![License](https://img.shields.io/badge/license-MIT-red.svg?style=flat-square)](LICENSE)


mrcli-go checks your gitlab instance for pending merge requests.

Get an overview about the total number of pending mr's by approvers and projects,
list mr's for your own, your colleagues or multiple projects.


## Installation

TODO

## Configuration

`mrcli-go` is configured using environment variables.
The `GITLAB_TOKEN` is required. Create a token with the `api` scope: https://gitlab.com/profile/personal_access_tokens
```dotenv
GITLAB_URL=https://gitlab.com # optional, set for on-premise installations
GITLAB_TOKEN=<gitlab token> # create gitlab token with api access: https://gitlab.com/profile/personal_access_tokens
```

## Basic Usage

```bash
mrcli-go overview
mrcli-go project <project names separated by space>
mrcli-go approver <approver names separated by space, leave empty for your mr's>
```

### Author

Daniel Pieper - <github@daniel-pieper.com>

### License

MrCli is licensed under the MIT License - see the `LICENSE` file for details

