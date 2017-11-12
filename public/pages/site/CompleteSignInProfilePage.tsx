import * as React from 'react';

import { SiteHomePage } from './SiteHomePage';
import { showModal, getQueryString } from '@fider/utils/page';
import { Form } from '@fider/components/common';

import { inject, injectables } from '@fider/di';
import { TenantService } from '@fider/services';

interface CompleteSignInProfilePageState {
  name: string;
}

export class CompleteSignInProfilePage extends React.Component<{}, CompleteSignInProfilePageState> {
  private form: Form;
  private key: string;

  @inject(injectables.TenantService)
  public service: TenantService;

  constructor(props: {}) {
    super(props);
    this.key = getQueryString('k');
    this.state = {
      name: '',
    };
  }

  public componentDidMount() {
    showModal('#signin-complete-modal', { closable: false });
  }

  private async submit() {
    const result = await this.service.completeProfile(this.key, this.state.name);
    if (result.ok) {
      location.href = '/';
    } else if (result.error) {
      this.form.setFailure(result.error);
    }
  }

  public render() {
    const modal = (
      <div id="signin-complete-modal" className="ui modal small">
        <div className="header">
            Complete your profile
        </div>
        <div className="content">
          <p>Because this is your first sign in, please input your display name.</p>
          <Form ref={(f) => { this.form = f!; }} onSubmit={() => this.submit()}>
            <div className="ui small action fluid input">
              <input
                onChange={(e) => this.setState({ name: e.currentTarget.value })}
                type="text"
                maxLength={100}
                placeholder="Your display name"
                className="small"
              />
              <button onClick={() => this.form.submit()} className={`ui small positive button ${this.state.name === '' && 'disabled'}`}>
                Submit
              </button>
            </div>
          </Form>
        </div>
      </div>
    );

    return (
      <div>
        {modal}
        <SiteHomePage />
      </div>
    );
  }
}
