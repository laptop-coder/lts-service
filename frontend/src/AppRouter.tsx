import { JSX } from "solid-js";

import { Router, Route, Navigate } from "@solidjs/router";

import HomePage from "./pages/HomePage";

const AppRouter = (): JSX.Element => {
  return (
    <Router>
      <Route path="/" component={() => <HomePage />} />
      <Route path="/*" component={() => <Navigate href="/" state />} />
    </Router>
  );
};

export default AppRouter;
