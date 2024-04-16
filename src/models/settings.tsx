
interface JSONMap {
    [key: string]: any;
  }

export interface Connector {
    connector_specific_config: JSONMap;
    created_date: string;
    credential_id: number;
    deleted_date?: string | null;
    disabled: boolean;
    id: number;
    input_type: string;
    last_attempt_status: string;
    last_successful_index_time?: string | null;
    name: string;
    refresh_freq: number;
    shared: boolean;
    source: string;
    tenant_id: string;
    total_docs_indexed: number;
    updated_date?: string | null;
    user_id: string;
  }
  
  export interface Credential {
    created_date: string;
    credential_json: JSONMap;
    deleted_date?: string | null;
    id: number;
    shared: boolean;
    source: string;
    tenant_id: string;
    updated_date?: string | null;
    user_id: string;
  }
  
  export interface EmbeddingModel {
    created_date: string;
    deleted_date?: string | null;
    id: number;
    index_name: string;
    is_active: boolean;
    model_dim: number;
    model_id: string;
    model_name: string;
    normalize: boolean;
    passage_prefix: string;
    query_prefix: string;
    tenant_id: string;
    updated_date?: string | null;
    url: string;
  }
  
  
  export interface LLM {
    endpoint: string;
    id: number;
    model_id: string;
    name: string;
    url: string;
  }
  
  export interface Persona {
    default_persona: boolean;
    description: string;
    display_priority: number;
    id: number;
    is_visible: boolean;
    llm: LLM;
    llm_id: number;
    name: string;
    starter_messages: number[];
    tenant_id: string;
  }