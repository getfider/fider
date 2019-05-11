package backup

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"

	"github.com/getfider/fider/app/pkg/errors"
)

func Create(ctx context.Context) (*bytes.Buffer, error) {

	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)

	for _, tableName := range []string{
		"attachments",
		"comments",
		"email_verifications",
		"notifications",
		"oauth_providers",
		"posts",
		"post_subscribers",
		"post_tags",
		"post_votes",
		"tags",
		"tenants",
		"user_providers",
		"users",
		"user_settings",
	} {
		err := addTableDataToZipFile(ctx, zipWriter, tableName)
		if err != nil {
			return nil, err
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close zip file")
	}

	return buffer, nil
}

func addTableDataToZipFile(ctx context.Context, zipWriter *zip.Writer, tableName string) error {
	tableData, err := exportTable(ctx, tableName)
	if err != nil {
		return errors.Wrap(err, "failed to export %s table", tableName)
	}

	fileWriter, err := zipWriter.Create(fmt.Sprintf("%s.json", tableName))
	if err != nil {
		return errors.Wrap(err, "failed to create %s.json in zip file", tableName)
	}
	_, err = fileWriter.Write(tableData)
	if err != nil {
		return errors.Wrap(err, "failed to write %s.json to zip file", tableName)
	}

	return nil
}
