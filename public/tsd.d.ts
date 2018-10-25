import { CurrentUser, SystemSettings, Tenant } from "@fider/models";

declare global {
  interface Window {
    ga?: (cmd: string, evt: string, args?: any) => void;
    set: (key: string, value: any) => void;
    __props: { [key: string]: any };
    __user: CurrentUser | undefined;
    __tenant: Tenant;
    __settings: SystemSettings;
  }

  var __webpack_nonce__:string;
  var __webpack_public_path__:string; 
}

declare var require: (id: string) => any;
