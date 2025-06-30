import React, { useEffect, useState } from "react"
import { Loader } from "@fider/components"
import ShowPostPage from "@fider/pages/ShowPost/ShowPost.page"
import { getPostDetails, PostDetailsResponse } from "@fider/services/actions/post-details"

interface PostDetailsProps {
  postNumber: number
}

export const PostDetails: React.FC<PostDetailsProps> = (props) => {
  const [loading, setLoading] = useState(true)
  const [data, setData] = useState<PostDetailsResponse | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (props.postNumber) {
      setLoading(true)
      setError(null)

      getPostDetails(props.postNumber).then((response) => {
        setLoading(false)
        if (response.ok) {
          setData(response.data)
        } else {
          setError("Failed to load post")
        }
      })
    }
  }, [props.postNumber])

  if (loading) {
    return (
      <div className="p-6 text-center">
        <Loader />
      </div>
    )
  }

  if (error) {
    return <div className="p-6 text-center text-red-500">{error}</div>
  }

  if (!data) {
    return null
  }

  return (
    <ShowPostPage post={data.post} comments={data.comments} tags={data.tags} votes={data.votes} subscribed={data.subscribed} attachments={data.attachments} />
  )
}
