package webhook

import (
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
)

// Props is a map of key:value
type Props map[string]interface{}

// SetUser describe the user prefixed by "keyPrefix"
func (p Props) SetUser(user *entity.User, keyPrefix string) Props {
	if user != nil {
		p[keyPrefix+"_id"] = user.ID
		p[keyPrefix+"_name"] = user.Name
		p[keyPrefix+"_email"] = user.Email
		p[keyPrefix+"_role"] = user.Role.String()
		p[keyPrefix+"_avatar"] = user.AvatarURL
	}
	return p
}

// SetTenant describe the tenant prefixed by "keyPrefix"
func (p Props) SetTenant(tenant *entity.Tenant, keyPrefix, baseURL, logoURL string) Props {
	if tenant != nil {
		p[keyPrefix+"_id"] = tenant.ID
		p[keyPrefix+"_name"] = tenant.Name
		p[keyPrefix+"_subdomain"] = tenant.Subdomain
		p[keyPrefix+"_status"] = tenant.Status.String()
		p[keyPrefix+"_locale"] = tenant.Locale
		p[keyPrefix+"_url"] = baseURL
		p[keyPrefix+"_logo"] = logoURL
	}
	return p
}

// SetPost describe the post prefixed by "keyPrefix"
func (p Props) SetPost(post *entity.Post, keyPrefix, baseURL string, includeAllFields, includeAuthor bool) Props {
	if post != nil {
		p[keyPrefix+"_id"] = post.ID
		p[keyPrefix+"_number"] = post.Number
		p[keyPrefix+"_title"] = post.Title
		p[keyPrefix+"_slug"] = post.Slug
		p[keyPrefix+"_description"] = post.Description
		p[keyPrefix+"_created_at"] = post.CreatedAt
		p[keyPrefix+"_url"] = post.Url(baseURL)

		if includeAuthor {
			p.SetUser(post.User, keyPrefix+"_author")
		}

		if includeAllFields {
			postResponse := post.Response
			p[keyPrefix+"_votes"] = post.VotesCount
			p[keyPrefix+"_comments"] = post.CommentsCount
			p[keyPrefix+"_status"] = post.Status.Name()
			p[keyPrefix+"_tags"] = post.Tags
			p[keyPrefix+"_response"] = postResponse != nil

			if postResponse != nil {
				keyPrefix := keyPrefix + "_response"
				p[keyPrefix+"_text"] = postResponse.Text
				p[keyPrefix+"_responded_at"] = postResponse.RespondedAt
				p.SetUser(postResponse.User, keyPrefix+"_author")

				originalPost := postResponse.Original
				if post.Status == enum.PostDuplicate && originalPost != nil {
					keyPrefix := keyPrefix + "_original"
					p[keyPrefix+"_number"] = originalPost.Number
					p[keyPrefix+"_title"] = originalPost.Title
					p[keyPrefix+"_slug"] = originalPost.Slug
					p[keyPrefix+"_status"] = originalPost.Status.Name()
					p[keyPrefix+"_url"] = originalPost.Url(baseURL)
				}
			}
		}
	}
	return p
}
