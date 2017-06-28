import * as React from 'react';

interface ButtonProps {
    classes?: string;
    href?: string;
    onClick?: () => Promise<any>;
}

interface ButtonState {
    clicked: boolean;
}

export class Button extends React.Component<ButtonProps, ButtonState> {

    public constructor(props: ButtonProps) {
        super(props);
        this.state = {
            clicked: false
        };
    }

    private async click() {
        this.setState({ clicked: true });
        if (this.props.onClick) {
            await this.props.onClick();
            this.setState({ clicked: false });
        }
    }

    public render() {
        const cssClasses = `ui button ${this.props.classes || ''} ${this.state.clicked ? 'loading disabled' : ''}`;
        if (this.props.href) {
            return <a href={this.props.href} className={cssClasses} onClick={() => this.click()}>
                        { this.props.children }
                   </a>;
        } else {
            return <button className={ cssClasses } onClick={() => this.click()}>
                        { this.props.children }
                   </button>;
        }
    }

}
