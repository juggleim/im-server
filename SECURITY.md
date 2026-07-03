# Security Policy

We take the security of JuggleIM seriously. Thank you for helping keep the project and its users safe.

## Supported versions

Security fixes are provided for the **latest released version** and the current default branch.
Please make sure you can reproduce the issue on an up-to-date checkout before reporting.

## Reporting a vulnerability

**Please do not report security vulnerabilities through public GitHub issues, discussions, or pull
requests.**

Instead, report them privately via GitHub's private vulnerability reporting:

👉 **[Report a vulnerability](https://github.com/juggleim/im-server/security/advisories/new)**

If you are unable to use that channel, reach out through the official website
(https://www.juggle.im) and ask to be put in contact with the security team.

### What to include

To help us triage quickly, please include as much of the following as you can:

- A description of the vulnerability and its potential impact.
- Steps to reproduce, or a proof-of-concept.
- Affected version / commit and deployment mode (Docker, source, binary).
- Any suggested remediation, if you have one.

### What to expect

- We aim to acknowledge new reports within a few business days.
- We will keep you informed of our progress toward a fix.
- Please give us a reasonable amount of time to release a fix before any public disclosure. We're
  happy to credit you in the release notes once the issue is resolved (unless you prefer to remain
  anonymous).

## Scope

Especially relevant to JuggleIM:

- The `app_secret` is **server-side only** and must never be shipped to clients — reports of leakage
  paths are in scope.
- Cross-tenant data access (data isolation is keyed by `app_key`).
- Authentication / authorization bypass on the API or admin gateways.
- Remote crashes or resource-exhaustion in message dispatch/storage.

Thank you for contributing to the security of JuggleIM. 🙏
