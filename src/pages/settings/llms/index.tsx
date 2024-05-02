import { Tabs, TabsContent } from "@/components/ui/tabs";
import { RenderTable } from "@/components/renderTable/render-table";
import { SettingHeader } from "@/components/ui/setting-header";
import { Controller } from "./llm.controller";
import { ConfirmDeleteDialog } from "@/components/dialogs/ConfirmDeleteDialog";
import { useEffect, useState } from "react";
import { CreateLLMDialog } from "@/components/dialogs/CreateLLMDialog";
import axios from "axios";
import { EditLLMDialog } from "@/components/dialogs/EditLLMDialog";

export function LLMManagementComponent() {
  const [llms, setLlms] = useState([]);
  const [selectedRow, setSelectedRow] = useState<string>("");
  const { columns, sortField, handleSortingChange } =
    Controller.useFilterHandler(llms);

  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showCreateDialogOpen, setShowCreateDialogOpen] =
    useState<boolean>(false);
  const [showEditDialogOpen, setShowEditDialogOpen] = useState<boolean>(false);

  async function getLLMs() {
    await axios
      .get(import.meta.env.VITE_PLATFORM_API_LLM_LIST_URL)
      .then(function (response) {
        if (response.status == 200) {
          setLlms(response.data.data);
        } else {
          setLlms([]);
        }
      })
      .catch(function (error) {
        console.error("Error fetching messages:", error);
      });
  }

  async function deleteLLM(id: string) {
    await axios.post(
      `${import.meta.env.VITE_PLATFORM_API_LLM_DELETE_URL}/${id}/delete`
    );
  }

  useEffect(() => {
    getLLMs();
  }, [showCreateDialogOpen, showDeleteDialog, showEditDialogOpen]);

  return (
    <div className="flex flex-grow flex-col m-8 overflow-x-hidden no-scrollbar">
      <SettingHeader
        title={"LLMs"}
        buttonTitle="New LLM"
        withBtn
        handleClick={() => {
          setShowCreateDialogOpen(true);
        }}
      />
      <>
        <Tabs defaultValue="personal">
          <TabsContent value="personal">
            <RenderTable
              columns={columns}
              handleSortingChange={handleSortingChange}
              sortField={sortField}
              tableData={llms}
              onDelete={(id: string) => {
                setShowDeleteDialog(true);
                setSelectedRow(id);
              }}
              onEdit={() => {
                setShowEditDialogOpen(true);
              }}
              withBtn
            />
          </TabsContent>
        </Tabs>
      </>
      {showDeleteDialog && (
        <div className="ml-auto">
          <ConfirmDeleteDialog
            description="Are you sure you want to delete this LLM?"
            deleteButtonText="Yes, Delete"
            onConfirm={() => {
              deleteLLM(selectedRow);
            }}
            open={showDeleteDialog}
            onOpenChange={setShowDeleteDialog}
          />
        </div>
      )}
      {showCreateDialogOpen && (
        <CreateLLMDialog
          open={showCreateDialogOpen}
          onOpenChange={setShowCreateDialogOpen}
        />
      )}
      {showEditDialogOpen && (
        <EditLLMDialog
          open={showEditDialogOpen}
          onOpenChange={setShowEditDialogOpen}
          // values={}
        />
      )}
    </div>
  );
}

export { LLMManagementComponent as Component };
