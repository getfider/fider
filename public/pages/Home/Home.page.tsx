/* eslint-disable prettier/prettier */
import "./Home.page.scss"
import NoDataIllustration from "@fider/assets/images/undraw-no-data.svg"
import HeroIllustration from "@fider/assets/images/Illustration.svg"

import React, { useState } from "react"
import { Post, Tag, PostStatus } from "@fider/models"
import { Markdown, Hint, PoweredByFider, Icon } from "@fider/components"
import { SimilarPosts } from "./components/SimilarPosts"
import { PostInput } from "./components/PostInput"
import { PostsContainer } from "./components/PostsContainer"
import { useFider } from "@fider/hooks"
import { VStack } from "@fider/components/layout"

export interface HomePageProps {
  posts: Post[]
  tags: Tag[]
  countPerStatus: { [key: string]: number }
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
          It&apos;s recommended that you create <strong>at least 3</strong> suggestions here before sharing this site. The initial content is important to start
          engaging your audience.
        </p>
      </Hint>
      <Icon sprite={NoDataIllustration} height="120" className="mt-6 mb-2" />
      <p className="text-muted">No posts have been created yet.</p>
    </div>
  )
}

const defaultWelcomeMessage = `We'd love to hear what you're thinking about. 

What can we do better? This is the place for you to vote, discuss and share ideas.`

const HomePage = (props: HomePageProps) => {
  const fider = useFider()
  const [title, setTitle] = useState("")

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

  return (
    <div id="p-home" className="page">

      <section className="hero">
        <div className="container">
          <div className="row">
            <div className="col col-10 col-md-6 col-caption">
              <MultiLineText className="welcome-message" text={fider.session.tenant.welcomeMessage || defaultWelcomeMessage} style="full" />
            </div>
            <div className="col col-2 col-md-6 col-image">
              <img src={HeroIllustration} alt="" />
            </div>
          </div>
        </div>
        <svg width="1440" height="302" viewBox="0 0 1440 302" fill="none" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none">
          <path d="M1.50678e-05 0H1440C1440 0 1440 134.41 1440 274.026C733.519 200.865 707.983 364.255 1.50678e-05 274.026C-3.09845e-05 128.913 1.50678e-05 0 1.50678e-05 0Z" fill="#6291EB" />
        </svg>
      </section>

      <section className="suggestions">
        <div className="container">
          <div className="row">
            <div className="col col-md-6">
              <PostInput placeholder={fider.session.tenant.invitation || "Enter your suggestion here..."} onTitleChanged={setTitle} />
            </div>
          </div>
        </div>
      </section>

      <section className="posts">
        <div className="container">
          <div className="row">
            <div className="col col-posts">
              {isLonely() ? (
                <Lonely />
              ) : title ? (
                <SimilarPosts title={title} tags={props.tags} />
              ) : (
                <PostsContainer posts={props.posts} tags={props.tags} countPerStatus={props.countPerStatus} />
              )}
            </div>
          </div>
        </div>
      </section>

      <section className="footer">
        <div className="container">
          <div className="row">
            <div className="col">
              <a rel="noopener" href="https://members.bcc.no/privacy-statement/" target="_blank">Privacy statement</a>
              <PoweredByFider />
            </div>
          </div>
        </div>
      </section>

    </div>
  )
}

export default HomePage



