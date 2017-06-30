import * as React from 'react';
import { Failure } from '@fider/services';
import { DisplayError } from './';

interface FormProps {
    onSubmit: () => Promise<any>;
}

interface FormState {
    failure?: Failure;
}

export class Form extends React.Component<FormProps, FormState> {

    public constructor(props: FormProps) {
        super(props);
    }

    public async submit() {
        if (this.props.onSubmit) {
            await this.props.onSubmit();
        }
    }

    public async clearFailure() {
        this.setState({ failure: undefined });
    }

    public async setFailure(failure: Failure) {
        this.setState({ failure });
    }

    public render() {
        return  <div className="ui form">
                    <DisplayError error={this.state && this.state.failure} />
                    { this.props.children }
                </div>;
    }

}
