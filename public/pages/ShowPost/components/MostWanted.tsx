import React from "react"
import { VStack } from "@fider/components/layout"
import { actions } from "@fider/services"
import { Post } from "@fider/models"
import { Trans } from "@lingui/react"

export const MostWanted = () => {
  const [posts, setPosts] = React.useState<Post[]>([])

  React.useEffect(() => {
    const fetchPosts = async () => {
      const result = await actions.searchPosts({ view: "most-wanted", limit: 5 })
      setPosts(result.data)
    }
    fetchPosts()
  }, [])

  return (
    <VStack spacing={4}>
      <h1 className="text-title">
        <Trans id="home.postfilter.option.mostwanted">Most Wanted</Trans>
      </h1>
      {posts.length > 0 &&
        posts.map((post) => (
          <div key={post.id} className="pt-1 hover">
            <span className="text-medium">{post.title}</span>
            <div className="pt-1 flex">
              <span className="text-muted" style={{ width: "150px" }}>
                <Trans id="showpost.mostwanted.comments" values={{ count: post.commentsCount }}>
                  comments
                </Trans>
              </span>
              <span className="text-muted">
                <Trans id="showpost.mostwanted.votes" values={{ count: post.votesCount }} message="{count, plural, one {# votes} other {# votes}}">
                  votes
                </Trans>
              </span>
            </div>
          </div>
        ))}
    </VStack>
  )
}
