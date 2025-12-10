# Contribution Guidlines

First of all, thanks for taking the time and effort to read these guidlines and support the project, we appriciate all the help we can get and value any work that is put in.


## Set up

Create a new issue (from the [contribution issue template](https://github.com/CTNOriginals/BitburnerGoFilesync/blob/stable/.github/ISSUE_TEMPLATE/contribution.md)) explaining your idea and how you are going to go about it.
You are expected to update this issue throughout your development to keep everyone informed of the state of your contribution,
think for example of an update where you had to pivot from one strategy to another as you realised it might be a better approach.
You may also create sub-issues under that initial contribution issue that further explain you process.

You will first have to create a fork of this repository,
make sure you fork the [development branch](https://github.com/CTNOriginals/BitburnerGoFilesync/tree/development) this is the branch that is the most up to date out of all and the only one allowed to merge back into.


## Commits

Commits should be done for each group of changes you make, 
all the changes in the commit have to be related to eachother in a clear way and should not require a description as to why it relates to one another.

More over less commits, so make sure to commit often.
Keep the amount of changed files in one single commit as low as possible.

Always add a clear title to the commit that describes the change nicely and is not too genaric so it can be found in a huge list of other commits that may or may not be genaric.
The title of the commit should be no longer then 40 characters at most (with some edge cases of course)

Describe the changes you make to the project clearly and in a consise way in the description of each commit you push.
Try to include issue references to the description that point to related issues once it also comes up in the description, for example: "The bug that was not a feature #123 has been patched".


# Style/Structure Guide

<!-- This project uses [gofumpt](https://github.com/mvdan/gofumpt?tab=readme-ov-file#installation) as its formatter. -->

<!-- Project spesific rules: -->
- Initial variable declerations use 'var' instead of walrus ':='
    * When the 'var' keyword is used, it is way easier to spot the decleration as it stand out from the surrounding text as opposed to a walrus ':=' operator.<br>
    If the 'var' keyword is formatted away on save you may need to add the following to you config: ```"gopls": {"formatting.gofumpt": false}```

If you are unsure on how to style something, try to look around for some other code that may relate to what you are doing and try to copy the style it applies.


# Pull Request (PR)

Once you are done making changes to your fork of the project and made everything nice and tidy again, you are ready to finally submit your PR and receive the praise you have been looking forward to.

When you create the PR, make sure to merge it back into the [development branch](https://github.com/CTNOriginals/BitburnerGoFilesync/tree/development).

After you submit the pull request, someone will have to review it, in this process it is very likely they will have further questions about the changes you made, so keep an eye on your PR.