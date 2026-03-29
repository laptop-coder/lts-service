import { Router, Route } from "@solidjs/router";
import { lazy, onMount } from "solid-js";
import { PublicRoute, ProtectedRoute } from "./components/Layouts";
import { useAuth } from "./lib/auth";

const Login = lazy(() => import("./pages/Login"));
const Register = lazy(() => import("./pages/Register"));
const PublicPosts = lazy(() => import("./pages/PublicPosts"));

function App() {
  const auth = useAuth();
  onMount(() => {
    auth.checkAuth();
  });

  return (
    <Router>
      {/* Public routes */}
      <Route path="/login" component={Login} />
      <Route path="/register" component={Register} />

      {/* Protected routes */}
      <Route
        path="/"
        component={ProtectedRoute}
        children={
          <>
            <Route path="/" component={PublicPosts} />
          </>
        }
      />
    </Router>
  );
}

export default App;
