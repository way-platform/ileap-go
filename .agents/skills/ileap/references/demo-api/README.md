# Demo API (PACT + iLEAP)

SINE Foundation's demo implementation of the HTTP REST API specified in the [PACT Data Exchange Protocol](https://wbcsd.github.io/data-exchange-protocol/v2/) and the [iLEAP Technical Specifications](https://sine-fdn.github.io/ileap-extension/).

> [!WARNING]
> The `demo-api` is currently WIP. Significant refactoring will take place in the near future.

## How to use the demo API

The demo API is available at https://api.ileap.sine.dev

### Credentials

Four user credentials are available:

| client_id                     | client_secret | Data                                                           |
| ----------------------------- | ------------- | -------------------------------------------------------------- |
| `hello`                       | `pathfinder`  | Global access to all sample data                               |
| `transport_service_user`      | `ileap`       | User is a (demo) shipper. Access to ShipmentFootprints and TAD |
| `transport_service_organizer` | `ileap`       | User is (demo) LSP. Access to TOCs, HOCs, and TAD              |
| `random_data`                 | `random_data` | Access to randomly generated iLEAP data                        |

### Endpoints

The following endpoints are available:

- PACT endpoints (with iLEAP demo data when applicable)
  - `/auth/token` implementing `Authenticate` action
  - `/.well-known/openid-configuration`: OpenId provider configuration document
  - `/2/jwks`: the JSON Web Key Set used to encode and sign the authentication token
  - `/2/footprints` implementing the `ListFootprints` action
  - `/2/footprints/<footprint-id>` implementing the `GetFootprint` action
  - `/2/events` implementing the `Events` action
- iLEAP endpoint
  - `/2/ileap/tad`
- Additional endpoints are:
  - `/openapi.json`: OpenAPI description file which is automatically generated from the types defined in [`api_types.rs`](src/api_types.rs) and endpoints defined in [`main.rs`](src/main.rs)
  - Swagger UI: `/swagger-ui/` if you fancy a visualization

No further endpoints are supported by this implementation and all return `{"message":"Bad Request","code":"BadRequest"`.

## Build instructions

### Build requirements

You need a working and up-to-date Rust toolchain installed. See [https://rustup.rs/](https://rustup.rs/) for details.

After having installed `rustup`, ensure you have the `stable toolchain` installed like this:

```sh
rustup update
rustup toolchain install stable
```

### Building

Once Rust is installed via rustup, just perform

```sh
cargo build
```

### Running locally

You first need to create a private key:

```sh
scripts/keygen.sh
```

Which will create the file `keypair.pem` for you.

Then, you can run the server like this:

```sh
PRIV_KEY=`cat keypair.pem` cargo run
```

To run it at a different port, e.g. 3333:

```sh
ROCKET_PORT=3333 PRIV_KEY=`cat keypair.pem` cargo run
```

### Running the server in a "Production" mode

```sh
## building
cargo build --release

## running
export ROCKET_SECRET_KEY=$(openssl rand -base64 32)
export PRIV_KEY="..."
cargo run --release
```

The resulting binary can be found in `target/release/demo-api`
