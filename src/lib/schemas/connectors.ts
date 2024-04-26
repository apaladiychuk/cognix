export type CreateConnectorSchema = {
    connector_id: number | null,
    source: string,
    connector_specific_config: object,
    refresh_freq: string,
    credential_id: number | null,
}

