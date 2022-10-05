# Github Best Practices

This is intended to be a living document as we test and iterate the best ways of working together. These practices were discussed at a cowork on 9/20/22 after a sprint retro on 9/16/22 that identified the lack of shared Github practices as a repeated issue.

## Git Branching Process

1. Write up a ticket to associate with a git branch
    - Create an issue
    - Create a task
    - Tag the task with the appropriate workstream
    - Can utilize bugfix tags or issue templates if relevant
    - Add the task to the story the task is working towards on the [board](https://github.com/IMLS/estimating-wifi/projects/3))
2. Create a branch with the ticket number in the title (123-example)
    - This can be done in the `Development` section after you've created the issue
    - Developers' names are unnecessary (e.g., smith-feature)
3. Use discretion in choosing a branch name (aim for 1-5 words)

## Git PRs

- Aim to keep PRs small
- The reviewer will approve and merge the PR, in the interest of time
- Expected timeline for reviewing code is around 2 days
- `main` currently serves as our development branch

## Git Comments
- Aim to keep comments concise and informative

## Deleting Branches
- Deleting a branch immediately after merging
- In general, developers should take ownership of deleting their branches when their work is completed