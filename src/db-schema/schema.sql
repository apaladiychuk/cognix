create table tenants
(
    id    uuid  not null    primary key,
    name  varchar(255),
    configuration jsonb default '{}'::JSONB not null
);
create table users
(
    id  uuid  not null    primary key,
    tenant_id   uuid  not null    references tenants,
    user_name   varchar(255)  not null    unique,
    first_name  varchar(255),
    last_name   varchar(255),
    external_id text,
    roles   text[] default '{}'::STRING[] not null
);

create table llm
(
    id   bigint    default unique_rowid() not null    primary key,
    tenant_id    uuid not null    references tenants,
    name varchar(255) not null,
    model_id varchar(255) not null,
    url  varchar  not null,
    api_key  varchar,
    endpoint varchar,
    created_date timestamp default now()  not null,
    updated_date timestamp,
    deleted_date timestamp
);

create table embedding_models
(
    id   bigint    default unique_rowid() not null    primary key,
    tenant_id    uuid not null,
    model_id varchar  not null,
    model_name   varchar  not null,
    model_dim    bigint   not null,
    url  varchar,
    is_active    boolean   default false  not null,
    created_date timestamp default now()  not null,
    updated_date timestamp,
    deleted_date timestamp
);

create table personas
(
    id   bigint    default unique_rowid() not null    primary key,
    name varchar  not null,
    llm_id   bigint    references llm,
    default_persona  boolean  not null,
    description  varchar  not null,
    tenant_id    uuid not null    references tenants,
    is_visible   boolean  not null,
    display_priority bigint,
    starter_messages jsonb default '{}'::JSONB    not null,
    created_date timestamp default now()  not null,
    updated_date timestamp,
    deleted_date timestamp
);

create table prompts
(
    id    bigint    default unique_rowid() not null    primary key,
    persona_id    bigint   not null    references personas,
    user_id   uuid not null    references users,
    name  varchar  not null,
    description   varchar  not null,
    system_prompt text not null,
    task_prompt   text not null,
    include_citations boolean  not null,
    datetime_aware    boolean  not null,
    default_prompt    boolean  not null,
    created_date  timestamp default now()  not null,
    updated_date  timestamp,
    deleted_date  timestamp
);

create table credentials
(
    id  bigint    default unique_rowid() not null    primary key,
    credential_json jsonb default '{}'::JSONB    not null,
    user_id uuid not null    references users,
    tenant_id   uuid not null    references tenants,
    source  varchar(50)  not null,
    created_date    timestamp default now()  not null,
    updated_date    timestamp,
    deleted_date    timestamp,
    shared  boolean  not null
);

create table connectors
(
    id bigint    default unique_rowid() not null    primary key,
    credential_id  bigint    references credentials,
    name   varchar  not null,
    source varchar(50)  not null,
    connector_specific_config  jsonb    not null,
    refresh_freq   bigint,
    user_id    uuid not null    references users,
    tenant_id  uuid not null    references tenants,
    shared boolean  not null,
    disabled   boolean  not null,
    last_successful_index_time timestamp,
    last_attempt_status    varchar,
    total_docs_indexed bigint   not null,
    created_date   timestamp default now()  not null,
    updated_date   timestamp,
    deleted_date   timestamp
);

create table chat_sessions
(
    id   bigint    default unique_rowid() not null    primary key,
    user_id  uuid not null    references users,
    description  text not null,
    created_date timestamp default now()  not null,
    deleted_date timestamp,
    persona_id   bigint   not null
    references personas,
    one_shot boolean  not null
);

create table chat_messages
(
    id   bigint    default unique_rowid() not null    primary key,
    chat_session_id  bigint   not null    references chat_sessions,
    message  text not null,
    message_type varchar(9)   not null,
    time_sent    timestamp default now()  not null,
    token_count  bigint   not null,
    parent_message   bigint,
    latest_child_message bigint,
    rephrased_query  text,
    citations    jsonb default '{}'::JSONB    not null,
    error    text
);

create table chat_message_feedbacks
(
    id  bigint  default unique_rowid() not null    primary key,
    chat_message_id bigint not null    references chat_messages,
    user_id uuid   not null    references users,
    up_votes    boolean    not null,
    feedback    varchar default ''::STRING not null
);

create table documents
(
    id   bigint   default unique_rowid() not null    primary key,
    document_id  varchar not null,
    connector_id bigint  not null    references connectors,
    link varchar,
    signature    text,
    created_date timestamp    default now()  not null,
    updated_date timestamp,
    deleted_date timestamp,
    status   varchar(100) default 'new'::STRING  not null
);

create table document_feedbacks
(
    id    bigint  default unique_rowid() not null    primary key,
    document_id   bigint    references documents,
    user_id   uuid    references users,
    document_rank bigint not null,
    up_votes  boolean    not null,
    feedback  varchar default ''::STRING not null
);

create table document_sets
(
    id    bigint    default unique_rowid() not null    primary key,
    user_id   uuid    references users,
    name  varchar  not null,
    description   varchar  not null,
    is_up_to_date boolean  not null,
    created_date  timestamp default now()  not null,
    updated_date  timestamp,
    deleted_date  timestamp
);

create table document_set_connector_pairs
(
    document_set_id bigint  not null    references document_sets,
    connector_id    bigint  not null    references connectors,
    is_current  boolean not null,
    primary key (document_set_id, connector_id)
);