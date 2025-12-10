import React from "react"
import { Post, Tag, CurrentUser } from "@fider/models"
import { ShowTag, Markdown, Icon, ResponseLozenge } from "@fider/components"
import IconChatAlt2 from "@fider/assets/images/heroicons-chat-alt-2.svg"
import { HStack, VStack } from "@fider/components/layout"
import { Trans } from "@lingui/react/macro"

interface ListPostsProps {
  posts?: Post[]
  tags: Tag[]
  emptyText: string
  minimalView?: boolean
}

const ListPostItem = (props: { post: Post; user?: CurrentUser; tags: Tag[] }) => {
  return (
    <a href={`/posts/${props.post.number}/${props.post.slug}`} className="c-posts-container__post-link">
      <VStack className="c-posts-container__post w-full" spacing={4}>
        {props.post.status !== "open" && (
          <div className="mb-1 align-self-start">
            <ResponseLozenge status={props.post.status} response={props.post.response} size={"small"} />
          </div>
        )}
        <HStack justify="between" align="start">
          <h3 className="c-posts-container__post-title text-break">{props.post.title}</h3>
          {props.post.commentsCount > 0 && (
            <HStack spacing={1} className="c-posts-container__post-comments flex-shrink-0">
              <span>{props.post.commentsCount}</span>
              <Icon sprite={IconChatAlt2} className="h-5 w-5" />
            </HStack>
          )}
        </HStack>
        <Markdown className="c-posts-container__postdescription" maxLength={300} text={props.post.description} style="plainText" />
        {props.tags.length >= 1 && (
          <HStack spacing={0} className="gap-2 flex-wrap">
            {props.tags.map((tag) => (
              <ShowTag key={tag.id} tag={tag} />
            ))}
          </HStack>
        )}
        <div className="c-posts-container__post-votes">
          <span className="text-semibold text-2xl">{props.post.votesCount}</span>{" "}
          <span className="text-gray-700">{props.post.votesCount === 1 ? <Trans id="label.vote">Vote</Trans> : <Trans id="label.votes">Votes</Trans>}</span>
        </div>
      </VStack>
    </a>
  )
}

const MinimalListPostItem = (props: { post: Post; tags: Tag[] }) => {
  return (
    <HStack spacing={4} align="start" className="c-posts-container__post-minimal">
      <HStack className="w-full" justify="between" align="start">
        <a className="text-link" href={`/posts/${props.post.number}/${props.post.slug}`}>
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
            <MinimalListPostItem key={post.id} post={post} tags={props.tags.filter((tag) => post.tags.indexOf(tag.slug) >= 0)} />
          ))}
        </VStack>
      ) : (
        <>
          {props.posts.map((post) => (
            <ListPostItem key={post.id} post={post} tags={props.tags.filter((tag) => post.tags.indexOf(tag.slug) >= 0)} />
          ))}
        </>
      )}
    </>
  )
}
