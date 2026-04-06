import { type Component, type JSX } from "solid-js";
import { useNavigate, A } from "@solidjs/router";
import { useAuth } from "../lib/auth";
import { usePermissions, PERMISSIONS, ROLES } from "../lib/permissions";

interface Props {
  children?: JSX.Element;
}

export const PublicRoute: Component<Props> = (props) => {
  const auth = useAuth();
  const { hasPermission, hasRole } = usePermissions();

  return (
    <div class="min-h-screen bg-gray-100">
      <header class="bg-white border-b border-gray-200 shadow-sm">
        <div class="container mx-auto px-4 py-3 flex justify-between items-center">
          <A href="/" class="flex items-center gap-3">
            <img
              class="w-10 h-10 rounded-full object-cover"
              src={`/storage/assets/logo.svg`}
              alt="Логотип"
            />
            <span class="text-xl font-bold text-gray-800">
              LostThingsSearch
            </span>
          </A>

          <div class="flex items-center gap-3">
            {hasRole(ROLES.ADMIN) && (
              <A
                href="/admin"
                class="px-3 py-1.5 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition"
              >
                Админка
              </A>
            )}

            {hasPermission(PERMISSIONS.POST_CREATE) && (
              <A
                href="/posts/new"
                class="w-9 h-9 flex items-center justify-center bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition"
              >
                <span class="text-xl font-bold">+</span>
              </A>
            )}

            {auth.user() && (
              <A
                href="/profile"
                class="w-15 h-15 flex bg-gray-100 rounded-full hover:bg-gray-200 transition"
              >
                <img
                  class="w-15 h-15 rounded-full object-cover border-2 border-blue-100 hover:brightness-95 transition"
                  src={`/storage/storage/avatars/${auth.user()?.hasAvatar ? auth.user()?.id : "default"}.jpeg`}
                  alt="Фото профиля"
                />
              </A>
            )}
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
