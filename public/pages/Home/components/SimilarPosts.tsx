import React, { useState, useEffect, useRef } from "react"
import { Post, Tag, CurrentUser } from "@fider/models"
import { ListPosts } from "./ListPosts"
import { actions } from "@fider/services"

import { i18n } from "@lingui/core"

import "./SimilarPosts.scss"

interface SimilarPostsProps {
  title: string
  tags: Tag[]
  user?: CurrentUser
}

export const SimilarPosts: React.FC<SimilarPostsProps> = (props) => {
  const [title, setTitle] = useState(props.title)
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [isVisible, setIsVisible] = useState(false)
  const timerRef = useRef<number>()

  useEffect(() => {
    loadSimilarPosts()

    return () => {
      window.clearTimeout(timerRef.current)
    }
  }, [])

  useEffect(() => {
    if (props.title !== title) {
      setTitle(props.title)
      setLoading(true)
    }
  }, [props.title])

  useEffect(() => {
    window.clearTimeout(timerRef.current)
    timerRef.current = window.setTimeout(loadSimilarPosts, 500)

    return () => {
      window.clearTimeout(timerRef.current)
    }
  }, [title])

  const preprocessSearchQuery = (query: string) => {
    const noiseWords = ["add", "support", "for", "implement", "create", "make", "allow", "enable", "provide"]
    return query
      .split(" ")
      .filter((x) => !noiseWords.includes(x))
      .join(" ")
  }

  const loadSimilarPosts = () => {
    if (loading) {
      if (title.length < 2) {
        setLoading(false)
        setIsVisible(false)
      } else {
        const query = preprocessSearchQuery(title)
        console.log("Query:", query)
        if (query.length < 2) {
          setLoading(false)
          setIsVisible(false)
          return
        }
        actions.findSimilarPosts(query).then((x) => {
          if (x.ok) {
            setLoading(false)
            setIsVisible(x.data.length > 0)
            setPosts(x.data)
          }
        })
      }
    }
  }

  const title_text = i18n._("home.similar.title", { message: "We have similar posts, is your idea already on the list?" })

  const animationClass = isVisible ? "similar-posts-visible" : "similar-posts-hidden"

  return (
    <>
      <div className={`similar-posts-container overflow-auto ${animationClass}`}>
        <div className="mb-4 text-gray-700">{title_text}</div>
        <div className="mb-6">
          <ListPosts posts={posts} tags={props.tags} emptyText="" minimalView={true} />
        </div>
      </div>
    </>
  )
}
