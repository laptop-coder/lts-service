import { Router, Route, Navigate } from "@solidjs/router";
import { lazy, onMount } from "solid-js";
import { PublicRoute, ProtectedRoute } from "./components/Layouts";
import AdminLayout from "./components/AdminLayout";
import { useAuth } from "./lib/auth";
import { ROLES, usePermissions } from "./lib/permissions";

const Login = lazy(() => import("./pages/Login"));
const Register = lazy(() => import("./pages/Register"));
const PublicPosts = lazy(() => import("./pages/PublicPosts"));
const PostsToVerify = lazy(() => import("./pages/admin/PostsToVerify"));
const CreatePost = lazy(() => import("./pages/CreatePost"));
const EditPost = lazy(() => import("./pages/EditPost"));
const Profile = lazy(() => import("./pages/Profile"));
const PublicProfile = lazy(() => import("./pages/PublicProfile"));
const Subjects = lazy(() => import("./pages/admin/Subjects"));
const Rooms = lazy(() => import("./pages/admin/Rooms"));
const StudentGroups = lazy(() => import("./pages/admin/StudentGroups"));
const InviteTokens = lazy(() => import("./pages/InviteTokens"));
const Users = lazy(() => import("./pages/Users"));
const StaffPositions = lazy(() => import("./pages/admin/StaffPositions"));
const InstitutionAdministratorPositions = lazy(
  () => import("./pages/admin/InstitutionAdministratorPositions"),
);
const About = lazy(() => import("./pages/About"));
const PostDetails = lazy(() => import("./pages/PostDetails"));
const ListOfConversations = lazy(() => import("./pages/ListOfConversations"));
const ConversationView = lazy(() => import("./pages/ConversationView"));

function App() {
  const auth = useAuth();
  onMount(() => {
    auth.checkAuth();
  });
  const { hasRole } = usePermissions();

  return (
    <Router>
      {/* Public routes */}
      <Route path="/login" component={Login} />
      <Route path="/register" component={Register} />
      <Route
        path="/"
        component={PublicRoute}
        children={<Route path="/" component={PublicPosts} />}
      />
      <Route
        path="/posts/:id"
        component={PublicRoute}
        children={<Route path="/" component={PostDetails} />}
      />
      <Route
        path="/about"
        component={PublicRoute}
        children={<Route path="/" component={About} />}
      />

      {/* Protected routes */}
      <Route
        path="/"
        component={ProtectedRoute}
        children={
          <>
            <Route path="/posts/new" component={CreatePost} />
            <Route path="/posts/:id/edit" component={EditPost} />
            <Route path="/profile" component={Profile} />
            <Route path="/users/:id" component={PublicProfile} />
            <Route path="/conversations" component={ListOfConversations} />
            <Route path="/conversations/:id" component={ConversationView} />
            <Route path="/admin" component={AdminLayout}>
              <Route
                path="/"
                component={() =>
                  hasRole(ROLES.ADMIN) ? (
                    <Navigate href="posts/verification" />
                  ) : (
                    <Navigate href="users" />
                  )
                }
              />
              <Route path="/posts/verification" component={PostsToVerify} />
              <Route path="/subjects" component={Subjects} />
              <Route path="/rooms" component={Rooms} />
              <Route path="/student_groups" component={StudentGroups} />
              <Route path="/invite_tokens" component={InviteTokens} />
              <Route path="/users" component={Users} />
              <Route path="/positions/staff" component={StaffPositions} />
              <Route
                path="/positions/institution_administrators"
                component={InstitutionAdministratorPositions}
              />
            </Route>
          </>
        }
      />
    </Router>
  );
}

export default App;
