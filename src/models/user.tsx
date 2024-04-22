import { JSONMap } from "./settings"

export interface Tenant {
    id: string;
    name: string;
    configuration: JSONMap;
}


export interface Credential {
    id: string;
    user_name: string;
    first_name: string;
    last_name: string;
    roles: string[];
    tenant: Tenant;
    tenant_id: string;
  }