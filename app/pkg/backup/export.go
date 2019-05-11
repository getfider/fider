package backup

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
)

func exportTable(ctx context.Context, tableName string) ([]byte, error) {
	trx := ctx.Value(app.TransactionCtxKey).(*dbx.Trx)
	tenant, _ := ctx.Value(app.TenantCtxKey).(*models.Tenant)
	columnName := "tenant_id"
	if tableName == "tenants" {
		columnName = "id"
	}

	rows, err := trx.Query(fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", tableName, columnName), tenant.ID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(jsonify(rows))
}

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path"
// 	"strings"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/credentials"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/s3"

// 	"github.com/jmoiron/sqlx"
// 	_ "github.com/lib/pq"
// )

// var db *sqlx.DB
// var tenantID int
// var s3Client *s3.S3

// func init() {
// 	var err error
// 	db, err = sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	tenantID = findTenantID(os.Getenv("TENANT"))

// 	os.Getenv("S3_BUCKET")

// 	awsSession := session.New(&aws.Config{
// 		Credentials:      credentials.NewStaticCredentials(os.Getenv("S3_ACCESS_KEY_ID"), os.Getenv("S3_SECRET_ACCESS_KEY"), ""),
// 		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT_URL")),
// 		Region:           aws.String(os.Getenv("S3_REGION")),
// 		DisableSSL:       aws.Bool(strings.HasSuffix(os.Getenv("S3_ENDPOINT_URL"), "http://")),
// 		S3ForcePathStyle: aws.Bool(true),
// 	})
// 	s3Client = s3.New(awsSession)
// }

// func main() {
// 	os.RemoveAll(path.Join("./exports", os.Getenv("TENANT")))
// 	export("attachments")
// 	export("comments")
// 	export("email_verifications")
// 	export("notifications")
// 	export("oauth_providers")
// 	export("posts")
// 	export("post_subscribers")
// 	export("post_tags")
// 	export("post_votes")
// 	export("tags")
// 	export("tenants")
// 	export("user_providers")
// 	export("users")
// 	export("user_settings")
// 	exportBlobs()
// }

// func exportBlobs() {
// 	response, err := s3Client.ListObjects(&s3.ListObjectsInput{
// 		Bucket:  aws.String(os.Getenv("S3_BUCKET")),
// 		MaxKeys: aws.Int64(3000),
// 		Prefix:  aws.String(fmt.Sprintf("tenants/%d/", tenantID)),
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, item := range response.Contents {
// 		key := *item.Key
// 		resp, err := s3Client.GetObject(&s3.GetObjectInput{
// 			Bucket: aws.String(os.Getenv("S3_BUCKET")),
// 			Key:    aws.String(key),
// 		})
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer resp.Body.Close()
// 		bytes, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			panic(err)
// 		}

// 		keyPart := strings.Split(key, "/")
// 		folder := path.Join("./exports", os.Getenv("TENANT"), "blobs", keyPart[2])
// 		os.MkdirAll(folder, 0777)

// 		err = ioutil.WriteFile(path.Join(folder, keyPart[3]), bytes, 0777)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

// func export(tableName string) {
// 	columnName := "tenant_id"
// 	if tableName == "tenants" {
// 		columnName = "id"
// 	}
// 	rows, err := db.Queryx(fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", tableName, columnName), tenantID)
// 	if err != nil {
// 		panic(err)
// 	}

// 	allResults := make([]map[string]interface{}, 0)

// 	for rows.Next() {
// 		results := make(map[string]interface{})
// 		err := rows.MapScan(results)
// 		if err != nil {
// 			panic(err)
// 		}
// 		allResults = append(allResults, results)
// 	}

// 	result, err := json.Marshal(allResults)
// 	if err != nil {
// 		panic(err)
// 	}

// 	folder := path.Join("./exports", os.Getenv("TENANT"))

// 	os.MkdirAll(folder, 0777)

// 	err = ioutil.WriteFile(path.Join(folder, tableName+".json"), result, 0777)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func findTenantID(tenant string) int {
// 	var tenantID int
// 	err := db.QueryRow("SELECT id FROM tenants WHERE subdomain = $1", tenant).Scan(&tenantID)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return tenantID
// }

func jsonify(rows *sql.Rows) []map[string]interface{} {
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]interface{}, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	results := make(map[string]interface{})
	allResults := make([]map[string]interface{}, 0)

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		for i, value := range values {
			switch value.(type) {
			case nil:
				results[columns[i]] = nil

			case []byte:
				s := string(value.([]byte))
				x, err := strconv.Atoi(s)

				if err != nil {
					results[columns[i]] = s
				} else {
					results[columns[i]] = x
				}

			default:
				results[columns[i]] = value
			}
		}

		allResults = append(allResults, results)
	}

	return allResults
}
