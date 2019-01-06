package stripe

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"path/filepath"
)

// FilePurpose is the purpose of a particular file.
type FilePurpose string

// List of values that FilePurpose can take.
const (
	FilePurposeBusinessLogo          FilePurpose = "business_logo"
	FilePurposeCustomerSignature     FilePurpose = "customer_signature"
	FilePurposeDisputeEvidence       FilePurpose = "dispute_evidence"
	FilePurposeFinanceReportRun      FilePurpose = "finance_report_run"
	FilePurposeFoundersStockDocument FilePurpose = "founders_stock_document"
	FilePurposeIdentityDocument      FilePurpose = "identity_document"
	FilePurposePCIDocument           FilePurpose = "pci_document"
	FilePurposeSigmaScheduledQuery   FilePurpose = "sigma_scheduled_query"
	FilePurposeTaxDocumentUserUpload FilePurpose = "tax_document_user_upload"
)

// FileParams is the set of parameters that can be used when creating a file.
// For more details see https://stripe.com/docs/api#create_file.
type FileParams struct {
	Params `form:"*"`

	// FileReader is a reader with the contents of the file that should be uploaded.
	FileReader io.Reader

	// Filename is just the name of the file without path information.
	Filename *string

	Purpose *string
}

// FileListParams is the set of parameters that can be used when listing
// files. For more details see https://stripe.com/docs/api#list_files.
type FileListParams struct {
	ListParams   `form:"*"`
	Created      *int64            `form:"created"`
	CreatedRange *RangeQueryParams `form:"created"`
	Purpose      *string           `form:"purpose"`
}

// File is the resource representing a Stripe file.
// For more details see https://stripe.com/docs/api#file_object.
type File struct {
	Created  int64         `json:"created"`
	ID       string        `json:"id"`
	Filename string        `json:"filename"`
	Links    *FileLinkList `json:"links"`
	Purpose  FilePurpose   `json:"purpose"`
	Size     int64         `json:"size"`
	Type     string        `json:"type"`
	URL      string        `json:"url"`
}

// FileList is a list of files as retrieved from a list endpoint.
type FileList struct {
	ListMeta
	Data []*File `json:"data"`
}

// GetBody gets an appropriate multipart form payload to use in a request body
// to create a new file.
func (f *FileParams) GetBody() (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if f.Purpose != nil {
		err := writer.WriteField("purpose", StringValue(f.Purpose))
		if err != nil {
			return nil, "", err
		}
	}

	if f.FileReader != nil && f.Filename != nil {
		part, err := writer.CreateFormFile("file", filepath.Base(StringValue(f.Filename)))
		if err != nil {
			return nil, "", err
		}

		_, err = io.Copy(part, f.FileReader)
		if err != nil {
			return nil, "", err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.Boundary(), nil
}

// UnmarshalJSON handles deserialization of a File.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (f *File) UnmarshalJSON(data []byte) error {
	if id, ok := ParseID(data); ok {
		f.ID = id
		return nil
	}

	type file File
	var v file
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = File(v)
	return nil
}
