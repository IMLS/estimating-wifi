# Development

## Secret leak prevention with pre-commit

Secrets such as API keys, tokens, credentials, etc. should not be
stored with this project's source code.  To help reduce the risk that
secrets are accidentally leaked, a configuration file for the
open-source tool [pre-commit](https://pre-commit.com/) is included in
the repository.  The hooks in this configuration file run scans that
will block developers from committing secrets.

Note: using pre-commit is an opt-in choice developers make in their
own local repositories in their own development environments.  If a
developer installs the pre-commit hooks, the hooks will run when
the developer attempts to run `git commit`; if they don't install
pre-commit, however, the process to create commits remains the same.

Once pre-commit is installed, the workflow remains the same from
the developer's perspective: edit files, stage them with `git add`,
create commits with `git commit`, and push commits with `git push`.
The difference is that when `git commit` is run, prior to the commit
being created, the pre-commit hook will run; this scans the files
in the commit for secrets and, if any are found, will prevent the
commit from being created.  It is possible to override the commit
workflow and force the creation of a commit without running pre-commit.

### Installing the pre-commit software

The pre-commit software is written in Python and may be obtained
from [PyPI](https://pypi.org/) by running:

```shell
pip install pre-commit
```

This should only have to be performed once.

More information is available on the [pre-commit website](https://pre-commit.com/#install)

### Installing the pre-commit hooks

Once the software is installed, the next step is to install
pre-commit's git hooks:

```shell
pre-commit install
```

This sould only have to be performed once.

More information is available on the [pre-commit website](https://pre-commit.com/#3-install-the-git-hook-scripts)

### Usage

The repository includes a pre-commit configuration file that's ready to
go.  The only steps the developer needs to take are to install the
pre-commit software and the pre-commit git hook.  After those two steps
are performed once, the developer does not need to do anything
differently than they normally do -- no workflow changes are required.
