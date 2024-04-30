import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { RenderTable } from "@/components/renderTable/render-table";
import { SettingHeader } from "@/components/ui/setting-header";
import { Controller } from "./existing-connectors.controller";
import { ConfirmDeleteDialog } from "@/components/dialogs/ConfirmDeleteDialog";
import { useState } from "react";
import { CreateConnectorDialog } from "@/components/dialogs/CreateConnectorDialog";

export function ConnectorsManagementComponent() {
  const { columns, tableData, sortField, handleSortingChange } =
    Controller.useFilterHandler();

  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showCreateTicketDialogOpen, setShowCreateTicketDialogOpen] =
    useState<boolean>(false);

  return (
    <div className="flex flex-grow flex-col m-8">
      <SettingHeader
        title={"Existing Connectors"}
        withBtn
        handleClick={() => {
          setShowCreateTicketDialogOpen(true);
        }}
      />
      <>
        <Tabs defaultValue="personal">
          <TabsList className="mb-4">
            <TabsTrigger value="personal">Personal</TabsTrigger>
            <TabsTrigger value="organizational">Organizational</TabsTrigger>
          </TabsList>
          <TabsContent value="personal">
            <RenderTable
              columns={columns}
              handleSortingChange={handleSortingChange}
              sortField={sortField}
              tableData={tableData}
              onDelete={() => {
                setShowDeleteDialog(true);
              }}
              onEdit={() => {}}
              onPause={() => {}}
              withBtn
            />
          </TabsContent>
          <TabsContent value="organizational">
            <RenderTable
              columns={columns}
              handleSortingChange={handleSortingChange}
              sortField={sortField}
              tableData={tableData}
              onDelete={() => {
                setShowDeleteDialog(true);
              }}
              onEdit={() => {}}
              onPause={() => {}}
            />
          </TabsContent>
        </Tabs>
      </>
      {showDeleteDialog && (
        <div className="ml-auto">
          <ConfirmDeleteDialog
            description="Are you sure you want to delete this Connector?"
            deleteButtonText="Yes, Delete"
            onConfirm={() => {
              console.log("Pressed");
            }}
            open={showDeleteDialog}
            onOpenChange={setShowDeleteDialog}
          />
        </div>
      )}
      {showCreateTicketDialogOpen && (
        <CreateConnectorDialog
          open={showCreateTicketDialogOpen}
          onOpenChange={setShowCreateTicketDialogOpen}
        />
      )}
    </div>
  );
}

export { ConnectorsManagementComponent as Component };
