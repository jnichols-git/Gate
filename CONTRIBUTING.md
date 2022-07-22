# Contributing to Auth

There are three main things you can do to help with the development of this project:

- Test the project
- Use the issue tracker
- Change the codebase

## Testing

auth is currently lacking testing in real-world situations; if you're interested in using the project
in the future, testing the server and how it functions in a non-production environment would be very helpful.

## Issue Tracker

auth uses GitHub Issues for tracking user feedback, feature requests, and bugs. If you find any of those things,
please put in an issue and the developers will get to it as soon as possible! You may also be assigned to issues
on request.

## Changing the Codebase

If you'd like to directly contribute to the code for the project, please make a fork, a branch for your changes,
and then enter a pull request. We have a set of guidelines regarding code quality, including a few key words.
If you SHOULD do something, it is encouraged, but not mandatory. If you MUST do something, it is a requirement for new code contributions.

Contribution Guidelines:

- Your code SHOULD be well commented to explain functionality where it is not self-explanatory, and MUST include
comments fully explaining the parameters and return values of functions.
- You SHOULD not use single-character variable names, with the exception of iteration. Your variable names SHOULD
explain the purpose of the data they hold.
- You MUST not export functions or data structures that do not need to be used outside of their modules. If your changes
export a function that was previously non-exported, you MUST justify the change in your PR.
- Your functions MUST return an error if one can result from a statement contained within. You MUST handle all
errors in exported functions. You SHOULD handle all errors in non-exported functions.
- Your code MUST include tests. Your test coverage SHOULD be at least 90%, and MUST be at least 80%. These guidelines
are excepted for server-related code at this time.
