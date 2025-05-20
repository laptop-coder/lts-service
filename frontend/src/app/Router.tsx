import { render } from "solid-js/web";
import { Router, Route, Navigate } from "@solidjs/router";
import { HomePage } from "../pages/home/index";

const root = document.getElementById("root");
render(
  () => (
    <Router>
      <Route
        path="/home"
        component={HomePage}
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
