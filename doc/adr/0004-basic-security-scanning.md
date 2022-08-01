# 1. Record architecture decisions

Date: 2022-07-28

## Status

Proposed

## Context

We need to add basic automated security scanning for the PiSpots code.

For Phase 3 we need security scanning which will have automated tooling that scans our code, containers, DB, etc. for the most common security concerns, dependency updates, etc. Based on the reserch we should have two good SAST (static aplication security testing) and two good DAST (dynamic application security testing) tools. 
Code scanning should include : 
    - looking for vulnerabilities in the source code
    - looking for bad/unsafe variabels
    - looking for vulnerable dependancies
    - input validation
    - execution control analysis
    - API scanning
    - key walidation
    (other tools and features can be added later)

## Decision
We are going to use tools that are supported in other GSA applications for Phase 3, those tools are:
    - semgrep.dev (already in MegaLinter)
    - snyk.io
    - codeql
    - dependabot

## Consequences

A license for snyk.io needs to be aquired. MegaLinter needs to be updated with any new SAST tools and checked for conflicts.