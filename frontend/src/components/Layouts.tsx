import { type Component, type JSX, Show, createEffect, createSignal, onMount, onCleanup } from "solid-js";
import { useNavigate, useLocation, A } from "@solidjs/router";
import { useAuth } from "../lib/auth";
import { usePermissions, PERMISSIONS, ROLES } from "../lib/permissions";
import { unreadMessagesCount, refreshUnreadMessagesCount } from "../lib/store";
import { MessageSquareText, Settings2, Plus } from "lucide-solid";

interface Props {
  children?: JSX.Element;
}

export const PublicRoute: Component<Props> = (props) => {
  const auth = useAuth();
  const location = useLocation();
  const { hasPermission, hasRole, hasAnyRole } = usePermissions();

  createEffect(async () => {
    if (auth.user()) {
      await refreshUnreadMessagesCount();
    }
  });

  const [isMobile, setIsMobile]= createSignal(window.innerWidth < 768)

  onMount(() => {
    const mq = window.matchMedia("(max-width: 767px)")
    const handler = (e: MediaQueryListEvent) => setIsMobile(e.matches)
    mq.addEventListener("change", handler)
    onCleanup(() => mq.removeEventListener("change", handler))
  })

  // /conversations/<any_string>
  const isMobileConversationViewPage = () =>
    /^\/conversations\/[^/]+$/.test(location.pathname) &&
    isMobile();

  return (
    <div class="min-h-screen bg-gray-50 flex flex-col">
      <header class="bg-white border-b border-gray-200 shadow-sm sticky top-0 z-30">
        <div class="container mx-auto px-4 h-16 flex justify-between items-center">
          <A href="/" class="flex items-center gap-3">
            <img
              class="w-10 h-10 rounded-full object-cover"
              src={`/storage/assets/logo.svg`}
              alt="Логотип"
            />
            <span class="text-xl font-bold text-gray-800">
              <span class="hidden md:block">LostThingsSearch</span>
              <span class="block md:hidden">LTS</span>
            </span>
          </A>

          <div class="flex items-center gap-3">
            {auth.user() ? (
              <>
                {hasAnyRole(ROLES.ADMIN, ROLES.SUPERADMIN) && (
                  <A
                    href="/admin"
                    class="bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition flex items-center justify-center w-10 h-10"
                  >
                    <Settings2 />
                  </A>
                )}

                {hasPermission(PERMISSIONS.POST_CREATE) && (
                  <A
                    href="/posts/new"
                    class="bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition flex items-center justify-center w-10 h-10"
                  >
                    <Plus />
                  </A>
                )}

                <A
                  href="/conversations"
                  class="bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition flex items-center justify-center w-10 h-10 relative"
                >
                  <MessageSquareText />
                  <Show when={unreadMessagesCount() > 0}>
                    <div class="bg-blue-600 text-white text-xs font-medium px-2 py-1 rounded-full absolute -top-1 -right-1">
                      {unreadMessagesCount()}
                    </div>
                  </Show>
                </A>

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
              </>
            ) : (
              <Show when={location.pathname !== "/login"}>
                <A
                  href="/login"
                  class="px-4 py-1.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
                >
                  Войти
                </A>
              </Show>
            )}
          </div>
        </div>
      </header>

      {/*Mobile chat page*/}
      <Show when={isMobileConversationViewPage()}>
        <div class="fixed inset-0 top-0 pt-16 z-10">{props?.children}</div>
      </Show>

      <Show when={!isMobileConversationViewPage()}>
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
      </Show>
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
