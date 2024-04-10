CREATE TYPE "searchtype" AS ENUM ('KEYWORD','SEMANTIC', 'HYBRID');

-- CREATE TABLE "connectors" (
--   "id" SERIAL PRIMARY KEY,
--   "credential_id" integer NOT NULL,
--   "name" varchar NOT NULL,
--   "source" varchar(50) NOT NULL,
--   "input_type" varchar(10),
--   "connector_specific_config" jsonb NOT NULL,
--   "refresh_freq" integer,
--   "user_id" uuid NOT NULL,
--   "tenant_id" uuid NOT NULL,
--   "last_successful_index_time" timestamp,
--   "last_attempt_status" varchar,
--   "total_docs_indexed" integer NOT NULL,
--   "created_date" timestamp NOT NULL DEFAULT (now()),
--   "deleted_date" timestamp
-- );
--
-- CREATE TABLE "credentials" (
--   "id" SERIAL PRIMARY KEY,
--   "credential_json" jsonb NOT NULL,
--   "user_id" uuid NOT NULL,
--   "tenant_id" uuid NOT NULL,
--   "source" varchar(50) NOT NULL,
--   "created_date" timestamp NOT NULL DEFAULT (now()),
--   "updated_date" timestamp,
--   "deleted_date" timestamp,
--   "admin_public" boolean NOT NULL
-- );

-- CREATE TABLE "embedding_models" (
--   "id" SERIAL PRIMARY KEY,
--   "tenant_id" uuid NOT NULL,
--   "model_id" varchar NOT NULL,
--   "model_name" varchar NOT NULL,
--   "model_dim" integer NOT NULL,
--   "normalize" boolean NOT NULL,
--   "query_prefix" varchar NOT NULL,
--   "passage_prefix" varchar NOT NULL,
--   "index_name" varchar NOT NULL,
--   "url" varchar,
--   "is_active" boolean NOT NULL DEFAULT false
-- );

-- CREATE TABLE "llm" (
--   "id" SERIAL PRIMARY KEY,
--   "name" varchar NOT NULL,
--   "model_id" varchar NOT NULL,
--   "url" varchar NOT NULL,
--   "api_key" varchar,
--   "endpoint" varchar
-- );

-- CREATE TABLE "personas" (
--   "id" SERIAL PRIMARY KEY,
--   "name" varchar NOT NULL,
--   "llm_id" integer,
--   "default_persona" boolean NOT NULL,
--   "description" varchar NOT NULL,
--   "tenant_id" uuid NOT NULL,
--   "search_type" searchtype NOT NULL,
--   "is_visible" boolean NOT NULL,
--   "display_priority" integer,
--   "starter_messages" jsonb
-- );
--
-- CREATE TABLE "prompts" (
--   "id" SERIAL PRIMARY KEY,
--   "persona_id" integer NOT NULL,
--   "user_id" uuid NOT NULL,
--   "name" varchar NOT NULL,
--   "description" varchar NOT NULL,
--   "system_prompt" text NOT NULL,
--   "task_prompt" text NOT NULL,
--   "include_citations" boolean NOT NULL,
--   "datetime_aware" boolean NOT NULL,
--   "default_prompt" boolean NOT NULL,
--   "created_date" timestamp NOT NULL DEFAULT (now()),
--   "deleted_date" timestamp
-- );

CREATE TABLE "chat_sessions" (
  "id" SERIAL PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "description" text NOT NULL,
  "created_date" timestamp NOT NULL DEFAULT (now()),
  "deleted_date" timestamp,
  "persona_id" integer NOT NULL,
  "one_shot" boolean NOT NULL
);

CREATE TABLE "chat_messages" (
  "id" SERIAL PRIMARY KEY,
  "chat_session_id" integer NOT NULL,
  "message" text NOT NULL,
  "message_type" varchar(9) NOT NULL,
  "time_sent" timestamp NOT NULL DEFAULT (now()),
  "token_count" integer NOT NULL,
  "parent_message" integer,
  "latest_child_message" integer,
  "rephrased_query" text,
  "citations" jsonb,
  "error" text
);

-- CREATE TABLE "users" (
--   "id" uuid PRIMARY KEY,
--   "tenant_id" uuid,
--   "user_name" text UNIQUE NOT NULL,
--   "first_name" text,
--   "last_name" text,
--   "external_id" text,
--   "roles" text[]
-- );

-- CREATE TABLE "tenants" (
--   "id" uuid PRIMARY KEY,
--   "name" text,
--   "configuration" jsonb
-- );

CREATE TABLE "document" (
  "id" varchar PRIMARY KEY NOT NULL,
  "connector_id" integer NOT NULL,
  "boost" integer NOT NULL,
  "hidden" boolean NOT NULL,
  "semantic_id" varchar NOT NULL,
  "link" varchar,
  "doc_updated_at" timestamp,
  "from_ingestion_api" boolean,
  "signature" text
);

CREATE TABLE "document_set" (
  "id" SERIAL PRIMARY KEY,
  "user_id" uuid,
  "name" varchar NOT NULL,
  "description" varchar NOT NULL,
  "is_up_to_date" boolean NOT NULL
);

CREATE TABLE "document_set_connector_pair" (
  "id" SERIAL PRIMARY KEY,
  "document_set_id" integer NOT NULL,
  "connector_id" integer NOT NULL,
  "is_current" boolean NOT NULL
);

ALTER TABLE "personas" ADD FOREIGN KEY ("llm_id") REFERENCES "llm" ("id");

ALTER TABLE "chat_sessions" ADD FOREIGN KEY ("persona_id") REFERENCES "personas" ("id");

ALTER TABLE "chat_messages" ADD FOREIGN KEY ("chat_session_id") REFERENCES "chat_sessions" ("id") ON DELETE CASCADE;

ALTER TABLE "chat_sessions" ADD FOREIGN KEY ("deleted_date") REFERENCES "personas" ("is_visible");

ALTER TABLE "users" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "personas" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "chat_sessions" ADD FOREIGN KEY ("id") REFERENCES "users" ("id");

ALTER TABLE "document" ADD FOREIGN KEY ("connector_id") REFERENCES "connectors" ("id");

ALTER TABLE "document_set" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "document_set_connector_pair" ADD FOREIGN KEY ("document_set_id") REFERENCES "document_set" ("id");

ALTER TABLE "document_set_connector_pair" ADD FOREIGN KEY ("connector_id") REFERENCES "connectors" ("id");

ALTER TABLE "embedding_models" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "connectors" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "credentials" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "credentials" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "connectors" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "prompts" ADD FOREIGN KEY ("persona_id") REFERENCES "personas" ("id");

ALTER TABLE "prompts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "connectors" ADD FOREIGN KEY ("credential_id") REFERENCES "credentials" ("id");
