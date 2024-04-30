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
  FormLabel,
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
  description: z.string(),
});

export function EditLLMDialog({
  values,
  children,
  open,
  onOpenChange,
}: {
  values?: DefaultValues<z.infer<typeof formSchema>>;
  children?: React.ReactNode;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      ...values,
    },
  });

  const { trigger: triggerEditLLM } = useMutation<LLMSchema>(
    import.meta.env.VITE_PLATFORM_API_LLM_CREATE_URL,
    "POST"
  );

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      await triggerEditLLM({
        name: values.name,
        model_id: values.model_id,
        url: values.url,
        api_key: values.api_key,
        endpoint: values.endpoint,
        system_prompt: values.system_prompt,
        task_prompt: values.task_prompt,
        description: values.description,
      });
    } catch (e) {
      console.log(e);
    }
  };

  return (
    <Dialog
      open={open}
      onOpenChange={(newOpen) => {
        onOpenChange?.(newOpen);
        form.reset();
      }}
    >
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Edit LLM</DialogTitle>
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
                  <FormLabel className="fixed ml-2 -mt-1.5 text-center px-1 bg-white text-muted-foreground">
                    Name
                  </FormLabel>
                  <FormControl>
                    <Input {...field} />
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
                  <FormLabel className="fixed ml-2 -mt-1.5 text-center px-1 bg-white text-muted-foreground">
                    Model ID
                  </FormLabel>
                  <FormControl>
                    <Input {...field} />
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
                  <FormLabel className="fixed ml-2 -mt-1.5 text-center px-1 bg-white text-muted-foreground">
                    URL
                  </FormLabel>

                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="api_key"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="fixed ml-2 -mt-1.5 text-center px-1 bg-white text-muted-foreground">
                    API Key
                  </FormLabel>

                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="endpoint"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="fixed ml-2 -mt-1.5 text-center px-1 bg-white text-muted-foreground">
                    Endpoint
                  </FormLabel>

                  <FormControl>
                    <Input {...field} />
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
                  <FormLabel className="fixed ml-2 mt-1 text-center px-1 bg-white text-muted-foreground">
                    System Prompt
                  </FormLabel>

                  <FormControl>
                    <TextArea
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
                  <FormLabel className="fixed ml-2 mt-1 text-center px-1 bg-white text-muted-foreground">
                    Task Prompt
                  </FormLabel>

                  <FormControl>
                    <TextArea
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
