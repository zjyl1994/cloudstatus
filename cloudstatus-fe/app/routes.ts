import { type RouteConfig, index } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  { path: ":nodeId", file: "routes/charts.tsx" }
] satisfies RouteConfig;
