import { CurrentUser, SystemSettings, Tenant } from "@fider/models";
import { Fider } from "./fider";

declare global {
  interface Window {
    ga?: (cmd: string, evt: string, args?: any) => void;
    set: (key: string, value: any) => void;
    __props: { [key: string]: any };
    __user: CurrentUser | undefined;
    __tenant: Tenant;
    __settings: SystemSettings;
  }
}

declare var require: (id: string) => any;
