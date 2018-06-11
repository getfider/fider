import { FiderPage, CurrentUser, SystemSettings, Tenant } from "@fider/models";

declare global {
  interface Window {
    ga?: (cmd: string, evt: string, args?: any) => void;
    set: (key: string, value: any) => void;
  }

  var page: FiderPage;
}

declare var require: (id: string) => any;
