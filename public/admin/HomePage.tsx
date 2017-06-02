import * as React from 'react';

import { Footer } from '../shared/Footer';
import { Header } from '../shared/Header';

export class HomePage extends React.Component<{}, {}> {

    constructor(props: {}) {
        super(props);
    }

    public render() {
      return <div>
                <Header />
                <div className="ui container">
                <h1 className="ui header">Administration</h1>
                <span>You're a Member or Admin!</span>
                </div>
                <Footer />
            </div>;
    }
}
