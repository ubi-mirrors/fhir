/*
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include "google/fhir/r4/extensions.h"

#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "google/fhir/test_helper.h"
#include "google/fhir/testutil/proto_matchers.h"
#include "proto/r4/core/datatypes.pb.h"
#include "proto/r4/core/extensions.pb.h"
#include "proto/r4/core/resources/composition.pb.h"
#include "proto/r4/google_extensions.pb.h"
#include "testdata/r4/profiles/test_extensions.pb.h"
#include "tensorflow/core/lib/core/status_test_util.h"

namespace google {
namespace fhir {
namespace r4 {

namespace {

using core::CapabilityStatementSearchParameterCombination;
using core::Composition;
using core::Extension;
using google::Base64BinarySeparatorStride;
using google::EventLabel;
using google::EventTrigger;
using google::PrimitiveHasNoValue;
using ::google::fhir::testutil::EqualsProto;
using testing::DigitalMediaType;

template <class T>
void ReadR4TestData(const std::string& type, T* message, Extension* extension) {
  *message = ReadR4Proto<T>(absl::StrCat("google/", type, ".message.prototxt"));
  *extension = ReadR4Proto<Extension>(
      absl::StrCat("google/", type, ".extension.prototxt"));
}

template <class T>
void TestExtensionToMessage(const std::string& name) {
  T message;
  Extension extension;
  ReadR4TestData(name, &message, &extension);

  T output;
  FHIR_ASSERT_OK(ExtensionToMessage(extension, &output));
  EXPECT_THAT(output, EqualsProto(message));
}

template <class T>
void TestConvertToExtension(const std::string& name) {
  T message;
  Extension extension;
  ReadR4TestData(name, &message, &extension);

  Extension output;
  FHIR_ASSERT_OK(ConvertToExtension(message, &output));
  EXPECT_THAT(output, EqualsProto(extension));
}

TEST(ExtensionsTest, ParseEventTrigger) {
  TestExtensionToMessage<EventTrigger>("trigger");
}

TEST(ExtensionsTest, PrintEventTrigger) {
  TestConvertToExtension<EventTrigger>("trigger");
}

TEST(ExtensionsTest, ParseEventLabel) {
  TestExtensionToMessage<EventLabel>("label");
}

TEST(ExtensionsTest, PrintEventLabel) {
  TestConvertToExtension<EventLabel>("label");
}

TEST(ExtensionsTest, ParsePrimitiveHasNoValue) {
  TestExtensionToMessage<PrimitiveHasNoValue>("primitive_has_no_value");
}

TEST(ExtensionsTest, PrintPrimitiveHasNoValue) {
  TestConvertToExtension<PrimitiveHasNoValue>("primitive_has_no_value");
}

TEST(ExtensionsTest, ParsePrimitiveHasNoValue_Empty) {
  TestExtensionToMessage<PrimitiveHasNoValue>("empty");
}

TEST(ExtensionsTest, PrintPrimitiveHasNoValue_Empty) {
  TestConvertToExtension<PrimitiveHasNoValue>("empty");
}

TEST(ExtensionsTest, ParseCapabilityStatementSearchParameterCombination) {
  TestExtensionToMessage<CapabilityStatementSearchParameterCombination>(
      "capability");
}

TEST(ExtensionsTest, PrintCapabilityStatementSearchParameterCombination) {
  TestConvertToExtension<CapabilityStatementSearchParameterCombination>(
      "capability");
}

TEST(ExtensionsTest, ParseBoundCodeExtension) {
  TestExtensionToMessage<DigitalMediaType>("digital_media_type");
}

TEST(ExtensionsTest, PrintBoundCodeExtension) {
  TestConvertToExtension<DigitalMediaType>("digital_media_type");
}

TEST(ExtensionsTest, ExtractOnlyMatchingExtensionOneFound) {
  Composition composition = PARSE_VALID_STU3_PROTO(R"pb(
    id { value: "1" }
    status { value: FINAL }
    subject { patient_id { value: "P0" } }
    encounter { encounter_id { value: "2" } }
    author { practitioner_id { value: "3" } }
    title { value: "Note" }
    type {
      coding {
        system { value: "type-system" }
        code { value: "RADIOLOGY" }
      }
    }
    date {
      value_us: 4608099660000000
      timezone: "America/New_York"
      precision: SECOND
    }
    extension {
      url { value: "https://g.co/fhir/StructureDefinition/random" }
      extension {
        url { value: "foo" }
        value { string_value { value: "barr" } }
      }
    }
    extension {
      url {
        value: "https://g.co/fhir/StructureDefinition/base64Binary-separatorStride"
      }
      extension {
        url { value: "separator" }
        value { string_value { value: "!" } }
      }
      extension {
        url { value: "stride" }
        value { positive_int { value: 5 } }
      }
    }
    extension {
      url { value: "https://g.co/fhir/StructureDefinition/more-random" }
      extension {
        url { value: "baz" }
        value { string_value { value: "quux" } }
      }
    }
  )pb");

  StatusOr<Base64BinarySeparatorStride> extracted =
      ExtractOnlyMatchingExtension<Base64BinarySeparatorStride>(composition);

  Base64BinarySeparatorStride expected;
  ASSERT_TRUE(::google::protobuf::TextFormat::ParseFromString(R"pb(
                                                      separator { value: "!" }
                                                      stride { value: 5 }
                                                    )pb",
                                                    &expected));

  FHIR_ASSERT_OK(extracted.status());
  EXPECT_THAT(extracted.ValueOrDie(), EqualsProto(expected));
}

TEST(ExtensionsTest, ExtractOnlyMatchingExtensionNoneFound) {
  Composition composition = PARSE_VALID_STU3_PROTO(R"pb(
    id { value: "1" }
    status { value: FINAL }
    subject { patient_id { value: "P0" } }
    encounter { encounter_id { value: "2" } }
    author { practitioner_id { value: "3" } }
    title { value: "Note" }
    type {
      coding {
        system { value: "type-system" }
        code { value: "RADIOLOGY" }
      }
    }
    date {
      value_us: 4608099660000000
      timezone: "America/New_York"
      precision: SECOND
    }
    extension {
      url { value: "https://g.co/fhir/StructureDefinition/random" }
      extension {
        url { value: "foo" }
        value { string_value { value: "barr" } }
      }
    }
    extension {
      url { value: "https://g.co/fhir/StructureDefinition/more-random" }
      extension {
        url { value: "baz" }
        value { string_value { value: "quux" } }
      }
    }
  )pb");

  StatusOr<Base64BinarySeparatorStride> extracted =
      ExtractOnlyMatchingExtension<Base64BinarySeparatorStride>(composition);

  EXPECT_FALSE(extracted.status().ok());
  EXPECT_EQ(absl::StatusCode::kNotFound, extracted.status().code());
}

TEST(ExtensionsTest, ExtractOnlyMatchingExtensionMultipleFound) {
  Composition composition = PARSE_VALID_STU3_PROTO(R"pb(
    id { value: "1" }
    status { value: FINAL }
    subject { patient_id { value: "P0" } }
    encounter { encounter_id { value: "2" } }
    author { practitioner_id { value: "3" } }
    title { value: "Note" }
    type {
      coding {
        system { value: "type:system" }
        code { value: "RADIOLOGY" }
      }
    }
    date {
      value_us: 4608099660000000
      timezone: "America/New_York"
      precision: SECOND
    }
    extension {
      url {
        value: "https://g.co/fhir/StructureDefinition/base64Binary-separatorStride"
      }
      extension {
        url { value: "separator" }
        value { string_value { value: "!" } }
      }
      extension {
        url { value: "stride" }
        value { positive_int { value: 5 } }
      }
    }
    extension {
      url {
        value: "https://g.co/fhir/StructureDefinition/base64Binary-separatorStride"
      }
      extension {
        url { value: "separator" }
        value { string_value { value: "*" } }
      }
      extension {
        url { value: "stride" }
        value { positive_int { value: 6 } }
      }
    }
  )pb");

  StatusOr<Base64BinarySeparatorStride> extracted =
      ExtractOnlyMatchingExtension<Base64BinarySeparatorStride>(composition);

  EXPECT_FALSE(extracted.status().ok());
  EXPECT_EQ(absl::StatusCode::kInvalidArgument, extracted.status().code());
}

}  // namespace

}  // namespace r4
}  // namespace fhir
}  // namespace google
