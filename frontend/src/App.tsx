import { Router, Route, Navigate } from "@solidjs/router";
import { lazy, onMount } from "solid-js";
import { PublicRoute, ProtectedRoute } from "./components/Layouts";
import AdminLayout from "./components/AdminLayout";
import { useAuth } from "./lib/auth";

const Login = lazy(() => import("./pages/Login"));
const Register = lazy(() => import("./pages/Register"));
const PublicPosts = lazy(() => import("./pages/PublicPosts"));
const PostsToVerify = lazy(() => import("./pages/admin/PostsToVerify"));
const CreatePost = lazy(() => import("./pages/CreatePost"));
const Profile = lazy(() => import("./pages/Profile"));

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
            <Route path="/posts/new" component={CreatePost} />
            <Route path="/profile" component={Profile} />
            <Route path="/admin" component={AdminLayout}>
              <Route
                path="/"
                component={() => <Navigate href='posts/verification'/>}
              />
              <Route
                path="/posts/verification"
                component={PostsToVerify}
              />
            </Route>
          </>
        }
      />
    </Router>
  );
}

export default App;
