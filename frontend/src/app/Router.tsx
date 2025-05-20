import { render } from "solid-js/web";
import { Router, Route, Navigate } from "@solidjs/router";
import { HomePage } from "../pages/home/index";
import { StatusPage } from "../pages/status/index";

const root = document.getElementById("root");
render(
  () => (
    <Router>
      <Route
        path="/home"
        component={HomePage}
      />
      <Route
        path="/status"
        component={StatusPage}
      />
      <Route
        path="/*"
        component={() => (
          <Navigate
            href="/home"
            state
          />
        )}
      />
    </Router>
  ),
  root!,
);
