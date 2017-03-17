import axios from "axios";
import * as React from "react";
import { SocialSignInButton } from "./SocialSignInButton";

const claims = (window as any)._claims;

interface IdeaInputState {
    idea: string;
}

export class IdeaInput extends React.Component<{}, IdeaInputState> {

    constructor() {
        super();
        this.state = { idea: "" };
    }

    public async submit() {
      const response = await axios.post("/api/ideas", {
        title: this.state.idea
      });
    }

    public render() {
        const details = claims ?
                        <button className="ui positive button" onClick={async () => { await this.submit(); } }>
                          Submit Idea
                        </button> :
                        <div>
                          <p>Please log in before posting an idea...</p>
                          <div className="ui list">
                            <div className="item">
                              <SocialSignInButton provider="facebook"/>
                            </div>
                            <div className="item">
                              <SocialSignInButton provider="google"/>
                            </div>
                          </div>
                        </div>;

        return <div>
                <div className="ui fluid input">
                    <input  id="new-idea-input"
                            type="text"
                            onKeyUp={(e) => { this.setState({ idea: e.currentTarget.value }); }}
                            placeholder="Enter your idea, new feature or suggestion here ..." />
                </div>

                <div id="new-idea-submit" className="ui grid">
                    <div className="four wide column">
                    { this.state.idea.length > 0 && details }
                    </div>
                </div>
               </div>;
    }
}
