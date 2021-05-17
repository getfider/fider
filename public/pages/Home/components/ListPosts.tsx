import React from "react"
import { Post, Tag, CurrentUser } from "@fider/models"
import { ShowTag, ShowPostResponse, VoteCounter, Markdown, Icon } from "@fider/components"
import IconChatAlt2 from "@fider/assets/images/heroicons-chat-alt-2.svg"
import { HStack, VStack } from "@fider/components/layout"

interface ListPostsProps {
  posts?: Post[]
  tags: Tag[]
  emptyText: string
}

const ListPostItem = (props: { post: Post; user?: CurrentUser; tags: Tag[] }) => {
  return (
    <HStack center={true}>
      <div className="align-self-start">
        <VoteCounter post={props.post} />
      </div>
      <VStack className="w-full" spacing={2}>
        <HStack justify="between">
          <a className="text-title hover:text-primary-base" href={`/posts/${props.post.number}/${props.post.slug}`}>
            {props.post.title}
          </a>
          {props.post.commentsCount > 0 && (
            <HStack className="text-muted">
              {props.post.commentsCount} <Icon sprite={IconChatAlt2} className="h-4 ml-1" />
            </HStack>
          )}
        </HStack>
        <Markdown className="text-gray-600" maxLength={300} text={props.post.description} style="plainText" />
        <ShowPostResponse status={props.post.status} response={props.post.response} />
        {props.tags.length >= 1 && (
          <HStack className="flex-wrap">
            {props.tags.map((tag) => (
              <ShowTag key={tag.id} tag={tag} />
            ))}
          </HStack>
        )}
      </VStack>
    </HStack>
  )
}

export const ListPosts = (props: ListPostsProps) => {
  if (!props.posts) {
    return null
  }

  if (props.posts.length === 0) {
    return <p className="text-center">{props.emptyText}</p>
  }

  return (
    <VStack spacing={4} divide={true} center={true} className="test">
      {props.posts.map((post) => (
        <ListPostItem key={post.id} post={post} tags={props.tags.filter((tag) => post.tags.indexOf(tag.slug) >= 0)} />
      ))}
    </VStack>
  )
}
