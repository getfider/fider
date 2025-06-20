import "./Home.page.scss"
import NoDataIllustration from "@fider/assets/images/undraw-no-data.svg"

import React, { useState } from "react"
import { Post, Tag, PostStatus, ImageUpload } from "@fider/models"
import { Markdown, Hint, PoweredByFider, Icon, Header, Button } from "@fider/components"
import { PostsContainer } from "./components/PostsContainer"
import { useFider } from "@fider/hooks"
import { VStack } from "@fider/components/layout"
import { ShareFeedback } from "./components/ShareFeedback"
import { cache } from "@fider/services"
import { i18n } from "@lingui/core"
import { Trans } from "@lingui/react/macro"
import { CACHE_TITLE_KEY, CACHE_DESCRIPTION_KEY, CACHE_ATTACHMENT_KEY } from "./components/ShareFeedback"

export interface HomePageProps {
  posts: Post[]
  tags: Tag[]
  searchNoiseWords: string[]
  countPerStatus: { [key: string]: number }
  draftPost?: {
    id: number
    code: string
    title: string
    description: string
  }
  draftAttachments?: string[]
}

export interface HomePageState {
  title: string
}

const Lonely = () => {
  const fider = useFider()

  return (
    <div className="text-center">
      <Hint permanentCloseKey="at-least-3-posts" condition={fider.session.isAuthenticated && fider.session.user.isAdministrator}>
        <p>
          <Trans id="home.lonely.suggestion">
            It&apos;s recommended that you create <strong>at least 3</strong> suggestions here before sharing this site. The initial content is important to
            start engaging your audience.
          </Trans>
        </p>
      </Hint>
      <Icon sprite={NoDataIllustration} height="120" className="mt-6 mb-2" />
      <p className="text-muted">
        <Trans id="home.lonely.text">No posts have been created yet.</Trans>
      </p>
    </div>
  )
}

const HomePage = (props: HomePageProps) => {
  if (props.draftPost) {
    // Need to store the details of the draft post in the cache
    cache.session.set(CACHE_TITLE_KEY, props.draftPost.title)
    cache.session.set(CACHE_DESCRIPTION_KEY, props.draftPost.description)
    if (props.draftAttachments?.length) {
      const images: ImageUpload[] = props.draftAttachments.map((bkey: string) => ({ bkey, remove: false }))
      cache.session.set(CACHE_ATTACHMENT_KEY, JSON.stringify(images))
    }
  }

  const fider = useFider()
  // const [title, setTitle] = useState("")
  const [isShareFeedbackOpen, setIsShareFeedbackOpen] = useState(props.draftPost !== undefined)

  const defaultWelcomeMessage = i18n._("home.form.defaultwelcomemessage", {
    message: `We'd love to hear what you're thinking about.

What can we do better? This is the place for you to vote, discuss and share ideas.`,
  })

  const defaultInvitation = i18n._("home.form.defaultinvitation", {
    message: "Enter your suggestion here...",
  })

  const isLonely = () => {
    const len = Object.keys(props.countPerStatus).length
    if (len === 0) {
      return true
    }

    if (len === 1 && PostStatus.Deleted.value in props.countPerStatus) {
      return true
    }

    return false
  }

  const handleNewPost = () => {
    setIsShareFeedbackOpen(true)
  }

  return (
    <>
      <ShareFeedback
        tags={props.tags}
        placeholder={fider.session.tenant.invitation || defaultInvitation}
        isOpen={isShareFeedbackOpen}
        onClose={() => setIsShareFeedbackOpen(false)}
      />
      <Header />
      <div id="p-home" className="page container">
        <div className="p-home__welcome-col">
          <VStack spacing={2} className="p-4">
            <Markdown text={fider.session.tenant.welcomeMessage || defaultWelcomeMessage} style="full" />
            <Button className="c-input" type="submit" variant="secondary" onClick={handleNewPost}>
              {fider.session.tenant.invitation || defaultInvitation}
            </Button>
          </VStack>
          <div onClick={() => setIsShareFeedbackOpen(true)}>
            <PoweredByFider slot="home-input" className="sm:hidden md:hidden lg:block mt-3" />
          </div>
        </div>
        <div className="p-home__posts-col p-4">
          {isLonely() && <Lonely />}
          <PostsContainer posts={props.posts} tags={props.tags} countPerStatus={props.countPerStatus} />
          <PoweredByFider slot="home-footer" className="lg:hidden xl:hidden mt-8" />
        </div>
      </div>
    </>
  )
}

export default HomePage
