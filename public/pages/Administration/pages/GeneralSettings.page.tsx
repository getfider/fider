import "./GeneralSettings.page.scss";

import React from "react";

import {
  Button,
  ButtonClickEvent,
  TextArea,
  Form,
  Input,
  ImageUploader,
  ImageUploadState,
  TenantLogoURL
} from "@fider/components/common";
import { actions, Failure, fileToBase64, Fider } from "@fider/services";
import { FaCogs } from "react-icons/fa";
import { AdminBasePage } from "../components/AdminBasePage";

interface GeneralSettingsPageProps {
  publicIP: string;
}

interface GeneralSettingsPageState {
  logo?: ImageUploadState;
  title: string;
  invitation: string;
  welcomeMessage: string;
  cname: string;
  error?: Failure;
}

export default class GeneralSettingsPage extends AdminBasePage<GeneralSettingsPageProps, GeneralSettingsPageState> {
  private fileSelector?: HTMLInputElement | null;

  public id = "p-admin-general";
  public name = "general";
  public icon = FaCogs;
  public title = "General";
  public subtitle = "Manage your site settings";

  constructor(props: GeneralSettingsPageProps) {
    super(props);

    this.state = {
      title: Fider.session.tenant.name,
      cname: Fider.session.tenant.cname,
      welcomeMessage: Fider.session.tenant.welcomeMessage,
      invitation: Fider.session.tenant.invitation
    };
  }

  private handleSave = async (e: ButtonClickEvent) => {
    const result = await actions.updateTenantSettings(this.state);
    if (result.ok) {
      e.preventEnable();
      location.href = `/`;
    } else if (result.error) {
      this.setState({ error: result.error });
    }
  };

  public fileChanged = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      const base64 = await fileToBase64(file);
      this.setState({
        logo: {
          upload: {
            content: base64,
            contentType: file.type,
            action: "upload"
          },
          ignore: false,
          remove: false
        }
      });
    }
  };

  public removeFile = async (e: ButtonClickEvent) => {
    this.setState({
      logo: {
        ignore: false,
        remove: true
      }
    });
  };

  public selectFile = async (e: ButtonClickEvent) => {
    if (this.fileSelector) {
      this.fileSelector.click();
    }
  };

  public dnsInstructions(): JSX.Element {
    const isApex = this.state.cname.split(".").length === 2;
    const recordType = isApex ? "A" : "CNAME";
    const publicIP = this.props.publicIP || "<error>";
    const targetRecord = isApex ? publicIP : `${Fider.session.tenant.subdomain}${Fider.settings.domain}`;
    return (
      <>
        <strong>{this.state.cname}</strong> {recordType} <strong>{targetRecord}</strong>
      </>
    );
  }

  private setTitle = (title: string): void => {
    this.setState({ title });
  };

  private setWelcomeMessage = (welcomeMessage: string): void => {
    this.setState({ welcomeMessage });
  };

  private setInvitation = (invitation: string): void => {
    this.setState({ invitation });
  };

  private setLogo = (logo: ImageUploadState): void => {
    this.setState({ logo });
  };

  private setCNAME = (cname: string): void => {
    this.setState({ cname });
  };

  public content() {
    return (
      <Form error={this.state.error}>
        <Input
          field="title"
          label="Title"
          maxLength={60}
          value={this.state.title}
          disabled={!Fider.session.user.isAdministrator}
          onChange={this.setTitle}
        />
        <TextArea
          field="welcomeMessage"
          label="Welcome Message"
          value={this.state.welcomeMessage}
          disabled={!Fider.session.user.isAdministrator}
          onChange={this.setWelcomeMessage}
        />
        <Input
          field="invitation"
          label="Invitation"
          maxLength={60}
          value={this.state.invitation}
          disabled={!Fider.session.user.isAdministrator}
          onChange={this.setInvitation}
        />

        <ImageUploader
          label="Logo"
          field="logo"
          defaultImageURL={TenantLogoURL(200)}
          previewMaxWidth={200}
          disabled={!Fider.session.user.isAdministrator}
          onChange={this.setLogo}
        >
          <p className="info">
            We accept JPG, GIF and PNG images, smaller than 100KB and with an aspect ratio of 1:1 with minimum
            dimensions of 200x200 pixels.
          </p>
        </ImageUploader>

        {!Fider.isSingleHostMode() && (
          <Input
            field="cname"
            label="Custom Domain"
            maxLength={100}
            placeholder="feedback.yourcompany.com"
            value={this.state.cname}
            disabled={!Fider.session.user.isAdministrator}
            onChange={this.setCNAME}
          >
            <div className="info">
              {this.state.cname ? (
                [
                  <p key={0}>Enter the following record into your DNS zone records:</p>,
                  <p key={1}>{this.dnsInstructions()}</p>,
                  <p key={2}>
                    Please note that it may take up to 72 hours for the change to take effect worldwide due to DNS
                    propagation.
                  </p>
                ]
              ) : (
                <p>
                  Custom domains allow you to access your app via your own domain name (for example,{" "}
                  <code>feedback.yourcompany.com</code>
                  ).
                </p>
              )}
            </div>
          </Input>
        )}

        <div className="field">
          <Button disabled={!Fider.session.user.isAdministrator} color="positive" onClick={this.handleSave}>
            Save
          </Button>
        </div>
      </Form>
    );
  }
}
