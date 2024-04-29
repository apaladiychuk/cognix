export type CreateConnectorSchema = {
  connector_id?: number;
  source: string;
  connector_specific_config: object;
  refresh_freq: string;
  credential_id?: string;
};
