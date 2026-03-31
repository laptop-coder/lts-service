import { type Component, type JSX } from "solid-js";
import { useNavigate, A } from "@solidjs/router";
import { useAuth } from "../lib/auth";
import { usePermissions, PERMISSIONS, ROLES } from "../lib/permissions";

interface Props {
  children?: JSX.Element;
}

export const PublicRoute: Component<Props> = (props) => {
  const navigate = useNavigate();
  const auth = useAuth();
  const { hasPermission, hasRole } = usePermissions();

  const handleLogout = async () => {
    await auth.logout();
    navigate("/login");
  };

  return (
    <div class="min-h-screen bg-gray-100">
      <header class="bg-white shadow">
        <div class="container mx-auto px-4 py-4 flex justify-between items-center">
          <A href="/">
            <h1 class="text-xl font-bold">Сервис поиска потерянных вещей</h1>
          </A>
          <div class="flex items-center gap-4">
            {hasRole(ROLES.ADMIN) && (
              <A class="text-xl font-bold" href="/admin">
                А
              </A>
            )}
            {hasPermission(PERMISSIONS.POST_CREATE) && (
              <A class="text-xl font-bold" href="/posts/new">
                +
              </A>
            )}
            <A href="/profile">
              {auth.user()?.firstName} {auth.user()?.lastName}
            </A>
            <button
              onClick={handleLogout}
              class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
            >
              Выйти
            </button>
          </div>
        </div>
      </header>
      <main class="container mx-auto px-4 py-8">{props?.children}</main>
    </div>
  );
};

export const ProtectedRoute: Component<Props> = (props) => {
  const navigate = useNavigate();
  const auth = useAuth();

  return (
    <>
      {auth.isLoading() ? (
        <div class="flex justify-center items-center h-screen">Загрузка...</div>
      ) : auth.isAuthenticated() ? (
        <PublicRoute>{props?.children}</PublicRoute>
      ) : (
        navigate("/login")
      )}
    </>
  );
};
