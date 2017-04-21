import axios from "axios";
import * as React from "react";
import { Idea } from "../models";

interface SupportCounterProps {
    idea: Idea;
}

interface SupportCounterState {
    supported: boolean;
    total: number;
}

export class SupportCounter extends React.Component<SupportCounterProps, SupportCounterState> {
    constructor(props: SupportCounterProps) {
        super(props);
        this.state = {
          supported: props.idea.totalSupporters >= 1 && props.idea.totalSupporters <= 10,
          total: props.idea.totalSupporters
        };
    }

    public async undo() {
        if (!this.state.supported) { return; }

        try {
            await axios.post(`/api/ideas/${this.props.idea.number}/unsupport`);

            this.setState({
                supported: false,
                total: this.state.total - 1
            });
        } catch (ex) {
            // TODO: handle this
        }

    }

    public async support() {
        if (this.state.supported) { return; }

        try {
            await axios.post(`/api/ideas/${this.props.idea.number}/support`);

            this.setState({
                supported: true,
                total: this.state.total + 1
            });
        } catch (ex) {
            // TODO: handle this
        }

    }

    public render() {
        const support = <div className="support-button ui mini violet inverted animated button"
                    onClick={async () => await this.support()}>
                    <div className="visible content">Want</div>
                    <div className="hidden content">
                        <i className="heart icon"></i>
                    </div>
                </div>;
        const undo = <div className="support-button ui mini violet animated button"
                    onClick={async () => await this.undo()}>
                    <div className="visible content"><i className="heart icon"></i></div>
                    <div className="hidden content">
                        Undo
                    </div>
                </div>;

        return <div className="ui small statistics">
                    <div className="statistic">
                        <div className="value">
                            { this.state.total }
                        </div>
                        { this.state.supported ? undo : support }
                    </div>
                </div>;
    }
}
