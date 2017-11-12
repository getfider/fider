import * as React from 'react';

interface ButtonProps {
    className?: string;
    href?: string;
    size?: 'mini' | 'tiny' | 'small' | 'large';
    onClick?: (event: ButtonClickEvent) => Promise<any>;
}

interface ButtonState {
    clicked: boolean;
}

export class ButtonClickEvent {
    private shouldEnable = true;
    public preventEnable(): void {
        this.shouldEnable = false;
    }
    public canEnable(): boolean {
        return this.shouldEnable;
    }
}

export class Button extends React.Component<ButtonProps, ButtonState> {
    private unmounted: boolean;

    public static defaultProps: Partial<ButtonProps> = {
        size: 'tiny'
    };

    public constructor(props: ButtonProps) {
        super(props);
        this.state = {
            clicked: false
        };
    }

    public componentWillUnmount() {
        this.unmounted = true;
    }

    private async click(e?: React.MouseEvent<HTMLButtonElement>) {
        if (e) {
            e.preventDefault();
            e.stopPropagation();
        }

        const event = new ButtonClickEvent();
        this.setState({ clicked: true });
        if (this.props.onClick) {
            await this.props.onClick(event);
            if (!this.unmounted && event.canEnable()) {
                this.setState({ clicked: false });
            }
        }
    }

    public render() {
        const cssClasses = `ui ${this.props.size} button ${this.props.className || ''} ${this.state.clicked ? 'loading disabled' : ''}`;
        if (this.props.href) {
            return <a href={ this.props.href } className={ cssClasses } onClick={ () => this.click() }>
                        { this.props.children }
                   </a>;
        } else {
            return <button type="button" className={ cssClasses } onClick={ (e) => this.click(e) }>
                        { this.props.children }
                   </button>;
        }
    }

}
