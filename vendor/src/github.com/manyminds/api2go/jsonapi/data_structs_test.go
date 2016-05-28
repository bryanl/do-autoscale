package jsonapi

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSONAPI Struct tests", func() {
	Context("Testing array and object data payload", func() {
		It("detects object payload", func() {
			sampleJSON := `
			{
				"data":
				{
					"type": "test",
					"id": "1",
					"attributes": {"foo": "bar"},
					"relationships": {
						"author": {
							"data": {"type": "author", "id": "1"}
						}
					}
				}
			}
`
			expectedData := &Data{
				Type:       "test",
				ID:         "1",
				Attributes: json.RawMessage([]byte(`{"foo": "bar"}`)),
				Relationships: map[string]Relationship{
					"author": {
						Data: &RelationshipDataContainer{
							DataObject: &RelationshipData{
								Type: "author",
								ID:   "1",
							},
						},
					},
				},
			}

			target := Document{}

			err := json.Unmarshal([]byte(sampleJSON), &target)
			Expect(err).ToNot(HaveOccurred())
			Expect(target.Data.DataObject).To(Equal(expectedData))
		})

		It("detects array payload", func() {
			sampleJSON := `
			{
				"data": [
					{
						"type": "test",
						"id": "1",
						"attributes": {"foo": "bar"},
						"relationships": {
							"comments": {
								"data": [
									{"type": "comments", "id": "1"},
									{"type": "comments", "id": "2"}
								]
							}
						}
					}
				]
			}
`
			expectedData := Data{
				Type:       "test",
				ID:         "1",
				Attributes: json.RawMessage([]byte(`{"foo": "bar"}`)),
				Relationships: map[string]Relationship{
					"comments": {
						Data: &RelationshipDataContainer{
							DataArray: []RelationshipData{
								{
									Type: "comments",
									ID:   "1",
								},
								{
									Type: "comments",
									ID:   "2",
								},
							},
						},
					},
				},
			}

			target := Document{}

			err := json.Unmarshal([]byte(sampleJSON), &target)
			Expect(err).ToNot(HaveOccurred())
			Expect(target.Data.DataArray).To(Equal([]Data{expectedData}))
		})
	})

	It("return an error for invalid relationship data format", func() {
		sampleJSON := `
		{
			"data": [
			{
				"type": "test",
				"id": "1",
				"attributes": {"foo": "bar"},
				"relationships": {
					"comments": {
						"data": "foo"
					}
				}
			}
			]
		}
		`

		target := Document{}

		err := json.Unmarshal([]byte(sampleJSON), &target)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("Invalid json for relationship data array/object"))
	})

	It("creates an empty slice for empty to-many relationships and nil for empty toOne", func() {
		sampleJSON := `
			{
				"data": [
					{
						"type": "test",
						"id": "1",
						"attributes": {"foo": "bar"},
						"relationships": {
							"comments": {
								"data": []
							},
							"author": {
								"data": null
							}
						}
					}
				]
			}
`
		expectedData := Data{
			Type:       "test",
			ID:         "1",
			Attributes: json.RawMessage([]byte(`{"foo": "bar"}`)),
			Relationships: map[string]Relationship{
				"comments": {
					Data: &RelationshipDataContainer{
						DataArray: []RelationshipData{},
					},
				},
				"author": {
					Data: nil,
				},
			},
		}

		target := Document{}

		err := json.Unmarshal([]byte(sampleJSON), &target)
		Expect(err).ToNot(HaveOccurred())
		Expect(target.Data.DataArray).To(Equal([]Data{expectedData}))
	})
})
