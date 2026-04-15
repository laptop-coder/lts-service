import { type Component, type JSX } from "solid-js";
import { useNavigate, A } from "@solidjs/router";
import { useAuth } from "../lib/auth";
import { usePermissions, PERMISSIONS, ROLES } from "../lib/permissions";

interface Props {
  children?: JSX.Element;
}

export const PublicRoute: Component<Props> = (props) => {
  const auth = useAuth();
  const { hasPermission, hasRole, hasAnyRole } = usePermissions();

  return (
    <div class="min-h-screen bg-gray-50 flex flex-col">
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
            {hasAnyRole(ROLES.ADMIN, ROLES.SUPERADMIN) && (
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

            {auth.user() ? (
              <A
                href="/profile"
                class="w-10 h-10 flex bg-gray-100 rounded-full hover:bg-gray-200 transition"
              >
                <img
                  class="w-10 h-10 rounded-full object-cover border-2 border-blue-100 hover:brightness-95 transition"
                  src={`/storage/storage/avatars/${auth.user()?.hasAvatar ? auth.user()?.id : "default"}.jpeg`}
                  alt="Фото профиля"
                />
              </A>
            ) : (
              <A
                href="/login"
                class="px-4 py-1.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
              >
                Войти
              </A>
            )}
          </div>
        </div>
      </header>

      <main class="container mx-auto px-4 py-8" flex-1>
        {props?.children}
      </main>

      <footer class="bg-white border-t border-gray-200 mt-auto">
        <div class="container mx-auto px-4 py-6">
          <div class="flex flex-col md:flex-row justify-between items-center gap-4">
            <div class="text-sm text-gray-500">
              © {new Date().getFullYear()} LostThingsSearch.
            </div>

            <div class="flex gap-6">
              <a
                href="/about"
                class="text-sm text-gray-500 hover:text-gray-700 transition"
              >
                О проекте
              </a>
            </div>
          </div>
        </div>
      </footer>
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
