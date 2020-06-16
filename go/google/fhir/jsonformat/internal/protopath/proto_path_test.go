// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO: migrate tests to use test proto instead of FHIR.
package protopath

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	pptpb "google/fhir/jsonformat/internal/protopath/protopathtest_go_proto"
	rdpb "google/fhir/proto/stu3/datatypes_go_proto"
	rfpb "google/fhir/proto/stu3/resources_go_proto"
)

func TestSet(t *testing.T) {
	tests := []struct {
		name  string
		path  Path
		value interface{}
		msg   proto.Message
		want  proto.Message
	}{
		{
			"single field",
			NewProtoPath("meta"),
			&rdpb.Meta{Id: &rdpb.String{Value: "id"}},
			&rfpb.Account{},
			&rfpb.Account{Meta: &rdpb.Meta{Id: &rdpb.String{Value: "id"}}},
		},
		{
			"nested field - parent exists",
			NewProtoPath("meta.id"),
			&rdpb.String{Value: "id"},
			&rfpb.Account{Meta: &rdpb.Meta{}},
			&rfpb.Account{Meta: &rdpb.Meta{Id: &rdpb.String{Value: "id"}}},
		},
		{
			"nested field - no parent",
			NewProtoPath("meta.id"),
			&rdpb.String{Value: "id"},
			&rfpb.Account{},
			&rfpb.Account{Meta: &rdpb.Meta{Id: &rdpb.String{Value: "id"}}},
		},
		{
			"nested field - set primitive",
			NewProtoPath("meta.id.value"),
			"id",
			&rfpb.Account{},
			&rfpb.Account{Meta: &rdpb.Meta{Id: &rdpb.String{Value: "id"}}},
		},
		{
			"repeated field - index exists",
			NewProtoPath("meta.security.0.id.value"),
			"code2",
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code1"}}}}},
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code2"}}}}},
		},
		{
			"repeated field - end",
			NewProtoPath("meta.security.-1.id.value"),
			"code",
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{}}},
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code"}}}}},
		},
		{
			"repeated field - no parent",
			NewProtoPath("meta.security.0.id.value"),
			"code",
			&rfpb.Account{},
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code"}}}}},
		},
		{
			"repeated field - clear",
			NewProtoPath("meta.security"),
			Zero,
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code"}}}}},
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{}}},
		},
		{
			"repeated field - set",
			NewProtoPath("meta.security"),
			[]*rdpb.Coding{{Id: &rdpb.String{Value: "code"}}, {Id: &rdpb.String{Value: "text"}}},
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code"}}}}},
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code"}}, {Id: &rdpb.String{Value: "text"}}}}},
		},
		{
			"repeated field element",
			NewProtoPath("meta.security.-1"),
			&rdpb.Coding{Id: &rdpb.String{Value: "code"}},
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{}}},
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code"}}}}},
		},
		{
			"missing field - zero value of pointer",
			NewProtoPath("meta.id"),
			Zero,
			&rfpb.Account{Meta: &rdpb.Meta{Id: &rdpb.String{Value: "id"}}},
			&rfpb.Account{Meta: &rdpb.Meta{}},
		},
		{
			"nil message",
			NewProtoPath("message_field"),
			(*pptpb.Message_InnerMessage)(nil),
			&pptpb.Message{},
			&pptpb.Message{},
		},
		{
			"missing field - zero value of string",
			NewProtoPath("meta.id.value"),
			Zero,
			&rfpb.Account{Meta: &rdpb.Meta{Id: &rdpb.String{Value: "id"}}},
			&rfpb.Account{Meta: &rdpb.Meta{Id: &rdpb.String{}}},
		},
		{
			"oneof",
			NewProtoPath("observation.id.value"),
			"id",
			&rfpb.ContainedResource{
				OneofResource: &rfpb.ContainedResource_Observation{
					Observation: &rfpb.Observation{
						Id:       &rdpb.Id{Value: ""},
						Language: &rdpb.LanguageCode{Value: "xyz"}}}},
			&rfpb.ContainedResource{
				OneofResource: &rfpb.ContainedResource_Observation{
					Observation: &rfpb.Observation{
						Id:       &rdpb.Id{Value: "id"},
						Language: &rdpb.LanguageCode{Value: "xyz"}}}},
		},
		{
			"oneof - empty",
			NewProtoPath("observation.id.value"),
			"id",
			&rfpb.ContainedResource{},
			&rfpb.ContainedResource{
				OneofResource: &rfpb.ContainedResource_Observation{
					Observation: &rfpb.Observation{
						Id: &rdpb.Id{Value: "id"}}}},
		},
		{
			"oneof - wrap resource",
			NewProtoPath("observation"),
			&rfpb.Observation{Id: &rdpb.Id{Value: "id"}},
			&rfpb.ContainedResource{},
			&rfpb.ContainedResource{
				OneofResource: &rfpb.ContainedResource_Observation{
					Observation: &rfpb.Observation{
						Id: &rdpb.Id{Value: "id"}}}},
		},
		{
			"oneof - primitive",
			NewProtoPath("oneof"),
			true,
			&pptpb.Message{},
			&pptpb.Message{Oneof: &pptpb.Message_OneofPrimitiveField{OneofPrimitiveField: true}},
		},
		{
			"oneof - message",
			NewProtoPath("oneof"),
			&pptpb.Message_InnerMessage{
				InnerField: 1,
			},
			&pptpb.Message{},
			&pptpb.Message{
				Oneof: &pptpb.Message_OneofMessageField{
					OneofMessageField: &pptpb.Message_InnerMessage{
						InnerField: 1,
					}}},
		},
		{
			"oneof - nil message",
			NewProtoPath("oneof_message_field"),
			(*pptpb.Message_InnerMessage)(nil),
			&pptpb.Message{},
			&pptpb.Message{},
		},
		{
			"enum",
			NewProtoPath("type"),
			pptpb.MessageType_TYPE_1,
			&pptpb.Message{},
			&pptpb.Message{Type: pptpb.MessageType_TYPE_1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg := proto.Clone(test.msg)
			if err := Set(msg, test.path, test.value); err != nil {
				t.Fatalf("Set failed with err: %v", err)
			}

			if diff := cmp.Diff(test.want, msg, protocmp.Transform()); diff != "" {
				t.Errorf("Set(_, %v, %v) => %v returned unexpected diff: (-want, +got) %v", test.value, test.path, test.msg, diff)
			}
		})
	}
}

func TestSet_Errors(t *testing.T) {
	tests := []struct {
		name  string
		path  Path
		value interface{}
		msg   proto.Message
	}{
		{
			"empty path part",
			NewProtoPath("foo."),
			Zero,
			&rfpb.Account{},
		},
		{
			"empty path",
			Path{},
			Zero,
			&rfpb.Account{},
		},
		{
			"invalid field",
			NewProtoPath("meta2"),
			&rdpb.Meta{Id: &rdpb.String{Value: "id"}},
			&rfpb.Account{},
		},
		{
			"invalid nested field - parent exists",
			NewProtoPath("meta.id2"),
			&rdpb.String{Value: "id"},
			&rfpb.Account{Meta: &rdpb.Meta{}},
		},
		{
			"nested field - incorrect type - primitive",
			NewProtoPath("meta.id.value"),
			1,
			&rfpb.Account{},
		},
		{
			"nested field - incorrect type - message",
			NewProtoPath("meta.id"),
			&rdpb.Uri{Value: "uri"},
			&rfpb.Account{},
		},
		{
			"repeated field - negative index",
			NewProtoPath("meta.security.-2.id.value"),
			"code2",
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code1"}}}}},
		},
		{
			"repeated field - invalid index",
			NewProtoPath("meta.security.foo.id.value"),
			"code2",
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code1"}}}}},
		},
		{
			"repeated field - index too high",
			NewProtoPath("meta.security.2.id.value"),
			"code",
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code1"}}}}},
		},
		{
			"oneof",
			NewProtoPath("foo.id.value"),
			Zero,
			&rfpb.ContainedResource{},
		},
		{
			"oneof - set oneof - missing proto",
			NewProtoPath("oneof"),
			&pptpb.Missing{},
			&pptpb.Message{},
		},
		{
			"oneof - set oneof - missing primitive",
			NewProtoPath("oneof"),
			float32(1.0),
			&pptpb.Message{},
		},
		{
			"oneof - set oneof - duplicate primitive type",
			NewProtoPath("oneof"),
			int32(1),
			&pptpb.Message{},
		},
		{
			"oneof - set oneof - duplicate message type",
			NewProtoPath("oneof"),
			&pptpb.Message_InnerMessage2{},
			&pptpb.Message{},
		},
		{
			"complex value",
			NewProtoPath("meta"),
			struct{}{},
			&rfpb.Account{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg := proto.Clone(test.msg)
			if err := Set(msg, test.path, test.value); err == nil {
				t.Fatalf("Set(%v, %v, %v) got error <nil>, expected error", test.msg, test.path, test.value)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name   string
		path   Path
		defVal interface{}
		msg    proto.Message
		want   interface{}
	}{
		{
			"single field",
			NewProtoPath("meta"),
			(*rdpb.Meta)(nil),
			&rfpb.Account{Meta: &rdpb.Meta{Id: &rdpb.String{Value: "id"}}},
			&rdpb.Meta{Id: &rdpb.String{Value: "id"}},
		},
		{
			"nested field",
			NewProtoPath("meta.id.value"),
			"id1",
			&rfpb.Account{Meta: &rdpb.Meta{Id: &rdpb.String{Value: "id"}}},
			"id",
		},
		{
			"repeated field - positive index",
			NewProtoPath("meta.security.0.id.value"),
			"",
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code1"}}}}},
			"code1",
		},
		{
			"repeated field - end",
			NewProtoPath("meta.security.-1.id.value"),
			"",
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code"}}}}},
			"code",
		},
		{
			"repeated field element",
			NewProtoPath("meta.security.-1"),
			Zero,
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code"}}}}},
			&rdpb.Coding{Id: &rdpb.String{Value: "code"}},
		},
		{
			"missing field",
			NewProtoPath("meta.id.value"),
			"code",
			&rfpb.Account{},
			"code",
		},
		{
			"missing field - nil",
			NewProtoPath("meta.id"),
			(*rdpb.String)(nil),
			&rfpb.Account{},
			(*rdpb.String)(nil),
		},
		{
			"missing field - repeated parent",
			NewProtoPath("meta.security.-1.id.value"),
			"code",
			&rfpb.Account{},
			"code",
		},
		{
			"missing field - untyped nil",
			NewProtoPath("meta"),
			nil,
			&rfpb.Account{},
			(*rdpb.Meta)(nil),
		},
		{
			"missing field - nested untyped nil",
			NewProtoPath("meta.id"),
			nil,
			&rfpb.Account{},
			(*rdpb.String)(nil),
		},
		{
			"missing field - repeated",
			NewProtoPath("meta.security"),
			Zero,
			&rfpb.Account{Meta: &rdpb.Meta{}},
			[]*rdpb.Coding{},
		},
		{
			"missing field - zero",
			NewProtoPath("meta.id"),
			Zero,
			&rfpb.Account{},
			(*rdpb.String)(nil),
		},
		{
			"oneof",
			NewProtoPath("observation.id.value"),
			Zero,
			&rfpb.ContainedResource{
				OneofResource: &rfpb.ContainedResource_Observation{
					Observation: &rfpb.Observation{
						Id: &rdpb.Id{Value: "id"}}}},
			"id",
		},
		{
			"oneof - zero",
			NewProtoPath("observation.id.value"),
			Zero,
			&rfpb.ContainedResource{},
			"",
		},
		{
			"oneof - message",
			NewProtoPath("oneof"),
			Zero,
			&pptpb.Message{
				Oneof: &pptpb.Message_OneofMessageField{
					OneofMessageField: &pptpb.Message_InnerMessage{InnerField: 1},
				},
			},
			&pptpb.Message_InnerMessage{InnerField: 1},
		},
		{
			"oneof - primitive",
			NewProtoPath("oneof"),
			Zero,
			&pptpb.Message{
				Oneof: &pptpb.Message_OneofPrimitiveField{
					OneofPrimitiveField: true,
				},
			},
			true,
		},
		{
			"oneof - multiple interfaces - ends on concrete type",
			NewProtoPath("observation.effective.date_time.value_us"),
			Zero,
			&rfpb.ContainedResource{
				OneofResource: &rfpb.ContainedResource_Observation{}},
			int64(0),
		},
		{
			"whole repeated field - not empty",
			NewProtoPath("meta.security"),
			Zero,
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code1"}}}}},
			[]*rdpb.Coding{{Id: &rdpb.String{Value: "code1"}}},
		},
		{
			"whole repeated field - empty",
			NewProtoPath("meta.security"),
			Zero,
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{}}},
			[]*rdpb.Coding{},
		},
		{
			"enum",
			NewProtoPath("type"),
			Zero,
			&pptpb.Message{Type: pptpb.MessageType_TYPE_1},
			pptpb.MessageType_TYPE_1,
		},
		{
			"enum - zero",
			NewProtoPath("type"),
			Zero,
			&pptpb.Message{},
			pptpb.MessageType_INVALID_UNINITIALIZED,
		},
		{
			"enum - default",
			NewProtoPath("type"),
			pptpb.MessageType_TYPE_1,
			&pptpb.Message{},
			pptpb.MessageType_TYPE_1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, err := Get(test.msg, test.path, test.defVal)
			if err != nil {
				t.Fatalf("Get(%v, %v, %v) got err %v, expected <nil>", test.msg, test.path, test.defVal, err)
			}

			if diff := cmp.Diff(val, test.want, protocmp.Transform(), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Get(%v, %v, _) => %v, want %v, diff:\n%s", test.msg, test.path, val, test.want, diff)
			}
		})
	}
}

func TestGet_Errors(t *testing.T) {
	tests := []struct {
		name   string
		path   Path
		defVal interface{}
		msg    proto.Message
	}{
		{
			"empty path part",
			NewProtoPath("foo."),
			Zero,
			&rfpb.Account{},
		},
		{
			"empty path",
			Path{},
			Zero,
			&rfpb.Account{},
		},
		{
			"invalid field",
			NewProtoPath("meta2"),
			nil,
			&rfpb.Account{},
		},
		{
			"invalid nested field - parent exists",
			NewProtoPath("meta.id2"),
			nil,
			&rfpb.Account{Meta: &rdpb.Meta{}},
		},
		{
			"nested field - incorrect default value type - primitive",
			NewProtoPath("meta.id.value"),
			1,
			&rfpb.Account{},
		},
		{
			"nested field - incorrect default value type - message",
			NewProtoPath("meta.id"),
			&rdpb.Uri{Value: "uri"},
			&rfpb.Account{},
		},
		{
			"nested field - trailing path",
			NewProtoPath("meta.id.value.foo"),
			Zero,
			&rfpb.Account{},
		},
		{
			"nested field - invalid path",
			NewProtoPath("meta.id.foo"),
			Zero,
			&rfpb.Account{},
		},
		{
			"repeated field - negative index",
			NewProtoPath("meta.security.-2.id.value"),
			"code2",
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code1"}}}}},
		},
		{
			"repeated field - invalid index",
			NewProtoPath("meta.security.foo.id.value"),
			"code2",
			&rfpb.Account{Meta: &rdpb.Meta{Security: []*rdpb.Coding{{Id: &rdpb.String{Value: "code1"}}}}},
		},
		{
			"unknown oneof",
			NewProtoPath("foo.id.value"),
			Zero,
			&rfpb.ContainedResource{},
		},
		{
			"oneof - empty",
			NewProtoPath("oneof"),
			Zero,
			&pptpb.Message{},
		},
		{
			"nested oneof - empty",
			NewProtoPath("extension.0.value.choice"),
			Zero,
			&rfpb.Account{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := Get(test.msg, test.path, test.defVal); err == nil {
				t.Fatalf("Get(%v, %v, %v) got error <nil>, expected error", test.msg, test.path, test.defVal)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []string{
		"normal",
		"multiple.level",
		"with.number.0.1.2.-1",
		"",
	}
	for _, tc := range tests {
		p := NewProtoPath(tc)
		if p.String() != tc {
			t.Errorf("%v.String() != %s", p, tc)
		}
	}
}