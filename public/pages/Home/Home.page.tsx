import "./Home.page.scss";

import React, { useState } from "react";
import { Post, Tag, PostStatus } from "@fider/models";
import { MultiLineText, Hint } from "@fider/components";
import { SimilarPosts } from "./components/SimilarPosts";
import { FaRegLightbulb } from "react-icons/fa";
import { PostInput } from "./components/PostInput";
import { PostsContainer } from "./components/PostsContainer";
import { useFider } from "@fider/hooks";
import { useTranslation, WithTranslation } from "react-i18next";

export interface HomePageProps extends WithTranslation {
  posts: Post[];
  tags: Tag[];
  countPerStatus: { [key: string]: number };
}

export interface HomePageState {
  title: string;
}

const Lonely = () => {
  const fider = useFider();
  const { t } = useTranslation();

  return (
    <div className="l-lonely center">
      <Hint
        permanentCloseKey="at-least-3-posts"
        condition={fider.session.isAuthenticated && fider.session.user.isAdministrator}
      >
        {t("home.lonelyHint")}
      </Hint>
      <p>
        <FaRegLightbulb />
      </p>
      <p dangerouslySetInnerHTML={{ __html: t("home.lonelyAdmin") }} />
    </div>
  );
};

const HomePage = (props: HomePageProps) => {
  const fider = useFider();
  const { t } = useTranslation();

  const [title, setTitle] = useState("");

  const isLonely = () => {
    const len = Object.keys(props.countPerStatus).length;
    if (len === 0) {
      return true;
    }

    if (len === 1 && PostStatus.Deleted.value in props.countPerStatus) {
      return true;
    }

    return false;
  };

  return (
    <div id="p-home" className="page container">
      <div className="row">
        <div className="l-welcome-col col-md-4">
          <MultiLineText
            className="welcome-message"
            text={fider.session.tenant.welcomeMessage || t("home.defaultWelcome")}
            style="full"
          />
          <PostInput
            placeholder={fider.session.tenant.invitation || t("home.suggestionPlaceholder")}
            onTitleChanged={setTitle}
          />
        </div>
        <div className="l-posts-col col-md-8">
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
  );
};

export default HomePage;
