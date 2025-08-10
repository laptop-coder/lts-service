import { JSX } from "solid-js";

import { Router, Route, Navigate } from "@solidjs/router";

import HomePage from "./pages/HomePage";
import { HOME_ROUTE } from "./utils/consts";

const AppRouter = (): JSX.Element => {
  return (
    <Router>
      <Route path={HOME_ROUTE} component={() => <HomePage />} />
      <Route path="/*" component={() => <Navigate href={HOME_ROUTE} state />} />
    </Router>
  );
};

export default AppRouter;
