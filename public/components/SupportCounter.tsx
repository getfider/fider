import * as React from "react";
import { Idea } from "../models";

interface SupportCounterProps {
    idea: Idea;
}

interface SupportCounterState {
    supported: boolean;
}

export class SupportCounter extends React.Component<SupportCounterProps, SupportCounterState> {
    constructor(props: SupportCounterProps) {
        super(props);
        this.state = {
          supported: props.idea.totalSupporters >= 1 && props.idea.totalSupporters <= 10,
        };
    }

    public undoButton() {
        return <div className="support-button ui mini violet animated button">
                    <div className="visible content"><i className="heart icon"></i></div>
                    <div className="hidden content">
                        Undo
                    </div>
                </div>;
    }

    public supportButton() {
        return <div className="support-button ui mini violet inverted animated button">
                    <div className="visible content">Want</div>
                    <div className="hidden content">
                        <i className="heart icon"></i>
                    </div>
                </div>;
    }

    public render() {
        const button = this.state.supported ? this.undoButton() : this.supportButton();

        return <div className="ui small statistics">
                    <div className="statistic">
                        <div className="value">
                            { this.props.idea.totalSupporters }
                        </div>
                        { button }
                    </div>
                </div>;
    }
}
