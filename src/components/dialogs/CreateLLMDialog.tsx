import { Button, buttonVariants } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { DefaultValues, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input";
import { useMutation } from "@/lib/mutation";
import { LLMSchema } from "@/lib/schemas/llms";
import { TextArea } from "../ui/textarea";

const formSchema = z.object({
  name: z.string(),
  model_id: z.string(),
  url: z.string(),
  api_key: z.string(),
  endpoint: z.string(),
  system_prompt: z.string(),
  task_prompt: z.string(),
  description: z.string()
});

export function CreateLLMDialog({
  defaultValues,  
  children,
  open,
  onOpenChange,
}: {
  defaultValues?: DefaultValues<z.infer<typeof formSchema>>;
  children?: React.ReactNode;
  open?: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      model_id: "",
      url: "",
      api_key: "",
      endpoint: "",
      system_prompt: "",
      task_prompt: "",
      description: "blank",
      ...defaultValues,
    },
  });

  const { trigger: triggerCreateLLM } =
    useMutation<LLMSchema>(
      import.meta.env.VITE_PLATFORM_API_LLM_CREATE_URL,
      "POST"
    );

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      await triggerCreateLLM({
        name: values.name,
        model_id: values.model_id,
        url: values.url,
        api_key: values.api_key,
        endpoint: values.endpoint,
        system_prompt: values.system_prompt,
        task_prompt: values.task_prompt,
        description: values.description,
      });
      onOpenChange(false)
    } catch (e) {
      console.log(e);
    }
  };

  return (
    <Dialog
      open={open}
      onOpenChange={() => {
        onOpenChange?.(false);
        form.reset();
      }}
    >
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Add New LLM</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="max-w-full space-y-4 overflow-hidden px-0.5 bg-white"
            onKeyDown={(e) => e.key === "Enter" && e.preventDefault()}
          >
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input placeholder="Name" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="model_id"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input placeholder="Model ID" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="url"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input placeholder="URL" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="api_key"
              render={({field}) => (
                <FormItem>
                  <FormControl>
                    <Input
                       placeholder="API Key"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="endpoint"
              render={({field}) => (
                <FormItem>
                  <FormControl>
                    <Input placeholder="Endpoint" {...field}/>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="system_prompt"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <TextArea
                      placeholder="System Prompt"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="task_prompt"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <TextArea
                      placeholder="Task Prompt"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <DialogClose asChild>
                <Button
                  variant="outline"
                  className={`${buttonVariants({
                    variant: "destructive",
                  })} text-md w-full h-10`}
                >
                  Cancel
                </Button>
              </DialogClose>
              <Button type="submit" className="text-md w-full h-10">
                Add
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
