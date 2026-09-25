package v1

import (
	"time"

	"github.com/ep0ll/rebaze/specs"
)

// PatchSpec is the core structure that defines a patch artifact for an OCI image
// or other OCI artifact. It consolidates all components of a patch into a single,
// verifiable specification.
type PatchSpec struct {
	specs.Versioned
	// MediaType specifies the type of this document data structure e.g. `application/vnd.oci.image.patch.v1+json`
	MediaType string `json:"mediaType,omitempty"`

	// ArtifactType specifies the IANA media type of artifact when the manifest is used for an artifact.
	ArtifactType string    `json:"artifactType,omitempty"`
	Created      time.Time `json:"created"`
	Author       string    `json:"author,omitempty"`
	Description  string    `json:"description,omitempty"`
	// Subject is an optional link from the image manifest to another manifest forming an association between the image manifest and the other manifest.
	Subject *Descriptor `json:"subject,omitempty"`
	// Data is an embedding of the targeted content. This is encoded as a base64
	// string when marshalled to JSON (automatically, by encoding/json). If
	// present, Data can be used directly to avoid fetching the targeted content.
	Data    []byte            `json:"data,omitempty"`
	Audit   []PatchAuditEvent `json:"audit,omitempty"`
	Users   []User            `json:"users"`
	Teams   []Team            `json:"teams,omitempty"`
	Patches []LayerAction     `json:"patches"`
	// Annotations contains arbitrary metadata for the image index.
	Annotations map[string]string `json:"annotations,omitempty"`
}
