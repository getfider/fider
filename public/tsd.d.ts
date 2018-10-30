import { CurrentUser, SystemSettings, Tenant } from "@fider/models";

declare global {
  interface Window {
    ga?: (cmd: string, evt: string, args?: any) => void;
    set: (key: string, value: any) => void;
    __props: { [key: string]: any };
    __contextID: string;
    __user: CurrentUser | undefined;
    __tenant: Tenant;
    __settings: SystemSettings;
  }
}

declare var require: (id: string) => any;
