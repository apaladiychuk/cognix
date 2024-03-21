CREATE TYPE "searchtype" AS ENUM (
  'KEYWORD',
  'SEMANTIC',
  'HYBRID'
);

CREATE TYPE "recencybiassetting" AS ENUM (
  'FAVOR_RECENT',
  'BASE_DECAY',
  'NO_DECAY',
  'AUTO'
);

CREATE TABLE "connectors" (
    "id" SERIAL PRIMARY KEY,
    "name" varchar NOT NULL,
    "source" varchar(50) NOT NULL,
    "input_type" varchar(10),
    "connector_specific_config" jsonb NOT NULL,
    "refresh_freq" integer,
    "user_id" uuid NOT NULL,
    "tenant_id" uuid NOT NULL,
    "created_date" timestamp NOT NULL DEFAULT (now()),
    "deleted_date" timestamp
);

CREATE TABLE "credentials" (
        "id" SERIAL PRIMARY KEY,
        "credential_json" jsonb NOT NULL,
        "user_id" uuid NOT NULL,
        "tenant_id" uuid NOT NULL,
        "source" varchar(50) NOT NULL,
        "created_date" timestamp NOT NULL DEFAULT (now()),
        "updated_date" timestamp,
        "deleted_date" timestamp,
        "admin_public" boolean NOT NULL
);

CREATE TABLE "connector_credential_pairs" (
                                              "connector_id" integer NOT NULL,
                                              "credential_id" integer NOT NULL,
                                              "last_successful_index_time" timestamp,
                                              "last_attempt_status" varchar,
                                              "total_docs_indexed" integer NOT NULL,
                                              "name" varchar NOT NULL,
                                              "is_public" boolean NOT NULL,
                                              PRIMARY KEY ("connector_id", "credential_id")
);

CREATE TABLE "embedding_models" (
                                    "id" SERIAL PRIMARY KEY,
                                    "tenant_id" uuid NOT NULL,
                                    "model_id" varchar NOT NULL,
                                    "model_name" varchar NOT NULL,
                                    "model_dim" integer NOT NULL,
                                    "normalize" boolean NOT NULL,
                                    "query_prefix" varchar NOT NULL,
                                    "passage_prefix" varchar NOT NULL,
                                    "index_name" varchar NOT NULL,
                                    "is_active" boolean NOT NULL DEFAULT false
);

CREATE TABLE "index_attempts" (
                                  "id" SERIAL PRIMARY KEY,
                                  "created_at" timestamp NOT NULL DEFAULT (now()),
                                  "status" varchar NOT NULL,
                                  "error_msg" varchar,
                                  "connector_id" integer,
                                  "credential_id" integer,
                                  "total_docs_indexed" integer,
                                  "time_started" timestamp,
                                  "new_docs_indexed" integer,
                                  "embedding_model_id" integer NOT NULL,
                                  "from_beginning" boolean NOT NULL,
                                  "full_exception_trace" text,
                                  "docs_removed_from_index" integer
);

CREATE TABLE "llm" (
                       "id" SERIAL PRIMARY KEY,
                       "name" varchar NOT NULL,
                       "model_id" varchar NOT NULL,
                       "url" varchar NOT NULL
);

CREATE TABLE "personas" (
                            "id" SERIAL PRIMARY KEY,
                            "name" varchar NOT NULL,
                            "llm_id" integer,
                            "default_persona" boolean NOT NULL,
                            "description" varchar NOT NULL,
                            "tenant_id" uuid NOT NULL,
                            "search_type" searchtype NOT NULL,
                            "is_visible" boolean NOT NULL,
                            "display_priority" integer,
                            "starter_messages" jsonb
);

CREATE TABLE "prompts" (
                           "id" SERIAL PRIMARY KEY,
                           "user_id" uuid NOT NULL,
                           "name" varchar NOT NULL,
                           "description" varchar NOT NULL,
                           "system_prompt" text NOT NULL,
                           "task_prompt" text NOT NULL,
                           "include_citations" boolean NOT NULL,
                           "datetime_aware" boolean NOT NULL,
                           "default_prompt" boolean NOT NULL,
                           "created_date" timestamp NOT NULL DEFAULT (now()),
                           "deleted_date" timestamp
);

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
                                 "prompt_id" integer,
                                 "citations" jsonb,
                                 "error" text
);

CREATE TABLE "users" (
                         "id" uuid PRIMARY KEY,
                         "tenant_id" uuid,
                         "user_name" text UNIQUE NOT NULL,
                         "first_name" text,
                         "last_name" text,
                         "roles" jsonb
);

CREATE TABLE "tenants" (
                           "id" uuid PRIMARY KEY,
                           "name" text
);

ALTER TABLE "connector_credential_pairs" ADD FOREIGN KEY ("connector_id") REFERENCES "connectors" ("id") ON DELETE CASCADE;

ALTER TABLE "connector_credential_pairs" ADD FOREIGN KEY ("credential_id") REFERENCES "credentials" ("id") ON DELETE CASCADE;

ALTER TABLE "index_attempts" ADD FOREIGN KEY ("connector_id") REFERENCES "connectors" ("id");

ALTER TABLE "index_attempts" ADD FOREIGN KEY ("credential_id") REFERENCES "credentials" ("id");

ALTER TABLE "index_attempts" ADD FOREIGN KEY ("embedding_model_id") REFERENCES "embedding_models" ("id");

ALTER TABLE "personas" ADD FOREIGN KEY ("llm_id") REFERENCES "llm" ("id");

ALTER TABLE "chat_sessions" ADD FOREIGN KEY ("persona_id") REFERENCES "personas" ("id");

ALTER TABLE "chat_messages" ADD FOREIGN KEY ("chat_session_id") REFERENCES "chat_sessions" ("id") ON DELETE CASCADE;

ALTER TABLE "chat_messages" ADD FOREIGN KEY ("prompt_id") REFERENCES "prompts" ("id");

ALTER TABLE "chat_sessions" ADD FOREIGN KEY ("deleted_date") REFERENCES "personas" ("is_visible");

ALTER TABLE "users" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");