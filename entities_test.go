/*
 * MALTEGO - Go package that provides datastructures for interacting with the Maltego graphical link analysis tool.
 * Copyright (c) 2021 Philipp Mieden <dreadl0ck [at] protonmail [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package maltego

import (
	"encoding/xml"
	"fmt"
	"testing"
)

var maltegoEntities = []EntityCoreInfo{
	{"ContentType", "category", "A MIME type describes different multi-media formats", "", nil},
	{"Email", "mail_outline", "An email message", "maltego.Email", nil},
	{"Interface", "router", "A network interface", "", []PropertyField{NewRequiredStringField("properties.interface", "Name of the network interface"), NewStringField("snaplen", "snap length for ethernet frames in bytes, default: 1514"), NewStringField("bpf", "berkeley packet filter to apply")}},
	{"PCAP", "sd_storage", "A packet capture dump file", "", []PropertyField{NewRequiredStringField("path", "Absolute path to the PCAP file")}},
}

func compareGeneratedXML(data []byte, expected string, t *testing.T) {

	if string(data) != expected {
		fmt.Println("-------------------RESULT--------------------------")
		fmt.Println(string(data))
		fmt.Println("------------------------------------------------")

		fmt.Println("-------------------EXPECTED--------------------------")
		fmt.Println(expected)
		fmt.Println("------------------------------------------------")

		t.Fatal("unexpected output")
	}
}

func TestGenerateTestEntityXMLEntity(t *testing.T) {
	expected := `<MaltegoEntity id="test.Entity" displayName="TestEntity" displayNamePlural="TestEntities" description="A test entity" category="Test" smallIconResource="Technology/WAN" largeIconResource="Technology/WAN" allowedRoot="true" conversionOrder="2147483647" visible="true">
   <Properties value="properties.test" displayValue="properties.test">
      <Groups></Groups>
      <Fields>
         <Field name="properties.test" type="string" nullable="true" hidden="false" readonly="false" description="" displayName="TestEntity">
            <SampleValue>-</SampleValue>
         </Field>
      </Fields>
   </Properties>
</MaltegoEntity>`
	e := MaltegoEntity{
		ID:                "test.Entity",
		DisplayName:       "TestEntity",
		DisplayNamePlural: "TestEntities",
		Description:       "A test entity",
		Category:          "Test",
		SmallIconResource: "Technology/WAN",
		LargeIconResource: "Technology/WAN",
		AllowedRoot:       true,
		ConversionOrder:   "2147483647",
		Visible:           true,
		Properties: EntityProperties{
			Value:        "properties.test",
			DisplayValue: "properties.test",
			Fields: Fields{
				Items: []PropertyField{
					{
						Name:        "properties.test",
						Type:        "string",
						Nullable:    true,
						Hidden:      false,
						Readonly:    false,
						Description: "",
						DisplayName: "TestEntity",
						SampleValue: "-",
					},
				},
			},
		},
	}

	data, err := xml.MarshalIndent(e, "", "   ")
	if err != nil {
		t.Fatal(err)
	}

	compareGeneratedXML(data, expected, t)
}

func TestToTransformDisplayName(t *testing.T) {
	res := ToTransformDisplayName("ToTCPServices")
	if res != "To TCP Services [NETCAP]" {
		t.Fatal("unexpected result", res)
	}

	res = ToTransformDisplayName("ToDHCP")
	if res != "To DHCP [NETCAP]" {
		t.Fatal("unexpected result", res)
	}

	res = ToTransformDisplayName("ToServerNameIndicators")
	if res != "To Server Name Indicators [NETCAP]" {
		t.Fatal("unexpected result", res)
	}

	res = ToTransformDisplayName("ToURLsForHost")
	if res != "To URLs For Host [NETCAP]" {
		t.Fatal("unexpected result", res)
	}

	res = ToTransformDisplayName("ToSourceIPs")
	if res != "To Source IPs [NETCAP]" {
		t.Fatal("unexpected result", res)
	}
}
