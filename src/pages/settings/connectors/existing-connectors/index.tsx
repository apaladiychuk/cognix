import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { RenderTable } from "@/components/renderTable/render-table";
import { SettingHeader } from "@/components/ui/setting-header";
import { Controller } from "./existing-connectors.controller";
import { ConfirmDeleteDialog } from "@/components/dialogs/ConfirmDeleteDialog";
import { useLayoutEffect, useState } from "react";
import { CreateConnectorDialog } from "@/components/dialogs/ConnectorDialog";
import axios from "axios";
import { Connector } from "@/models/settings";

export function ConnectorsManagementComponent() {
  const [ connectors, setConnectors ] = useState<Connector[]>([]);
  const [selectedRow, setSelectedRow] = useState<Connector>();
  const { columns, tableData, sortField, handleSortingChange } =
    Controller.useFilterHandler(connectors);

  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showConnectorDialogOpen, setShowConnectorDialogOpen] =
    useState<boolean>(false);

  async function getConnectors() {
    await axios
      .get(import.meta.env.VITE_PLATFORM_API_CONNECTOR_LIST_URL)
      .then(function (response) {
        if (response.status == 200) {
          setConnectors(response.data.data);
        } else {
          setConnectors([]);
        }
      })
      .catch(function (error) {
        console.error("Error fetching messages:", error);
      });
  }

  async function deleteConnector(id: string) {
    await axios.post(
      `${import.meta.env.VITE_PLATFORM_API_CONNECTOR_DELETE_URL}/${id}/delete`
    ).then((response) => {
      if (response.status == 200){
        setConnectors(connectors.filter((connector) => connector.id !== id));
      }
    });
  }

  async function setRow(id: string) {
    return await axios.get(`${import.meta.env.VITE_PLATFORM_API_CONNECTOR_DETAIL_URL}/${id}`)
    .then(function (response) {
      if (response.status === 200) {
        setSelectedRow(response.data.data);
      } else {
        throw new Error("Failed to fetch detailed connector: " + response.status);
      }
    })
    .catch(function (error) {
      console.error("Error fetching detailed connector:", error);
      throw error; // Re-throw the error to be caught by the caller
    });
  }

  useLayoutEffect(() => {
    getConnectors()
  }, [showConnectorDialogOpen]);

  return (
    <div className="flex flex-grow flex-col m-8 overflow-x-hidden no-scrollbar">
      <SettingHeader
        title={"Existing Connectors"}
        withBtn
        handleClick={() => {
          setShowConnectorDialogOpen(true);
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
              tableData={connectors}
              onDelete={async (id: string) => {
                await setRow(id);
                setShowDeleteDialog(true);
              }}
              onEdit={async (id: string) => {
                await setRow(id)
                setShowConnectorDialogOpen(true);
              }}
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
              deleteConnector(selectedRow!.id)
            }}
            open={showDeleteDialog}
            onOpenChange={setShowDeleteDialog}
          />
        </div>
      )}
      {showConnectorDialogOpen && (
        <CreateConnectorDialog
          open={showConnectorDialogOpen}
          onOpenChange={() => {
            setShowConnectorDialogOpen(false);
            setSelectedRow(undefined);
          }}
          instance={selectedRow}
        />
      )}
    </div>
  );
}

export { ConnectorsManagementComponent as Component };
