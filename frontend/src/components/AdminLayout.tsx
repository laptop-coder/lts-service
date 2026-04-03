import { For, Show, createSignal } from "solid-js";
import { useNavigate, useLocation } from "@solidjs/router";
import { usePermissions, ROLES } from "../lib/permissions";

interface Props {
  children?: any;
}

interface SubTab {
  key: string;
  label: string;
  path: string;
}

interface Tab {
  key: string;
  label: string;
  path?: string;
  subTabs?: SubTab[];
}

const tabs: Tab[] = [
  {
    key: "posts",
    label: "Объявления",
    subTabs: [
      {
        key: "posts-verification",
        label: "Верификация",
        path: "/admin/posts/verification",
      },
    ],
  },
  {
    key: "users",
    label: "Пользователи",
    path: "/admin/users",
  },
  {
    key: "subjects",
    label: "Предметы",
    path: "/admin/subjects",
  },
  {
    key: "rooms",
    label: "Кабинеты",
    path: "/admin/rooms",
  },
  {
    key: "student-groups",
    label: "Учебные группы",
      path: "/admin/student_groups",
  },
  {
    key: "positions",
    label: "Должности",
    subTabs: [
      {
      key: "positions-institution-administrators",
      label: "Администрация",
      path: "/admin/positions/institution_administrators"
      },
      {
      key: "positions-staff",
      label: "Сотрудники",
      path: "/admin/positions/staff"
      }
    ]
  },
  {
    key: "invite-tokens",
    label: "Инвайт-токены",
    path: "/admin/invite_tokens",
  },
];

const AdminLayout = (props: Props) => {
  const navigate = useNavigate();
  const location = useLocation();
  const [openSubmenu, setOpenSubmenu] = createSignal<string | null>(null);
  const { hasRole } = usePermissions();

  const isActive = (path: string) => location.pathname === path;
  const isParentActive = (tab: Tab) => {
    if (tab.path && isActive(tab.path)) return true;
    if (tab.subTabs) {
      return tab.subTabs.some((sub) => isActive(sub.path));
    }
    return false;
  };

  const toggleSubmenu = (key: string) => {
    setOpenSubmenu(openSubmenu() === key ? null : key);
  };

  return (
    <>
      {hasRole(ROLES.ADMIN) && (
        <div class="flex min-h-screen bg-gray-100">
          <aside class="w-64 bg-white shadow-lg overflow-y-auto">
            <div class="p-4 border-b">
              <h2 class="text-xl font-bold">Панель администратора</h2>
            </div>
            <nav class="p-2">
              <div class="space-y-1">
                <For each={tabs}>
                  {(tab) => (
                    <div>
                      {tab.path ? (
                        <button
                          onClick={() => navigate(tab.path!)}
                          class={`
                        w-full text-left px-4 py-2 rounded-lg transition-colors
                        ${isActive(tab.path) ? "bg-blue-50 text-blue-700" : "text-gray-600 hover:bg-gray-100"}
                      `}
                        >
                          {tab.label}
                        </button>
                      ) : (
                        <button
                          onClick={() => toggleSubmenu(tab.key)}
                          class={`
                        w-full text-left px-4 py-2 rounded-lg transition-colors flex justify-between items-center
                        ${isParentActive(tab) ? "bg-blue-50 text-blue-700" : "text-gray-600 hover:bg-gray-100"}
                      `}
                        >
                          <span>{tab.label}</span>
                          <span>{openSubmenu() === tab.key ? "▼" : "▶"}</span>
                        </button>
                      )}

                      <Show when={tab.subTabs && openSubmenu() === tab.key}>
                        <div class="ml-4 mt-1 space-y-1">
                          <For each={tab.subTabs}>
                            {(sub) => (
                              <button
                                onClick={() => navigate(sub.path)}
                                class={`
                              w-full text-left px-4 py-2 rounded-lg text-sm transition-colors
                              ${isActive(sub.path) ? "bg-blue-50 text-blue-700" : "text-gray-600 hover:bg-gray-100"}
                            `}
                              >
                                {sub.label}
                              </button>
                            )}
                          </For>
                        </div>
                      </Show>
                    </div>
                  )}
                </For>
              </div>
            </nav>
          </aside>

          <main class="flex-1 p-6 overflow-auto">
            <div class="bg-white rounded-lg shadow-sm p-6">
              {props.children}
            </div>
          </main>
        </div>
      )}
    </>
  );
};

export default AdminLayout;
