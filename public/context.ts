import * as React from "react";
import { CurrentUser, Tenant, SystemSettings } from "@fider/models";

interface PageContext {
  tenant: Tenant;
  user?: CurrentUser;
  system: SystemSettings;
}

export const PageContext = React.createContext<PageContext>({} as any);
