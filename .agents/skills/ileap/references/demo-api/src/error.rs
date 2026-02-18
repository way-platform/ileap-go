/*
 * Copyright (c) Martin Pompéry
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the crate's root directory of this source tree.
 */
use rocket::http::Status;
use rocket::response::{self, Responder};
use rocket::serde::json::Json;
use rocket::serde::{Deserialize, Serialize};
use rocket::{Request, Response};
use rocket_okapi::gen::OpenApiGenerator;
use rocket_okapi::okapi::openapi3::{MediaType, Responses};
use rocket_okapi::okapi::{self, schemars};
use rocket_okapi::response::OpenApiResponderInner;
use rocket_okapi::{JsonSchema, OpenApiError};

/// Error code enum for NoSuchFootprint responses
#[derive(Serialize, Deserialize, JsonSchema, PartialEq, Debug, Clone, Copy)]
#[serde(crate = "rocket::serde")]
pub(crate) enum NoSuchFootprintCode {
    NoSuchFootprint,
}

/// Error code enum for AccessDenied responses
#[derive(Serialize, Deserialize, JsonSchema, PartialEq, Debug, Clone, Copy)]
#[serde(crate = "rocket::serde")]
pub(crate) enum AccessDeniedCode {
    AccessDenied,
}

/// Error code enum for BadRequest responses
#[derive(Serialize, Deserialize, JsonSchema, PartialEq, Debug, Clone, Copy)]
#[serde(crate = "rocket::serde")]
pub(crate) enum BadRequestCode {
    BadRequest,
}

/// Error code enum for NotImplemented responses
#[derive(Serialize, Deserialize, JsonSchema, PartialEq, Debug, Clone, Copy)]
#[serde(crate = "rocket::serde")]
pub(crate) enum NotImplementedCode {
    NotImplemented,
}

/// Error code enum for Unauthorized responses
#[derive(Serialize, Deserialize, JsonSchema, PartialEq, Debug, Clone, Copy)]
#[serde(crate = "rocket::serde")]
pub(crate) enum UnauthorizedCode {
    Unauthorized,
}

#[derive(Serialize, JsonSchema, PartialEq, Debug)]
#[serde(crate = "rocket::serde")]
pub(crate) enum GetPfError {
    NoSuchFootprint(NoSuchFootprint),
    BadRequest(BadRequest),
}

#[derive(Serialize, Deserialize, JsonSchema, PartialEq, Debug)]
#[serde(crate = "rocket::serde")]
#[allow(dead_code)] // TODO: remove struct if not used
/// Response with an error code of `NoSuchFootprint`. See Chapter "Error Codes" of the Tech Specs for mor details.
pub(crate) struct NoSuchFootprint {
    pub(crate) message: &'static str,
    pub(crate) code: NoSuchFootprintCode,
}

#[derive(Serialize, Deserialize, JsonSchema, PartialEq, Debug)]
#[serde(crate = "rocket::serde")]
/// Response with an error code of `AccessDenied`. See Chapter "Error Codes" of the Tech Specs for mor details.
pub(crate) struct AccessDenied {
    pub(crate) message: &'static str,
    pub(crate) code: AccessDeniedCode,
}

/// RFC 6749 OAuth 2.0 Error Response
#[derive(Serialize, JsonSchema, PartialEq, Debug)]
#[serde(crate = "rocket::serde")]
pub(crate) struct OAuth2ErrorMessage {
    pub(crate) error: &'static str,
    pub(crate) error_description: &'static str,
}

#[derive(Serialize, Deserialize, JsonSchema, PartialEq, Debug)]
#[serde(crate = "rocket::serde")]
/// Response with an error code of `BadRequest`. See Chapter "Error Codes" of the Tech Specs for mor details.
pub(crate) struct BadRequest {
    pub(crate) message: &'static str,
    pub(crate) code: BadRequestCode,
}

#[derive(Serialize, Deserialize, JsonSchema, PartialEq, Debug)]
#[serde(crate = "rocket::serde")]
/// Response with an error code of `NotImplemented`. See Chapter "Error Codes" of the Tech Specs for mor details.
pub(crate) struct NotImplemented {
    pub(crate) message: &'static str,
    pub(crate) code: NotImplementedCode,
}

#[derive(Serialize, Deserialize, JsonSchema, PartialEq, Debug)]
#[serde(crate = "rocket::serde")]
#[allow(dead_code)] // TODO: remove struct if not used
/// Response with an error code of `Unauthorized`, used for iLEAP TransportActivityData
pub(crate) struct Unauthorized {
    pub(crate) message: &'static str,
    pub(crate) code: UnauthorizedCode,
}

impl Default for AccessDenied {
    fn default() -> Self {
        Self {
            message: "Access Denied",
            code: AccessDeniedCode::AccessDenied,
        }
    }
}

impl Default for BadRequest {
    fn default() -> Self {
        Self {
            message: "Bad Request",
            code: BadRequestCode::BadRequest,
        }
    }
}

impl Default for NoSuchFootprint {
    fn default() -> Self {
        NoSuchFootprint {
            message: "The specified footprint does not exist",
            code: NoSuchFootprintCode::NoSuchFootprint,
        }
    }
}

impl Default for NotImplemented {
    fn default() -> Self {
        NotImplemented {
            message: "Not Implemented",
            code: NotImplementedCode::NotImplemented,
        }
    }
}

impl Default for Unauthorized {
    fn default() -> Self {
        Unauthorized {
            message: "Unauthorized",
            code: UnauthorizedCode::Unauthorized,
        }
    }
}

impl<'r, 'o: 'r> Responder<'r, 'o> for GetPfError {
    fn respond_to(self, req: &'r Request<'_>) -> response::Result<'o> {
        match self {
            GetPfError::BadRequest(e) => e.respond_to(req),
            GetPfError::NoSuchFootprint(e) => e.respond_to(req),
        }
    }
}

impl<'r, 'o: 'r> Responder<'r, 'o> for NoSuchFootprint {
    fn respond_to(self, request: &'r Request<'_>) -> response::Result<'o> {
        Response::build()
            .merge(Json(self).respond_to(request)?)
            .status(Status::NotFound)
            .ok()
    }
}

impl<'r, 'o: 'r> Responder<'r, 'o> for AccessDenied {
    fn respond_to(self, request: &'r Request<'_>) -> response::Result<'o> {
        Response::build()
            .merge(Json(self).respond_to(request)?)
            .status(Status::Forbidden)
            .ok()
    }
}

impl<'r, 'o: 'r> Responder<'r, 'o> for BadRequest {
    fn respond_to(self, request: &'r Request<'_>) -> response::Result<'o> {
        Response::build()
            .merge(Json(self).respond_to(request)?)
            .status(Status::BadRequest)
            .ok()
    }
}

impl<'r, 'o: 'r> Responder<'r, 'o> for NotImplemented {
    fn respond_to(self, request: &'r Request<'_>) -> response::Result<'o> {
        Response::build()
            .merge(Json(self).respond_to(request)?)
            .status(Status::BadRequest)
            .ok()
    }
}

impl<'r, 'o: 'r> Responder<'r, 'o> for Unauthorized {
    fn respond_to(self, request: &'r Request<'_>) -> response::Result<'o> {
        Response::build()
            .merge(Json(self).respond_to(request)?)
            .status(Status::Unauthorized)
            .ok()
    }
}

impl<'r, 'o: 'r> Responder<'r, 'o> for OAuth2ErrorMessage {
    fn respond_to(self, request: &'r Request<'_>) -> response::Result<'o> {
        Response::build()
            .merge(Json(self).respond_to(request)?)
            .status(Status::BadRequest)
            .ok()
    }
}

impl OpenApiResponderInner for GetPfError {
    fn responses(gen: &mut OpenApiGenerator) -> Result<Responses, OpenApiError> {
        let mut responses = Responses::default();

        let bad_request_responses = BadRequest::responses(gen)?;
        for (code, resp) in bad_request_responses.responses {
            responses.responses.insert(code, resp);
        }

        let no_such_footprint_responses = NoSuchFootprint::responses(gen)?;
        for (code, resp) in no_such_footprint_responses.responses {
            responses.responses.insert(code, resp);
        }

        Ok(responses)
    }
}

impl OpenApiResponderInner for BadRequest {
    fn responses(gen: &mut OpenApiGenerator) -> Result<Responses, OpenApiError> {
        let resp = openapi_response::<BadRequest>(
            gen,
            "400".to_owned(),
            "\
            # 400 Bad Request\n\
            The request given is wrongly formatted or data was missing. \
            "
            .to_owned(),
        );
        Ok(resp)
    }
}

impl OpenApiResponderInner for NotImplemented {
    fn responses(gen: &mut OpenApiGenerator) -> Result<Responses, OpenApiError> {
        let resp = openapi_response::<NotImplemented>(
            gen,
            "501".to_owned(),
            "\
            # 501 Not Implemented\n\
            The request given is not implemented by the server. \
            "
            .to_owned(),
        );
        Ok(resp)
    }
}

impl OpenApiResponderInner for NoSuchFootprint {
    fn responses(gen: &mut OpenApiGenerator) -> Result<Responses, OpenApiError> {
        let resp = openapi_response::<NoSuchFootprint>(
            gen,
            "404".to_owned(),
            "# 404 Not Found".to_owned(),
        );
        Ok(resp)
    }
}

impl OpenApiResponderInner for Unauthorized {
    fn responses(gen: &mut OpenApiGenerator) -> Result<Responses, OpenApiError> {
        let resp = openapi_response::<Unauthorized>(
            gen,
            "401".to_owned(),
            "\
            # 401 Unauthorized\n\
            The access token is not valid. \
            "
            .to_owned(),
        );
        Ok(resp)
    }
}

impl OpenApiResponderInner for AccessDenied {
    fn responses(gen: &mut OpenApiGenerator) -> Result<Responses, OpenApiError> {
        let resp =
            openapi_response::<AccessDenied>(gen, "403".to_owned(), "# 403 Forbidden".to_owned());
        Ok(resp)
    }
}

fn openapi_response<T: JsonSchema>(
    gen: &mut OpenApiGenerator,
    code: String,
    description: String,
) -> Responses {
    use okapi::openapi3::RefOr;

    let schema = gen.json_schema::<T>();
    let resp = okapi::openapi3::Response {
        description,
        content: okapi::map! {
            "application/json".to_owned() => MediaType {
                schema: Some(schema),
                ..Default::default()
            }
        },
        ..Default::default()
    };

    Responses {
        responses: okapi::map! {
            code => RefOr::Object(resp),
        },
        ..Default::default()
    }
}
