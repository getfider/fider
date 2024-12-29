import React from "react"
import { VStack } from "@fider/components/layout"
import { actions } from "@fider/services"
import { Post } from "@fider/models"

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
      <h1 className="text-title">Most Wanted</h1>
      {posts.length > 0 &&
        posts.map((post) => (
          <div key={post.id} className="pt-1">
            <span className="text-medium">{post.title}</span>
            <div className="pt-1 flex">
              <span className="text-muted" style={{ width: "150px" }}>{post.commentsCount} comments</span>
              <span className="text-muted">{post.votesCount} votes</span>
            </div>
          </div>
        ))}
    </VStack>
  )
}
