import { Persona } from "src/models/settings";

export interface ChatMessage {
    chat_session_id: number;
    citations?: number[];
    error?: string;
    id: number;
    latest_child_message?: number;
    message: string;
    message_type: string;
    parent_message?: number;
    rephrased_query?: string;
    time_sent: string;
    token_count: number;
  }
  
  export interface ChatSession {
    created_date: string;
    deleted_date?: string | null;
    description: string;
    id: number;
    messages: ChatMessage[];
    one_shot: boolean;
    persona: Persona;
    persona_id: number;
    user_id: string;
  }

  