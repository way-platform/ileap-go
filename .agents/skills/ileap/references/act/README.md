# ACT (beta) for iLEAP v1.0.0

ACT (Automated Conformance Testing) tests your implementationf [iLEAP Technical
Specifications v1.0.0](https://sine-fdn.github.io/ileap-extension/).

For this, ACT performs the tests included in [Appendix C: iLEAP Conformance
Testing](https://sine-fdn.github.io/ileap-extension/#appendix-c).

By integrating ACT into your CI/CD pipeline and development processes, you can streamline the iLEAP
Tech Specs implementation and conformance process in general.

A web interface is also available at [act.sine.dev](https://act.sine.dev).


> [!NOTE]
> ACT is now a stable beta.
>
> If you encounter any issues or difficulties, give us [feedback](#contact), open an issue on the [ACT GitHub repository](https://github.com/sine-fdn/act/issues).


## Usage

You can use ACT in the CLI or in a GitHub workflow.

All you need to do is to adapt below commands to your environment and run them in your terminal or in your GitHub workflow:

```sh
curl -sSf https://raw.githubusercontent.com/sine-fdn/act/main/act.sh |\
  bash -s -- test -b "<url>" -u "<user>" -p "<password>"
```

(i.e. `<user>`, `<password>` etc. are placeholders for your actual client credentials)

Example usage with SINE's iLEAP demo API
```sh
curl -sSf https://raw.githubusercontent.com/sine-fdn/act/main/act.sh |\
  bash -s -- test -b "https://api.ileap.sine.dev" -u "hello" -p "pathfinder"
```

### Options

```sh
Options:
  -b, --basepath <BASEPATH>            URL of the API
      --auth-basepath <AUTH_BASEPATH>  Auth Base URL (if different from the API endpoint)
  -u, --user <USER>                    Basic auth username
  -p, --password <PASSWORD>            Basic auth password
      --expired-token <EXPIRED_TOKEN>  Expired token (some tests are skipped if not provided)
  -j, --json [<FILE>]                  Export results to JSON (optionally specify output file)
  -h, --help                           Print help
  -V, --version                        Print version
```

## GitHub Workflow

Adding ACT to your CI/CD pipeline with GitHub is as simple as including the following job in your GitHub workflow, replacing `<url>`, `<user>`, and `<password>` with the URL of your API and the Basic Auth credentials (user and password).

```TOML
  act_test:
    name: ACT Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run ACT (test)
        run: |
          set -o pipefail
          curl -sSf https://raw.githubusercontent.com/sine-fdn/act/main/act.sh | \
          bash -s -- test \
          -b `<url>` \
          -u `<user>` \
          -p `<password>` \
```


### Security Considerations

If you intend to use this action to test live or otherwise production-like systems, you should make
     sure that credentials are passed in as secrets (e.g., `${{secrets.ACT_USER}}` and
     `${{secrets.ACT_PASSWORD}}`)


## Limitations

> [!IMPORTANT]
>
> There are builds ready for ARM64 and x86_64 architectures and the `ubuntu-latest` runner only.
> If you need support for a different architecture or runner, please let us know.

### Coverage

PACT Test Cases

ACT uses a self-hosted version of the PACT Conformance Service v1.4.0 to run PACT tests. iLEAP
currently uses PACT version v2.2.

PACT tests are written and maintained by the PACT community. For more information, see the [pact-conformance-service
repo](https://github.com/wbcsd/pact-conformance-service).

[iLEAP Test Cases](https://sine-fdn.github.io/ileap-extension/#ileap-specific-conformance-tests)

- [x] iLEAP Test Case 001: Get ProductFootprint with ShipmentFootprint
- [x] iLEAP Test Case 002: Get ProductFootprint with TOC
- [x] iLEAP Test Case 003: Get ProductFootprint with HOCs
- [x] iLEAP Test Case 004: Get All TransportActivityData
- [x] iLEAP Test Case 005: Get Filtered List of TransportActivityData
- [x] iLEAP Test Case 006: Get Limited List of TransportActivityData
- [x] iLEAP Test Case 007: Attempt TransportActivityData with Invalid Token
- [x] iLEAP Test Case 008: Attempt TransportActivityData with Expired Token


# Contact

For any questions, feedback, or issues, please contact us at [act-feedback@sine.dev](act-feedback@sine.dev).
