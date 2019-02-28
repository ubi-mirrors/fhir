// Copyright 2019 Google LLC
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

#include "google/fhir/stu3/annotations.h"

#include <string>

#include "google/protobuf/descriptor.pb.h"
#include "google/protobuf/message.h"
#include "proto/stu3/annotations.pb.h"

namespace google {
namespace fhir {
namespace stu3 {

const string& GetFhirProfileBase(const google::protobuf::Descriptor* descriptor) {
  return descriptor->options().GetExtension(stu3::proto::fhir_profile_base);
}

const string& GetStructureDefinitionUrl(const google::protobuf::Descriptor* descriptor) {
  return descriptor->options().GetExtension(
      stu3::proto::fhir_structure_definition_url);
}

const bool IsChoiceType(const google::protobuf::FieldDescriptor* field) {
  return field->type() == google::protobuf::FieldDescriptor::Type::TYPE_MESSAGE &&
         field->message_type()->options().GetExtension(
             stu3::proto::is_choice_type);
}

const bool IsPrimitive(const google::protobuf::Descriptor* descriptor) {
  return descriptor->options().GetExtension(
             stu3::proto::structure_definition_kind) ==
         stu3::proto::StructureDefinitionKindValue::KIND_PRIMITIVE_TYPE;
}

const bool IsResource(const google::protobuf::Descriptor* descriptor) {
  return descriptor->options().GetExtension(
             stu3::proto::structure_definition_kind) ==
         stu3::proto::StructureDefinitionKindValue::KIND_RESOURCE;
}

const bool IsReference(const google::protobuf::Descriptor* descriptor) {
  return descriptor->options().ExtensionSize(stu3::proto::fhir_reference_type) >
         0;
}

const bool HasValueset(const google::protobuf::Descriptor* descriptor) {
  return descriptor->options().HasExtension(stu3::proto::fhir_valueset_url);
}

const string& GetInlinedCodingSystem(const ::google::protobuf::FieldDescriptor* field) {
  return field->options().GetExtension(stu3::proto::fhir_inlined_coding_system);
}

const string& GetInlinedCodingCode(const ::google::protobuf::FieldDescriptor* field) {
  return field->options().GetExtension(stu3::proto::fhir_inlined_coding_code);
}

}  // namespace stu3
}  // namespace fhir
}  // namespace google