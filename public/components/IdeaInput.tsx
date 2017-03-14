import * as React from "react";

export class IdeaInput extends React.Component<{}, {}> {
    render() {
        return <div className="ui fluid input">
                    <input id="new-idea-input" type="text" placeholder="Enter your idea, new feature or suggestion here ..." />
               </div>;
    }
}