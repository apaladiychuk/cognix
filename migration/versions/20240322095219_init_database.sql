-- +goose Up
-- +goose StatementBegin
CREATE TABLE tenants (
    id uuid PRIMARY KEY,     
    name varchar(255),
    configuration jsonb not null default '{}'::jsonb
);
CREATE TABLE users (
    id uuid PRIMARY KEY,
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    user_name varchar(255) UNIQUE NOT NULL,
    first_name varchar(255),
    last_name varchar(255),
    external_id text,
    roles text[] NOT NULL DEFAULT '{}'::text[]
);
CREATE TABLE llm (
    id SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    model_id varchar(255) NOT NULL,
    url varchar NOT NULL
);

CREATE TABLE embedding_models (
    id SERIAL PRIMARY KEY,
    tenant_id uuid NOT NULL,
    model_id varchar NOT NULL,
    model_name varchar NOT NULL,
    model_dim integer NOT NULL,
    normalize boolean NOT NULL,
    query_prefix varchar NOT NULL,
    passage_prefix varchar NOT NULL,
    index_name varchar NOT NULL,
    is_active boolean NOT NULL DEFAULT false
);

CREATE TABLE personas (
    id SERIAL PRIMARY KEY,
    name varchar NOT NULL,
    llm_id integer references llm(id),
    default_persona boolean NOT NULL,
    description varchar NOT NULL,
    tenant_id uuid NOT NULL references tenants(id),
    search_type searchtype NOT NULL,
    is_visible boolean NOT NULL,
    display_priority integer,
    starter_messages jsonb NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE prompts (
    id SERIAL PRIMARY KEY,
    persona_id integer NOT NULL REFERENCES personas(id),
    name varchar NOT NULL,
    description varchar NOT NULL,
    system_prompt text NOT NULL,
    task_prompt text NOT NULL,
    include_citations boolean NOT NULL,
    datetime_aware boolean NOT NULL,
    default_prompt boolean NOT NULL,
    created_date timestamp NOT NULL DEFAULT (now()),
    deleted_date timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE embedding_models ;
DROP TABLE prompts;
DROP TABLE personas ;
DROP TABLE llm;
DROP TABLE users ;
DROP TABLE tenants ;
-- +goose StatementEnd
