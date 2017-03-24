import * as moment from "moment";
import * as React from "react";
import { Idea } from "../models";
import * as storage from "../storage";
import { Gravatar, MultiLineText } from "./Common";

import { Footer } from "./Footer";
import { Header } from "./Header";
import { IdeaInput } from "./IdeaInput";

export class ShowIdeaRoot extends React.Component<{}, {}> {
    public render() {
        const idea = storage.get<Idea>("idea");

        return <div>
                  <Header />
                  <div className="ui container">
                    <h1 className="ui header">{ idea.title }</h1>

                    <p>{ idea.description }</p>

                    <p>
                      <Gravatar email={idea.user.email}/> <u>{idea.user.name}</u>
                      &nbsp;shared <u title={idea.createdOn}>{ moment(idea.createdOn).fromNow() }</u>
                    </p>

                    <p><a href="/">Back</a></p>
                  </div>
                  <Footer />
               </div>;
    }
}
