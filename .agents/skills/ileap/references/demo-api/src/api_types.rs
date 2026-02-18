/*
 * Copyright (c) Martin Pompéry
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the crate's root directory of this source tree.
 */
//! Use Case 001 REST API-related type definitions

#![allow(renamed_and_removed_lints)]
#![allow(clippy::blocks_in_conditions)]

use chrono::{DateTime, Utc};
use ileap_data_model::{ILeapType, Tad};
use okapi::openapi3::Responses;
use pact_data_model::{PfId, ProductFootprint, UuidError};
use rocket::serde::json::Json;
use rocket::{
    http::Header,
    serde::{Deserialize, Serialize},
};
use rocket_okapi::gen::OpenApiGenerator;
use rocket_okapi::response::OpenApiResponderInner;
use rocket_okapi::OpenApiError;
use schemars::JsonSchema;
use uuid::Uuid;

#[derive(Debug, Responder)]
pub(crate) enum PfListingResponse {
    Finished(Json<PfListingResponseInner>),
    Cont(Json<PfListingResponseInner>, Header<'static>),
}

#[derive(Debug, Serialize, Deserialize, Clone, JsonSchema, PartialEq)]
#[serde(crate = "rocket::serde", rename_all = "camelCase")]
/// HTTP Body of Action `GetFootprint`
pub(crate) struct ProductFootprintResponse {
    pub(crate) data: ProductFootprint<ILeapType>,
}

#[derive(Debug, Serialize, Deserialize, JsonSchema, PartialEq)]
#[serde(crate = "rocket::serde", rename_all = "camelCase")]
/// HTTP Body of Action `ListFootprints`
pub(crate) struct PfListingResponseInner {
    pub(crate) data: Vec<ProductFootprint<ILeapType>>,
}

#[derive(Debug, Responder)]
pub(crate) enum TadListingResponse {
    Finished(Json<TadListingResponseInner>),
    Cont(Json<TadListingResponseInner>, Header<'static>),
}

#[derive(Debug, Serialize, Deserialize, Clone, JsonSchema, PartialEq)]
#[serde(crate = "rocket::serde", rename_all = "camelCase")]
/// HTTP Body of Action `TransportActivityData`
pub(crate) struct TadListingResponseInner {
    pub(crate) data: Vec<Tad>,
}

#[derive(Responder, JsonSchema, Debug)]
pub(crate) enum EventsApiResponse {
    #[response(status = 200)]
    Ok(()),
    #[response(status = 400, content_type = "application/json")]
    NoAuth(crate::error::BadRequest),
    #[response(status = 501, content_type = "application/json")]
    _NotImpl(crate::error::NotImplemented),
    #[response(status = 400, content_type = "application/json")]
    BadReq(crate::error::BadRequest),
}

#[derive(Serialize, Deserialize, Debug, JsonSchema, PartialEq)]
#[serde(crate = "rocket::serde")]
pub(crate) struct PathfinderEvent {
    pub(crate) specversion: String,
    pub(crate) id: String,
    pub(crate) source: String,
    pub(crate) time: DateTime<Utc>,
    #[serde(flatten)]
    pub(crate) data: PathfinderEventData,
}

#[derive(Serialize, Deserialize, Debug, JsonSchema, PartialEq)]
#[serde(crate = "rocket::serde")]
#[serde(tag = "type")]
/// possible contents of `data` property of events - see tech specs section 6 for details
pub(crate) enum PathfinderEventData {
    #[serde(rename = "org.wbcsd.pathfinder.ProductFootprint.Published.v1")]
    /// the contents of the `data` field of a `PF Update Event` – see Tech Specs section 6.8.3
    PFUpdateEvent(PFEventData<PFUpdateEventBody>),

    #[serde(rename = "org.wbcsd.pathfinder.ProductFootprintRequest.Created.v1")]
    /// the contents of the `data` field of a `PF Request Event` – see Tech Specs section 6.8.4.1
    PFRequestEvent(PFEventData<PFRequestEventBody>),
    //todo: add event types PF Response Event and PF Response Error Event
}

impl<T> From<T> for PFEventData<T> {
    fn from(data: T) -> Self {
        Self { data }
    }
}

#[derive(Serialize, Deserialize, Debug, JsonSchema, PartialEq)]
#[serde(crate = "rocket::serde")]
pub(crate) struct PFEventData<T> {
    pub(crate) data: T,
}

#[derive(Serialize, Deserialize, Debug, JsonSchema, PartialEq)]
#[serde(crate = "rocket::serde", rename_all = "camelCase")]
pub(crate) struct PFUpdateEventBody {
    pub(crate) pf_ids: Vec<PfId>,
}

#[derive(Serialize, Deserialize, Debug, JsonSchema, PartialEq)]
#[serde(crate = "rocket::serde", rename_all = "camelCase")]
pub(crate) struct PFRequestEventBody {
    pub(crate) pf: rocket::serde::json::Value,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub(crate) comment: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, JsonSchema, PartialEq)]
#[serde(crate = "rocket::serde", rename_all = "camelCase")]
pub(crate) struct PfIdParam(pub PfId);

impl<'a> rocket::request::FromParam<'a> for PfIdParam {
    type Error = UuidError;

    fn from_param(param: &'a str) -> Result<Self, Self::Error> {
        let uuid = Uuid::parse_str(param).map_err(UuidError::ParseError)?;
        if uuid.get_version_num() != 4 {
            Err(UuidError::VersionError)
        } else {
            Ok(PfIdParam(PfId(uuid)))
        }
    }
}

fn openapi_link_header() -> okapi::openapi3::Header {
    okapi::openapi3::Header {
        description: None,
        value: okapi::openapi3::ParameterValue::Schema {
            style: None,
            explode: None,
            allow_reserved: false,
            example: Some(
                "https://api.example.com/2/footprints?[...]"
                    .to_owned()
                    .into(),
            ),
            examples: None,
            schema: okapi::openapi3::SchemaObject {
                instance_type: Some(schemars::schema::InstanceType::String.into()),
                ..Default::default()
            },
        },
        required: false,
        deprecated: false,
        allow_empty_value: false,
        extensions: Default::default(),
    }
}

impl OpenApiResponderInner for PfListingResponse {
    fn responses(
        gen: &mut rocket_okapi::gen::OpenApiGenerator,
    ) -> rocket_okapi::Result<okapi::openapi3::Responses> {
        use okapi::openapi3::RefOr;

        let mut responses: okapi::openapi3::Responses =
            <Json<PfListingResponseInner>>::responses(gen)?;

        match &mut responses.responses["200"] {
            RefOr::Object(response) => {
                let header = okapi::openapi3::Header {
                    description: Some(
                        "Link header to next result set. See Tech Specs section 6.6.1".to_owned(),
                    ),
                    ..openapi_link_header()
                };
                let header = RefOr::Object(header);
                response.headers.insert("link".to_owned(), header);
            }
            _ => {
                panic!("expected object");
            }
        }

        Ok(responses)
    }
}

impl OpenApiResponderInner for TadListingResponse {
    fn responses(
        gen: &mut rocket_okapi::gen::OpenApiGenerator,
    ) -> rocket_okapi::Result<okapi::openapi3::Responses> {
        use okapi::openapi3::RefOr;

        let mut responses: okapi::openapi3::Responses =
            <Json<TadListingResponseInner>>::responses(gen)?;

        match &mut responses.responses["200"] {
            RefOr::Object(response) => {
                let header = okapi::openapi3::Header {
                    description: Some("Link header to next result set.".to_owned()),
                    ..openapi_link_header()
                };
                let header = RefOr::Object(header);
                response.headers.insert("link".to_owned(), header);
            }
            _ => {
                panic!("expected object");
            }
        }

        Ok(responses)
    }
}

impl OpenApiResponderInner for EventsApiResponse {
    fn responses(gen: &mut OpenApiGenerator) -> Result<Responses, OpenApiError> {
        Ok(Responses {
            responses: okapi::map! {
                "200".to_owned() => <()>::responses(gen)?.responses["200"].clone(),
                "400".to_owned() => crate::error::BadRequest::responses(gen)?.responses["400"].clone(),
                "403".to_owned() => crate::error::AccessDenied::responses(gen)?.responses["403"].clone(),
                "501".to_owned() => crate::error::NotImplemented::responses(gen)?.responses["501"].clone(),
            },
            ..Default::default()
        })
    }
}

#[test]
fn test_pathfinder_event_deser() {
    use chrono::prelude::*;
    use uuid::uuid;

    let time = Utc.with_ymd_and_hms(2022, 5, 31, 17, 31, 00).unwrap();
    let event = PathfinderEvent {
        specversion: "1.0".to_owned(),
        id: "123".to_owned(),
        source: "https://example.com".to_owned(),
        time,
        data: PathfinderEventData::PFUpdateEvent(
            PFUpdateEventBody {
                pf_ids: vec![
                    PfId(uuid!("52B87062-1506-455C-B521-5212212959A8")),
                    PfId(uuid!("8C5D709E-F3A0-4B90-889D-91BF2A68FA19")),
                ],
            }
            .into(),
        ),
    };

    // test deserialization with a reference string
    {
        let event2: PathfinderEvent = rocket::serde::json::from_str(
            r#"{
            "type": "org.wbcsd.pathfinder.ProductFootprint.Published.v1",
            "specversion": "1.0",
            "id": "123",
            "source": "https://example.com",
            "time": "2022-05-31T17:31:00Z",
            "data": {
              "pfIds": [
                "52B87062-1506-455C-B521-5212212959A8",
                "8C5D709E-F3A0-4B90-889D-91BF2A68FA19"
              ]
            }
        }"#,
        )
        .unwrap();

        assert_eq!(event, event2);
    }

    // test serialize->deserialize roundtrip equality of input/output
    {
        let json = rocket::serde::json::to_string(&event).unwrap();
        let event2: PathfinderEvent = rocket::serde::json::from_str(&json).unwrap();
        assert_eq!(event, event2);
    }
}
