import * as React from 'react';
import { get } from '../storage';
import { SocialSignInButton } from './SocialSignInButton';

interface AuthSettings {
    endpoint: string;
    providers: {
        google: boolean,
        facebook: boolean
    };
}

interface SocialSignInListProps {
    size: 'small' | 'normal';
    orientation: 'horizontal' | 'vertical';
}

export const SocialSignInList = (props: SocialSignInListProps) => {
    const settings = get<AuthSettings>('auth');
    const cssClasses = props.orientation === 'horizontal' ? 'horizontal divided' : '';

    const google = settings.providers.google &&
                    <div className="item">
                        <SocialSignInButton provider="google" size={props.size} />
                    </div>;
    const facebook = settings.providers.facebook &&
                    <div className="item">
                        <SocialSignInButton provider="facebook" size={props.size} />
                    </div>;

    const noAuth = !facebook && !google &&
                    <div className="item">
                        <div className="ui tertiary inverted red segment">
                            There are no authentication methods enabled.
                        </div>
                    </div>;

    return <div className={`ui list ${cssClasses}`}>
                { facebook }
                { google }
                { noAuth }
            </div>;
};
