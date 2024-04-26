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
import { CreateConnectorSchema } from "@/lib/schemas/connectors";

const formSchema = z.object({
  connector_id: z.number(),
  source: z.string(),
  connector_specific_config: z.record(z.string()),
  refresh_freq: z.union([z.string().email(), z.literal("")]),
  credential_id: z.number(),
});

export function CreateConnectorDialog({
  defaultValues,
  children,
  open,
  onOpenChange,
}: {
  defaultValues?: DefaultValues<z.infer<typeof formSchema>>;
  children?: React.ReactNode;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
}) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      source: "",
      connector_specific_config: {},
      refresh_freq: "",
      credential_id: undefined,
      connector_id: undefined,
      ...defaultValues,
    },
  });

  //   useEffect(() => {
  //     if (activeTab === 'existing') {
  //       form.setValue('new_customer_email', '');
  //     } else if (activeTab === 'create') {
  //       form.setValue('customer_id', '');
  //     }
  //   }, [activeTab]);

  const { trigger: triggerCreateConnector } =
    useMutation<CreateConnectorSchema>(
      import.meta.env.VITE_PLATFORM_API_CONNECTOR_CREATE_URL,
      "POST"
    );

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      await triggerCreateConnector({
        source: values.source,
        connector_specific_config: values.connector_specific_config,
        refresh_freq: values.refresh_freq,
        credential_id: values.credential_id,
        connector_id: values.connector_id,
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
          <DialogTitle>Add Connector</DialogTitle>
        </DialogHeader>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="max-w-full space-y-4 overflow-hidden px-0.5"
            onKeyDown={(e) => e.key === "Enter" && e.preventDefault()}
          >
            <FormField
              control={form.control}
              name="connector_id"
              render={(field) => (
                <FormItem>
                  <FormLabel>Connector</FormLabel>
                  <FormControl>
                    <Input value={field.field.value} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="source"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input placeholder="Source" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="connector_specific_config"
              render={({...field}) => (
                <FormItem>
                  <FormControl>
                    <textarea
                      className="h-28 w-full p-3 rounded-md"
                      style={{ resize: "none" }}
                      placeholder="Connector Specific Configuration"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="refresh_freq"
              render={(field) => (
                <FormItem>
                  <FormControl>
                    <Input placeholder="Refresh Frequency" type="text" {...field}/>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="credential_id"
              render={(field) => (
                <FormItem>
                  <FormControl>
                    <textarea
                      className="h-28 w-full p-3 rounded-md"
                      style={{ resize: "none" }}
                      placeholder="Connector credential"
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
                  })} w-full`}
                >
                  Cancel
                </Button>
              </DialogClose>
              <Button type="submit" className="w-full">
                Add
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
