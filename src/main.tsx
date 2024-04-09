import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
// import "@/global.css";
// import "@fontsource-variable/inter";

const router = createBrowserRouter([
  {
    path: "/",
    children: [
      // {
      //   path: "login",
      //   children: [
      //     {
      //       path: "sign-in",
      //       lazy: () => import("@/pages/login/sign-in"),
      //     },
      //     {
      //       path: "sign-up",
      //       lazy: () => import("@/pages/login/sign-up"),
      //     },
      //     {
      //       path: "onboarding",
      //       lazy: () => import("@/pages/login/onboarding"),
      //     },
      //     {
      //       path: "user",
      //       lazy: () => import("@/pages/login/user"),
      //     },
      //   ],
      // },
      {
        path: "chat",
        lazy: () => import("@/pages/chat"),
      },
      {
        path: "settings",
        children: [
          {
            path: "connectors",
            children: [
              {
                path: "existing-connectors",
                lazy: () =>
                  import("@/pages/settings/connectors/existing-connectors"),
              },
              {
                path: "add-connector",
                lazy: () => import("@/pages/settings/connectors/add-connector"),
              },
            ],
          },
          {
            path: "documents",
            children: [
              {
                path: "document-sets",
                lazy: () => import("@/pages/settings/documents/sets"),
              },
              {
                path: "explorer",
                lazy: () => import("@/pages/settings/documents/explorer"),
              },
              {
                path: "feedback",
                lazy: () => import("@/pages/settings/documents/feedback"),
              },
            ],
          },
          {
            path: "custom-assistant",
            children: [
              {
                path: "personas",
                lazy: () =>
                  import("@/pages/settings/custom-assistant/personas"),
              },
              {
                path: "slack-bots",
                lazy: () =>
                  import("@/pages/settings/custom-assistant/slack-bots"),
              },
              {
                path: "teams",
                lazy: () => import("@/pages/settings/custom-assistant/teams"),
              },
            ],
          },
          {
            path: "model-config",
            children: [
              {
                path: "llms",
                lazy: () => import("@/pages/settings/model-config/llms"),
              },
              {
                path: "embedding",
                lazy: () => import("@/pages/settings/model-config/embedding"),
              },
            ],
          },
          {
            path: "user-management",
            children: [
              {
                path: "users",
                lazy: () => import("@/pages/settings/user-management/users"),
              },
            ],
          },
        ],
      },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
