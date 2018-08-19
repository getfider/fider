package apiv1

import (
	"time"

	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
)

// Authorize exchanges given input for an Auth Token
func Authorize() web.HandlerFunc {
	return func(c web.Context) error {
		input := new(actions.APIAuthorize)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		token, err := jwt.Encode(jwt.FiderClaims{
			UserID:    input.User.ID,
			UserName:  input.User.Name,
			UserEmail: input.User.Email,
			Origin:    jwt.FiderClaimsOriginAPI,
			Metadata: jwt.Metadata{
				ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			},
		})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"authToken": token,
		})
	}
}

// * User Flow

// Request:
// {
// 	"apiKey": "..."
// }

// Response:
// {
// 	"authToken": "..."
// }

// * Application Flow
// Request:
// {
// 	"appID": "..."
// 	"appUserToken": "..."
// }

// Response:
// {
// 	"authToken": "..."
// }
