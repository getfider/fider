import * as React from 'react';

interface ButtonProps {
    className?: string;
    href?: string;
    size?: 'tiny' | 'large';
    onClick?: () => Promise<any>;
}

interface ButtonState {
    clicked: boolean;
}

export class Button extends React.Component<ButtonProps, ButtonState> {
    public static defaultProps: Partial<ButtonProps> = {
        size: 'tiny'
    };

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
        const cssClasses = `ui ${this.props.size} button ${this.props.className || ''} ${this.state.clicked ? 'loading disabled' : ''}`;
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
