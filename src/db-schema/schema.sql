create table connectors
(
    id                        serial        primary key,
    name                      varchar                                not null,
    source                    varchar(50)                            not null,
    input_type                varchar(10),
    connector_specific_config jsonb                                  not null,
    refresh_freq              integer,
    user_id                   uuid not null,
    tenant_id                 uuid not null,
    created_date              timestamp without time zone default now() not null,
    deleted_date              timestamp without time zone null,
);

create table credentials
(
    id              serial   primary key,
    credential_json jsonb                                  not null,
    user_id         uuid not null,
    tenant_id       uuid not null,
    source          varchar(50)                            not null,
    created_date    timestamp without time zone default now() not null,
    updated_date    timestamp without time zone null,
    deleted_date    timestamp without time zone null,
    admin_public    boolean                                not null
);

create table connector_credential_pairs
(
    connector_id               integer not null references connectors on delete cascade,
    credential_id              integer not null references credentials on delete cascade,
    last_successful_index_time timestamp without time zone,
    last_attempt_status        varchar,
    total_docs_indexed         integer not null,
    name                       varchar                                                               not null,
    is_public                  boolean                                                               not null,
    primary key (connector_id, credential_id)
);

create table embedding_models
(
    id             serial        primary key,
    tenant_id      uuid             not null,
    model_id       varchar          not null,
    model_name     varchar          not null,
    model_dim      integer          not null,
    normalize      boolean          not null,
    query_prefix   varchar          not null,
    passage_prefix varchar          not null,
    index_name     varchar          not null,
    is_active      boolean  not null default false
);

create table index_attempts
(
    id                      serial primary key,
    created_at              timestamp without time zone default now() not null,
    status                  varchar not null,
    error_msg               varchar,
    connector_id            integer references connectors,
    credential_id           integer references credentials,
    total_docs_indexed      integer,
    time_started            timestamp without time zone,
    new_docs_indexed        integer,
    embedding_model_id      integer not null references embedding_models,
    from_beginning          boolean                                not null,
    full_exception_trace    text,
    docs_removed_from_index integer
);

create index ix_index_attempt_latest_for_connector_credential_pair
    on index_attempt (connector_id, credential_id, time_created);

create table llm (
                     id serial primary key ,
                     name varchar  not null,
                     model_id varchar not null,
                     url varchar   not null
);

create type searchtype as enum ('KEYWORD', 'SEMANTIC', 'HYBRID');

create table searchtypes (
    id integer not null primary key,
    key varchar not null unique,
    name varchar not null
);
insert into searchtypes (id, key, name) values(1, 'KEYWORD', 'Keyword');
insert into searchtypes (id, key, name) values(2, 'SEMANTIC', 'Semantic search');
insert into searchtypes (id, key, name) values(2, 'HYBRID', 'Hybrid search');

create type recencybiassetting as enum ('FAVOR_RECENT', 'BASE_DECAY', 'NO_DECAY', 'AUTO');

create table personas
(
    id                         serial       primary key,
    name                       varchar            not null,
    llm_id                     integer references llm(id) ,
    default_persona            boolean            not null,
    description                varchar            not null,
    tenant_id                  uuid not null,
    search_type                searchtype         not null,
    search_type_id             int references searchtypes,
    is_visible                 boolean            not null,
    display_priority           integer,
    starter_messages           jsonb
);

create table prompts
(
    id                serial       primary key,
    user_id           uuid    not null,
    name              varchar not null,
    description       varchar not null,
    system_prompt     text    not null,
    task_prompt       text    not null,
    include_citations boolean not null,
    datetime_aware    boolean not null,
    default_prompt    boolean not null,
    created_date      timestamp without time zone default now() not null,
    deleted_date      timestamp without time zone null
);


create table chat_sessions
(
    id           serial       primary key,
    user_id      uuid not null,
    description  text not null,
    created_date timestamp without time zone default now() not null,
    deleted_date timestamp without time zone null,
    persona_id   integer not null references persona,
    one_shot     boolean not null
);



create table chat_messages
(
    id                   serial primary key,
    chat_session_id      integer not null references chat_sessions on delete cascade,
    message              text                                   not null,
    message_type         varchar(9)                             not null,
    time_sent            timestamp without time zone default now() not null,
    token_count          integer                                not null,
    parent_message       integer,
    latest_child_message integer,
    rephrased_query      text,
    prompt_id            integer  references prompt,
    citations            jsonb,
    error                text
);
