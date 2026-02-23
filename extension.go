package ileap

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/way-platform/ileap-go/ileapv1pb"
)

// Extension schema constants for iLEAP DataModelExtension.
const (
	DataSchemaShipmentFootprint = "https://api.ileap.sine.dev/shipment-footprint.json"
	DataSchemaTOC              = "https://api.ileap.sine.dev/toc.json"
	DataSchemaHOC              = "https://api.ileap.sine.dev/hoc.json"
	ExtensionSpecVersion       = "2.0.0"
	ExtensionDocumentation    = "https://sine-fdn.github.io/ileap-extension/"
)

// NewShipmentFootprintExtension converts a typed ShipmentFootprint to a DataModelExtension.
func NewShipmentFootprintExtension(sf *ileapv1pb.ShipmentFootprint) (*ileapv1pb.DataModelExtension, error) {
	return newExtension(sf, DataSchemaShipmentFootprint)
}

// NewTOCExtension converts a typed TOC to a DataModelExtension.
func NewTOCExtension(toc *ileapv1pb.TOC) (*ileapv1pb.DataModelExtension, error) {
	return newExtension(toc, DataSchemaTOC)
}

// NewHOCExtension converts a typed HOC to a DataModelExtension.
func NewHOCExtension(hoc *ileapv1pb.HOC) (*ileapv1pb.DataModelExtension, error) {
	return newExtension(hoc, DataSchemaHOC)
}

func newExtension(m proto.Message, dataSchema string) (*ileapv1pb.DataModelExtension, error) {
	data, err := protojson.Marshal(m)
	if err != nil {
		return nil, err
	}
	s := &structpb.Struct{}
	if err := protojson.Unmarshal(data, s); err != nil {
		return nil, err
	}
	ext := &ileapv1pb.DataModelExtension{}
	ext.SetSpecVersion(ExtensionSpecVersion)
	ext.SetDataSchema(dataSchema)
	ext.SetDocumentation(ExtensionDocumentation)
	ext.SetData(s)
	return ext, nil
}
