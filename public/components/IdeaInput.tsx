import axios from "axios";
import * as React from "react";
import { SocialSignInButton } from "./SocialSignInButton";

const claims = (window as any)._claims;

interface IdeaInputState {
    idea: string;
    clicked: boolean;
}

export class IdeaInput extends React.Component<{}, IdeaInputState> {
    private description: HTMLTextAreaElement;

    constructor() {
        super();
        this.state = {
          idea: "",
          clicked: false
        };
    }

    public async submit() {
      this.setState({ clicked: true });

      const response = await axios.post("/api/ideas", {
        title: this.state.idea,
        description: this.description.value
      });

      location.reload();
    }

    public render() {
        const buttonClasses = `ui positive button ${this.state.clicked && "loading disabled"}`;
        const details = claims ?
                        <div>
                          <div className="field">
                            <textarea ref={(ref) => this.description = ref }
                                      rows={6}
                                      placeholder="Describe your idea (optional)"></textarea>
                            </div>
                              <button className={buttonClasses} onClick={async () => { await this.submit(); } }>
                                Submit Idea
                              </button>
                            </div> :
                          <div>
                          <div className="ui message">
                            <div className="header">
                              Please log in before posting an idea
                            </div>
                            <p>This only takes a second and you'll be good to go!</p>
                            <SocialSignInButton provider="facebook" small={true} />
                            <SocialSignInButton provider="google" small={true} />
                          </div>
                        </div>;

        return <div className="ui form">
                <div className="ui fluid input">
                    <input  id="new-idea-input"
                            type="text"
                            maxLength={60}
                            onKeyUp={(e) => { this.setState({ idea: e.currentTarget.value }); }}
                            placeholder="Enter your idea, new feature or suggestion here ..." />
                </div>

                { this.state.idea.length > 0 && details }
               </div>;
    }
}
