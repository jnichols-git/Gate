# Contributing to Gate

There are three main things you can do to help with the development of this project:

- Test the project
- Use the issue tracker
- Change the codebase

## Testing

Gate is currently lacking testing in real-world situations; if you're interested in using the project
in the future, testing the server and how it functions in a non-production environment would be very helpful.

## Issue Tracker

Gate uses GitHub Issues for tracking user feedback, feature requests, and bugs. If you find any of those things,
please put in an issue and the developers will get to it as soon as possible! You may also be assigned to issues
on request.

## Changing the Codebase

If you'd like to directly contribute to the code for the project, please make a fork, a branch for your changes,
and then enter a pull request. We have a set of guidelines regarding code quality, including a few key words.
- SHOULD: Encouraged, but not mandatory. Violations of SHOULD rules should be justified in PRs.
- MUST: Mandatory. Violations of MUST rules will prevent a PR from being approved.
- EXCEPT: Invalidates rules for specific enumerated conditions.

Contribution Guidelines:

1. Comments and Documentation
    1. You SHOULD include comments explaining functionality if it is not self-evident from the code.
    2. You MUST include comments explaining the parameters, callers, and return values of functions, EXCEPT in test files.
    3. You SHOULD follow up-to-date Golang [commenting conventions](https://go.dev/doc/comment).
2. Variable and Type Naming
    1. You SHOULD not use single-character variable or type names.
    2. Your variable and type names MUST be camelCase if not exported, or PascalCase if exported.
    3. You SHOULD use variable names that explain the purpose of the data they hold.
    4. Rules 1-3 apply to all variables EXCEPT for temporary `int` variables used in iteration, which should use `i`, `j`, and `k`.
    If you need more than 3 nested loops, you SHOULD re-evaluate the algorithm you are using.
3. Exported and Non-Exported Functions and Data Types
    1. You SHOULD not export functions or data types that do not need to be used outside of their modules.
    2. You MUST not export data types that contain private information.
    3. If a data type is not exported, exported functions MUST not return that data type.
4. Error Handling
    1. Your functions MUST handle errors from functions they call by returning that error.
    2. Your functions MUST handle invalid state by returning an error. The error SHOULD describe the invalid state, EXCEPT if that state contains private data.
    3. You MUST handle all errors that occur in an exported function.
    4. You SHOULD handle all errors that occur in a non-exported function.
5. Testing
    1. Your code MUST include tests.
    2. Your test coverage SHOULD be at least 90%, and MUST be at least 80%, EXCEPT for server-related code.
6. Style and Formatting
    1. Your code MUST be run through `gofmt`.
    2. Your code MUST have valid output from `godoc`.

## Templates

### Function Comment
```
// Description.
//
// Input:
//   - 
// Output:
//   - 
func exampleFunc() {...}
```
