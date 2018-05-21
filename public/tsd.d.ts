import { CurrentUser, SystemSettings, Tenant } from "@fider/models";

interface PageProps {
  user?: CurrentUser;
  system: SystemSettings;
  tenant: Tenant;
  [key: string]: any;
}

declare global {
  interface Window {
    ga?: (cmd: string, evt: string, args?: any) => void;
    props: PageProps;
    set: (key: string, value: any) => void;
  }
}

declare var require: (id: string) => any;
