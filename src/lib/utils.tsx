import { Persona } from "@/models/settings";
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function capitalize(str: string) {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

export const virtuosoClassName = cn(
  "[&::-webkit-scrollbar]:bg-transparent [&::-webkit-scrollbar]:w-2",
  "[&::-webkit-scrollbar-thumb]:dark:bg-muted [&::-webkit-scrollbar-thumb]:bg-muted-foreground/40 [&::-webkit-scrollbar-thumb]:rounded"
);

export function dataConverter(dateString: string) {
  return new Date(dateString).toLocaleDateString("en-US", {
    month: "long",
    day: "numeric",
    year: "numeric",
  });
}

export function reassembleLLMData (data: Persona[]) {
  for (const record in data){
    data[record].model_id = data[record].llm.model_id
    data[record].endpoint = data[record].llm.endpoint
  }
  return data
}

// export function reassembleLLMInstance (data: Persona) {
//   data.model_id = data.llm.model_id
//   data.endpoint = data.llm.endpoint
//   data.task_prompt = data.prompt?.task_prompt
//   data.system_prompt = data.prompt?.system_prompt
//   data.url = data.prompt.url
//   return data
// }