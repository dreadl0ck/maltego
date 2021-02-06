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
	"github.com/davecgh/go-spew/spew"
	"strings"
	"testing"
)

func TestParseMaltegoToTDS(t *testing.T) {

	var (
		tr = &Transform{}
		// Sample request XML going from Maltego client to TDS when running the example "DNSToIP" Transform.
		maltegoToTDS = `<MaltegoMessage>
		<MaltegoTransformRequestMessage>
			<Entities>
				<Entity Type="DNSName">
					<Genealogy>
						<Type Name="maltego.DNSName" OldName="DNSName"/>
					</Genealogy>
					<AdditionalFields>
						<Field Name="fqdn" DisplayName="DNS Name">alpine.paterva.com</Field>
					</AdditionalFields>
					<Value>alpine.paterva.com</Value>
					<Weight>0</Weight>
				</Entity>
			</Entities>
			<Limits SoftLimit="256" HardLimit="256"/>
		</MaltegoTransformRequestMessage>
	</MaltegoMessage>`
	)

	err := xml.Unmarshal([]byte(maltegoToTDS), tr)
	if err != nil {
		t.Fatal(err)
	}

	if tr.RequestMessage == nil || len(tr.RequestMessage.Entities.Items) != 1 {
		parseFailure(t, " len(tr.RequestMessage.Entities.Items) != 1", maltegoToTDS, tr)
	}

	if strings.TrimSpace(tr.RequestMessage.Entities.Items[0].Value) != "alpine.paterva.com" {
		parseFailure(t, "tr.RequestMessage.Entities.Items[0].Value != alpine.paterva.com", maltegoToTDS, tr)
	}

	if strings.TrimSpace(tr.RequestMessage.Entities.Items[0].Type) != "DNSName" {
		parseFailure(t, "tr.RequestMessage.Entities.Items[0].Type != DNSName", maltegoToTDS, tr)
	}

	if len(tr.RequestMessage.Entities.Items[0].Fields.Items) != 1 {
		parseFailure(t, "len(tr.RequestMessage.Entities.Items[0].Fields.Items) != 1", maltegoToTDS, tr)
	}

	if strings.TrimSpace(tr.RequestMessage.Entities.Items[0].Genealogy.Type.Name) != "maltego.DNSName" {
		parseFailure(t, "tr.RequestMessage.Entities.Items[0].Genealogy.Type.Name != maltego.DNSName", maltegoToTDS, tr)
	}

	if strings.TrimSpace(tr.RequestMessage.Entities.Items[0].Genealogy.Type.OldName) != "DNSName" {
		parseFailure(t, "tr.RequestMessage.Entities.Items[0].Genealogy.Type.OldName != DNSName", maltegoToTDS, tr)
	}

	if strings.TrimSpace(tr.RequestMessage.Entities.Items[0].Fields.Items[0].Name) != "fqdn" {
		parseFailure(t, "tr.RequestMessage.Entities.Items[0].Fields.Items[0].Name != fqdn", maltegoToTDS, tr)
	}

	if strings.TrimSpace(tr.RequestMessage.Entities.Items[0].Fields.Items[0].Text) != "alpine.paterva.com" {
		parseFailure(t, "tr.RequestMessage.Entities.Items[0].Fields.Items[0].Text != alpine.paterva.com", maltegoToTDS, tr)
	}

	if strings.TrimSpace(tr.RequestMessage.Entities.Items[0].Fields.Items[0].DisplayName) != "DNS Name" {
		parseFailure(t, "tr.RequestMessage.Entities.Items[0].Fields.Items[0].DisplayName != DNS Name", maltegoToTDS, tr)
	}

	if strings.TrimSpace(tr.RequestMessage.Entities.Items[0].Weight) != "0" {
		parseFailure(t, "tr.RequestMessage.Entities.Items[0].Weight != 0", maltegoToTDS, tr)
	}

	if tr.RequestMessage.Limits.SoftLimit != "256" {
		parseFailure(t, "tr.RequestMessage.Limits.SoftLimit != 256", maltegoToTDS, tr)
	}

	if tr.RequestMessage.Limits.HardLimit != "256" {
		parseFailure(t, "tr.RequestMessage.Limits.HardLimit != 256", maltegoToTDS, tr)
	}
}

func TestParseTDSToMaltego(t *testing.T) {

	var (
		tr = &Transform{}
		// Sample response XML of the above request going from TDS to Maltego client when running the example "DNSToIP" Transform.
		tdsToMaltego = `<MaltegoMessage>
		<MaltegoTransformResponseMessage>
			<Entities>
				<Entity Type="maltego.IPv4Address">
					<Value><![CDATA[173.230.156.137]]></Value>
					<Weight>100</Weight>
				</Entity>
			</Entities>
			<UIMessages>
				<UIMessage MessageType="Inform">Slider value is at: 256</UIMessage>
			</UIMessages>
		</MaltegoTransformResponseMessage>
	</MaltegoMessage>`
	)

	err := xml.Unmarshal([]byte(tdsToMaltego), tr)
	if err != nil {
		t.Fatal(err)
	}

	if tr.ResponseMessage == nil || len(tr.ResponseMessage.UIMessages.Items) != 1 {
		parseFailure(t, "len(tr.ResponseMessage.UIMessages.Items) != 1", tdsToMaltego, tr)
	}

	if tr.ResponseMessage == nil || len(tr.ResponseMessage.Entities.Items) != 1 {
		parseFailure(t, "len(message.ResponseMessage.Entities.Items) != 1", tdsToMaltego, tr)
	}

	if strings.TrimSpace(tr.ResponseMessage.Entities.Items[0].Type) != "maltego.IPv4Address" {
		parseFailure(t, "tr.ResponseMessage.Entities.Items[0].Type != maltego.IPv4Address", tdsToMaltego, tr)
	}

	if strings.TrimSpace(tr.ResponseMessage.Entities.Items[0].Value) != "173.230.156.137" {
		parseFailure(t, "tr.ResponseMessage.Entities.Items[0].Value != 173.230.156.137", tdsToMaltego, tr)
	}

	if strings.TrimSpace(tr.ResponseMessage.Entities.Items[0].Weight) != "100" {
		parseFailure(t, "tr.ResponseMessage.Entities.Items[0].Weight != 100", tdsToMaltego, tr)
	}

	if strings.TrimSpace(tr.ResponseMessage.UIMessages.Items[0].MessageType) != "Inform" {
		parseFailure(t, "tr.ResponseMessage.UIMessages.Items[0].MessageType != Inform", tdsToMaltego, tr)
	}

	if strings.TrimSpace(tr.ResponseMessage.UIMessages.Items[0].Text) != "Slider value is at: 256" {
		parseFailure(t, "tr.ResponseMessage.UIMessages.Items[0].Text != Slider value is at: 256", tdsToMaltego, tr)
	}
}

func parseFailure(t *testing.T, reason, expected string, transform *Transform) {
	fmt.Println("=========== OUTPUT ==========")
	spew.Dump(transform)
	fmt.Println("=========== EXPECTED ==========")
	fmt.Println(expected)
	t.Fatal("unexpected output: " + reason)
}

// helper to compare output against expected result
// and help diagnose issues.
func compare(t *testing.T, data []byte, exp string) {
	if string(data) != exp {
		fmt.Println("=========== OUTPUT ==========")
		fmt.Println(string(data))
		fmt.Println("=========== EXPECTED ==========")
		fmt.Println(exp)
		fmt.Println("=========== DETAIL ==========")
		for i, c := range string(data) {
			if string(exp[i]) != string(c) {
				fmt.Println("\n", i, ":", string(exp[i]), "!=", string(c))
				t.Fatal("unexpected out")
			} else {
				fmt.Print(string(c))
			}
		}
		t.Fatal("unexpected out")
	}
}

func TestTransformFromStructure(t *testing.T) {
	m := Transform{
		ResponseMessage: &ResponseMessage{
			Entities: Entities{
				Items: []*Entity{
					{
						Type:  "type",
						Value: "value",
					},
					{
						Type:  "type2",
						Value: "value2",
					},
				},
			},
			UIMessages: UIMessages{
				Items: []*UIMessage{
					{
						Text:        "text",
						MessageType: UIMessageDebug,
					},
					{
						Text:        "text2",
						MessageType: UIMessageDebug,
					},
				},
			},
		},
	}

	data, err := xml.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	exp := `<MaltegoMessage><MaltegoTransformResponseMessage><Entities><Entity Type="type"><Value>value</Value><Weight></Weight></Entity><Entity Type="type2"><Value>value2</Value><Weight></Weight></Entity></Entities><UIMessages><UIMessage MessageType="Debug">text</UIMessage><UIMessage MessageType="Debug">text2</UIMessage></UIMessages></MaltegoTransformResponseMessage></MaltegoMessage>`

	compare(t, data, exp)
}

func TestTransformViaHelpers(t *testing.T) {
	trx := Transform{}

	trx.AddEntity("type", "value")
	trx.AddEntity("type2", "value2")

	trx.AddUIMessage("message", UIMessageDebug)
	trx.AddUIMessage("message2", UIMessageDebug)

	out := `<MaltegoMessage><MaltegoTransformResponseMessage><Entities><Entity Type="type"><Value>value</Value><Weight>100</Weight></Entity><Entity Type="type2"><Value>value2</Value><Weight>100</Weight></Entity></Entities><UIMessages><UIMessage MessageType="Debug">message</UIMessage><UIMessage MessageType="Debug">message2</UIMessage></UIMessages></MaltegoTransformResponseMessage></MaltegoMessage>`
	compare(t, []byte(trx.ReturnOutput()), out)
}

func TestTransformEntity(t *testing.T) {
	trx := Entity{
		Type:    "type",
		Value:   "value",
		IconURL: "http://asdf.com",
		Weight:  "10",
		Info: &DisplayInformation{
			Labels: []*DisplayLabel{
				NewDisplayLabel("name", "text"),
				NewDisplayLabel("name2", "text2"),
			},
		},
	}

	data, err := xml.Marshal(trx)
	if err != nil {
		t.Fatal(err)
	}

	exp := `<Entity Type="type"><Value>value</Value><Weight>10</Weight><DisplayInformation><Label Name="text" Type="text/html"><![CDATA[name]]></Label><Label Name="text2" Type="text/html"><![CDATA[name2]]></Label></DisplayInformation><IconURL>http://asdf.com</IconURL></Entity>`
	compare(t, data, exp)
}

func TestTransformException(t *testing.T) {
	msg := Transform{
		ExceptionMessage: &ExceptionMessage{
			Exceptions: Exceptions{
				Items: []*Exception{
					{
						Text: "oops",
						Code: "errorCode",
					},
				},
			},
		},
	}

	data, err := xml.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	exp := `<MaltegoMessage><MaltegoTransformExceptionMessage><Exceptions><Exception code="errorCode">oops</Exception></Exceptions></MaltegoTransformExceptionMessage></MaltegoMessage>`
	compare(t, data, exp)
}

func TestTransformThrowException(t *testing.T) {
	trx := Transform{}
	trx.AddException("oops", "errorCode")

	out := `<MaltegoMessage><MaltegoTransformExceptionMessage><Exceptions><Exception code="errorCode">oops</Exception></Exceptions></MaltegoTransformExceptionMessage></MaltegoMessage>`
	compare(t, []byte(trx.ThrowExceptions()), out)
}

func TestLabel(t *testing.T) {
	l := NewDisplayLabel("text", "name")

	data, err := xml.Marshal(l)
	if err != nil {
		t.Fatal(err)
	}

	str := `<Label Name="name" Type="text/html"><![CDATA[text]]></Label>`
	compare(t, data, str)
}

func TestEscape(t *testing.T) {
	fmt.Println(EscapeText("\n"))
}