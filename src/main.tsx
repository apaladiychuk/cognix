import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import "@/global.css";
import "@/lib/axios";

export const router = createBrowserRouter([
  {
    path: "/",
    lazy: () => import("@/pages/login"),
  },
  {
    path: "/",
    children: [
      {
        path: "google/callback",
        lazy: () => import("@/pages/login/redirect"),
      },
      {
        path: "/platform",
        lazy: () => import("@/pages/platform"),
        children: [
          {
            path: "/platform",
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
                    lazy: () =>
                      import("@/pages/settings/connectors/add-connector"),
                  },
                ],
              },
              {
                path: "feedback",
                lazy: () => import("@/pages/settings/feedback"),
              },
              {
                path: "embeddings",
                lazy: () => import("@/pages/settings/embeddings"),
              },
              {
                path: "llms",
                lazy: () => import("@/pages/settings/llms"),
              },
              {
                path: "users",
                lazy: () => import("@/pages/settings/users"),
              },
              {
                path: "config",
                lazy: () => import("@/pages/settings/config"),
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
