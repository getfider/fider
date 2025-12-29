import React from "react"
import { Post, Tag, CurrentUser } from "@fider/models"
import { ShowTag, VoteCounter, Markdown, Icon, ResponseLozenge } from "@fider/components"
import IconChatAlt2 from "@fider/assets/images/heroicons-chat-alt-2.svg"
import { HStack, VStack } from "@fider/components/layout"

interface ListPostsProps {
  posts?: Post[]
  tags: Tag[]
  emptyText: string
  minimalView?: boolean
  onPostClick?: (postNumber: number, slug: string) => void
}

const ListPostItem = (props: { post: Post; user?: CurrentUser; tags: Tag[]; onPostClick?: (postNumber: number, slug: string) => void }) => {
  const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
    if (props.onPostClick) {
      e.preventDefault()
      props.onPostClick(props.post.number, props.post.slug)
    }
  }

  return (
    <HStack spacing={4} align="start" className="c-posts-container__post">
      <div>
        <VoteCounter post={props.post} />
      </div>
      <VStack className="w-full" spacing={2}>
        {props.post.status !== "open" && (
          <div className="mb-2 align-self-start">
            <ResponseLozenge status={props.post.status} response={props.post.response} size={"small"} />
          </div>
        )}
        <HStack justify="between">
          <a className="text-title text-break hover:text-primary-base" href={`/posts/${props.post.number}/${props.post.slug}`} onClick={handleClick}>
            {props.post.title}
          </a>
          {props.post.commentsCount > 0 && (
            <HStack className="text-muted">
              {props.post.commentsCount} <Icon sprite={IconChatAlt2} className="h-4 ml-1" />
            </HStack>
          )}
        </HStack>
        <Markdown className="c-posts-container__postdescription" maxLength={300} text={props.post.description} style="plainText" />
        {props.tags.length >= 1 && (
          <HStack spacing={0} className="gap-2 flex-wrap">
            {props.tags.map((tag) => (
              <ShowTag key={tag.id} tag={tag} link />
            ))}
          </HStack>
        )}
      </VStack>
    </HStack>
  )
}

const MinimalListPostItem = (props: { post: Post; tags: Tag[]; onPostClick?: (postNumber: number, slug: string) => void }) => {
  const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
    if (props.onPostClick) {
      e.preventDefault()
      props.onPostClick(props.post.number, props.post.slug)
    }
  }

  return (
    <HStack spacing={4} align="start" className="c-posts-container__post">
      <HStack className="w-full" justify="between" align="start">
        <a className="text-link" href={`/posts/${props.post.number}/${props.post.slug}`} onClick={handleClick}>
          {props.post.title}
        </a>
        {props.post.status !== "open" ? (
          <div>
            <ResponseLozenge status={props.post.status} response={props.post.response} size={"micro"} />
          </div>
        ) : (
          <span className="text-gray-700 text-sm">+{props.post.votesCount}</span>
        )}
      </HStack>
    </HStack>
  )
}

export const ListPosts = (props: ListPostsProps) => {
  const { minimalView = false } = props

  if (!props.posts) {
    return null
  }

  if (props.posts.length === 0) {
    return <p className="text-center">{props.emptyText}</p>
  }

  return (
    <>
      {minimalView ? (
        <VStack spacing={2}>
          {props.posts.map((post) => (
            <MinimalListPostItem key={post.id} post={post} tags={props.tags.filter((tag) => post.tags.indexOf(tag.slug) >= 0)} onPostClick={props.onPostClick} />
          ))}
        </VStack>
      ) : (
        <VStack spacing={4} divide>
          {props.posts.map((post) => (
            <ListPostItem key={post.id} post={post} tags={props.tags.filter((tag) => post.tags.indexOf(tag.slug) >= 0)} onPostClick={props.onPostClick} />
          ))}
        </VStack>
      )}
    </>
  )
}
