import * as React from 'react';

import { Header, Footer } from '@fider/components/common';

export class AdminHomePage extends React.Component<{}, {}> {

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
