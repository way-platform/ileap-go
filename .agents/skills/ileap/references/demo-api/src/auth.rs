/*
 * Copyright (c) Martin Pompéry
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the crate's root directory of this source tree.
 */

#![allow(renamed_and_removed_lints)]
#![allow(clippy::blocks_in_conditions)]

use std::collections::HashSet;

use jsonwebtoken::errors::Result;
use jsonwebtoken::{Algorithm, TokenData};
use jsonwebtoken::{DecodingKey, EncodingKey};
use jsonwebtoken::{Header, Validation};

use rocket::http::Status;
use rocket::outcome::Outcome;
use rocket::request::{self, FromRequest, Request};
use rocket::response::status;
use rocket::serde::{Deserialize, Serialize};

use base64::{engine::general_purpose::STANDARD, Engine as _};
use rocket_okapi::okapi::map;
use rocket_okapi::okapi::openapi3::{
    Object, SecurityRequirement, SecurityScheme, SecuritySchemeData,
};
use rocket_okapi::{
    gen::OpenApiGenerator,
    request::{OpenApiFromRequest, RequestHeaderInput},
};
use rsa::pkcs8::DecodePrivateKey;
use rsa::pkcs8::EncodePublicKey;
use rsa::{pkcs8::LineEnding, RsaPrivateKey, RsaPublicKey};

#[derive(Clone)]
pub struct KeyPair {
    pub pub_key: RsaPublicKey,
    enc_key: EncodingKey,
    dec_key: DecodingKey,
}

#[derive(Clone, Debug, Serialize, Deserialize)]
#[serde(crate = "rocket::serde", rename_all = "camelCase")]
pub struct UserToken {
    pub username: String,
}

#[derive(Debug, PartialEq)]
pub struct OAuth2ClientCredentials {
    pub id: String,
    pub secret: String,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(crate = "rocket::serde")]
pub enum OAuth2TokenType {
    #[serde(rename = "bearer")]
    Bearer,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(crate = "rocket::serde")]
pub struct OAuth2TokenReply {
    pub access_token: String,
    pub token_type: OAuth2TokenType,
    pub scope: Option<String>,
    //pub expires_in: u32,
}

#[derive(rocket::form::FromForm)]
pub struct OAuth2ClientCredentialsBody<'r> {
    #[allow(dead_code)]
    pub grant_type: &'r str,
    pub scope: Option<&'r str>,
}

#[derive(Debug)]
pub enum UserTokenError {
    //Missing,
    Invalid,
    Missing,
}

#[rocket::async_trait]
impl<'r> FromRequest<'r> for OAuth2ClientCredentials {
    type Error = ();

    async fn from_request(req: &'r Request<'_>) -> request::Outcome<Self, Self::Error> {
        if let Some(auth_header) = req.headers().get_one("Authorization") {
            let split = auth_header.split_whitespace().collect::<Vec<_>>();
            if split.len() != 2 {
                return Outcome::Error((Status::BadRequest, ()));
            }
            let (basic, payload) = (split[0], split[1]);
            if basic != "Basic" {
                return Outcome::Error((Status::Unauthorized, ()));
            }
            if let Some(credentials) = decode_basic_auth(payload) {
                return Outcome::Success(credentials);
            }
        }

        Outcome::Error((Status::BadRequest, ()))
    }
}

fn decode_basic_auth(raw_auth_info: &str) -> Option<OAuth2ClientCredentials> {
    if let Ok(auth_info) = STANDARD.decode(raw_auth_info) {
        if let Ok(decoded_str) = String::from_utf8(auth_info) {
            let username_password = decoded_str.split(':').collect::<Vec<_>>();
            if username_password.len() == 2 {
                let credentials = OAuth2ClientCredentials {
                    id: String::from(username_password[0]),
                    secret: String::from(username_password[1]),
                };
                return Some(credentials);
            }
        }
    }

    None
}

const BEARER_TOKEN_START: &str = "Bearer ";

#[rocket::async_trait]
impl<'r> FromRequest<'r> for UserToken {
    type Error = status::Custom<UserTokenError>;

    async fn from_request(
        req: &'r Request<'_>,
    ) -> request::Outcome<Self, status::Custom<UserTokenError>> {
        if let Some(authen_header) = req.headers().get_one("Authorization") {
            let authen_str = authen_header.to_string();
            if authen_str.starts_with(BEARER_TOKEN_START) {
                let token = authen_str[BEARER_TOKEN_START.len()..authen_str.len()].trim();
                let key_pair = req.rocket().state::<KeyPair>().unwrap();
                if let Ok(token_data) = decode_token(token.to_string(), key_pair) {
                    return Outcome::Success(token_data.claims);
                }
            }
            Outcome::Error((
                Status::BadRequest,
                status::Custom(Status::BadRequest, UserTokenError::Invalid),
            ))
        } else {
            Outcome::Error((
                Status::Unauthorized,
                status::Custom(Status::Unauthorized, UserTokenError::Missing),
            ))
        }
    }
}

pub fn load_keys() -> KeyPair {
    #[cfg(test)]
    let priv_key = include_str!("../keypair.pem");
    #[cfg(not(test))]
    let priv_key: &str = &std::env::var("PRIV_KEY").expect("PRIV_KEY must be set");

    let private_key = RsaPrivateKey::from_pkcs8_pem(priv_key)
        .unwrap_or_else(|err| panic!("Could not deserialize private key: {}", err));
    let public_key = RsaPublicKey::from(&private_key);

    let pub_key = public_key
        .to_public_key_pem(LineEnding::default())
        .unwrap_or_else(|err| panic!("Could not serialize public key: {}", err));

    let dec_key = DecodingKey::from_rsa_pem(pub_key.as_bytes()).unwrap();
    let enc_key = EncodingKey::from_rsa_pem(priv_key.as_bytes()).unwrap();

    KeyPair {
        pub_key: public_key,
        enc_key,
        dec_key,
    }
}

fn decode_token(token: String, key_pair: &KeyPair) -> Result<TokenData<UserToken>> {
    let mut v = Validation::new(Algorithm::RS256);
    v.validate_exp = false;
    v.required_spec_claims = HashSet::new();

    jsonwebtoken::decode::<UserToken>(&token, &key_pair.dec_key, &v)
}

pub fn encode_token(u: &UserToken, key_pair: &KeyPair) -> Result<String> {
    let header = Header::new(Algorithm::RS256);

    jsonwebtoken::encode(&header, u, &key_pair.enc_key)
}

impl OpenApiFromRequest<'_> for UserToken {
    fn from_request_input(
        _gen: &mut OpenApiGenerator,
        _name: String,
        _required: bool,
    ) -> rocket_okapi::Result<RequestHeaderInput> {
        // Setup global requirement for Security scheme
        let security_scheme = SecurityScheme {
            description: Some("OAuth2 Client Credentials Grant (RFC6749 4.4)".to_owned()),
            // Setup data requirements.
            // In this case the header `Authorization: mytoken` needs to be set.
            data: SecuritySchemeData::OAuth2 {
                flows: rocket_okapi::okapi::openapi3::OAuthFlows::ClientCredentials {
                    token_url: "/auth/token".into(),
                    refresh_url: None,
                    scopes: map! {
                        "footprint:list".to_owned() => "Ability to list footprints".to_owned(),
                        "footprint:read".to_owned() => "Ability to access individual footprints".to_owned(),
                    },
                    extensions: Default::default(),
                },
            },
            extensions: Object::default(),
        };
        // Add the requirement for this route/endpoint
        // This can change between routes.
        let mut security_req = SecurityRequirement::new();
        // Each security requirement needs to be met before access is allowed.
        security_req.insert("BearerAuth".to_owned(), Vec::new());
        // These vvvvvvv-----^^^^^^^^ values need to match exactly!
        Ok(RequestHeaderInput::Security(
            "BearerAuth".to_owned(),
            security_scheme,
            security_req,
        ))
    }
}

#[test]
fn decode_basic_auth_test() {
    assert_eq!(
        Some(OAuth2ClientCredentials {
            id: "martin".into(),
            secret: "secret".into()
        }),
        decode_basic_auth(&STANDARD.encode(b"martin:secret"))
    );
}
