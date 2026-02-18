use rocket::serde::{Deserialize, Serialize};
use url::Url;

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq)]
#[serde(crate = "rocket::serde")]
pub(crate) struct OpenIdConfiguration {
    pub(crate) issuer: Url,
    pub(crate) authorization_endpoint: String,
    pub(crate) token_endpoint: String,
    pub(crate) jwks_uri: String,
    pub(crate) response_types_supported: Vec<String>,
    pub(crate) subject_types_supported: Vec<String>,
    pub(crate) id_token_signing_alg_values_supported: Vec<String>,
}
