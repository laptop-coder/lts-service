import { createSignal, Show, For } from "solid-js";
import { api } from "../lib/api";
import { usePermissions, PERMISSIONS, ROLES } from "../lib/permissions";
import QRCodeButton from "../components/QRCode";

const InviteTokens = () => {
  const [count, setCount] = createSignal(1);
  const [selectedRoles, setSelectedRoles] = createSignal<number[]>([]);
  const [tokens, setTokens] = createSignal<{ token: string; index: number }[]>(
    [],
  );
  const [creating, setCreating] = createSignal(false);
  const [error, setError] = createSignal("");
  const [progress, setProgress] = createSignal({ current: 0, total: 0 });
  const [buttonCopiedIndex, setButtonCopiedIndex] = createSignal<number | null>(
    null,
  );

  const { hasPermission, hasRole } = usePermissions();

  const downloadTokensFile = () => {
    if (tokens().length === 0 || selectedRoles().length === 0) return;
    // Assemble content
    const content = `# Индивидуальные пригласительные ссылки для регистрации аккаунтов\n\nРоли:\n${selectedRoles()
      .map((id) => `- ${roles.find((r) => r.id === id)!.name}`)
      .join("\n")}\n\n${tokens()
      .map(
        (item) =>
          `8< ----------------------\n\n${window.location.protocol}//${window.location.host}/register?inviteToken=${item.token}\n`,
      )
      .join("\n")}`;
    // Create file in browser
    const blob = new Blob([content], { type: "text/plain;charset=utf-8" });
    // Create download link and click
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `invite-tokens_${new Date().toISOString().slice(0, 19)}.md`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const delay = (ms: number) =>
    new Promise((resolve) => setTimeout(resolve, ms));
  const DELAY_MS = 500;

  const roles = [
    {
      id: 2,
      name: "Администратор",
      permission: PERMISSIONS.TOKEN_INVITE_ADMIN_CREATE,
    },
    {
      id: 3,
      name: "Администрация ОУ",
      permission: PERMISSIONS.TOKEN_INVITE_USER_CREATE,
    },
    {
      id: 4,
      name: "Сотрудник ОУ",
      permission: PERMISSIONS.TOKEN_INVITE_USER_CREATE,
    },
    {
      id: 5,
      name: "Преподаватель",
      permission: PERMISSIONS.TOKEN_INVITE_USER_CREATE,
    },
    {
      id: 6,
      name: "Родитель",
      permission: PERMISSIONS.TOKEN_INVITE_USER_CREATE,
    },
    {
      id: 7,
      name: "Обучающийся",
      permission: PERMISSIONS.TOKEN_INVITE_USER_CREATE,
    },
  ];

  const availableRoles = roles.filter((role) => hasPermission(role.permission));

  const toggleRole = (roleId: number) => {
    if (selectedRoles().includes(roleId)) {
      setSelectedRoles(selectedRoles().filter((id) => id !== roleId));
    } else {
      setSelectedRoles([...selectedRoles(), roleId]);
    }
  };

  const createToken = async (roleIds: number[]): Promise<string> => {
    const formData = new URLSearchParams();
    roleIds.forEach((id) => formData.append("roleId", id.toString()));
    const data = await api.post<{ inviteToken: string }>(
      "/tokens/invite",
      formData,
    );
    return data.inviteToken;
  };

  const handleCreate = async () => {
    if (hasRole(ROLES.SUPERADMIN)) {
      setSelectedRoles([2])
    }

    if (selectedRoles().length === 0) {
      setError("Выберите хотя бы одну роль");
      return;
    }

    setCreating(true);
    setError("");
    setTokens([]);
    setProgress({ current: 0, total: count() });

    const results: { token: string; index: number }[] = [];
    for (let i = 0; i < count(); i++) {
      try {
        await delay(DELAY_MS);
        const token = await createToken(selectedRoles());
        results.push({ token, index: i + 1 });
        setProgress({ current: i + 1, total: count() });
      } catch (err) {
        setError(
          `Ошибка при создании токена ${i + 1}: ${err instanceof Error ? err.message : err}`,
        );
        break;
      }
    }
    setTokens(results);
    setCreating(false);
  };

  const copyToClipboard = async (text: string, index: number) => {
    await navigator.clipboard.writeText(text);
    setButtonCopiedIndex(index);
    setTimeout(() => setButtonCopiedIndex(null), 2000);
  };

  return (
    <div class="space-y-6 p-4">
      <div class="mb-6">
        <h1 class="text-3xl font-bold text-gray-800">Инвайт-токены</h1>
        <p class="text-gray-500 mt-1">
          Создание пригласительных ссылок для регистрации{" "}
          {hasRole(ROLES.SUPERADMIN) && "админов"}
        </p>
      </div>

      <Show when={error()}>
        <div class="bg-red-50 border border-red-200 text-red-600 p-3 rounded-xl">
          {error()}
        </div>
      </Show>

      <div class="bg-white rounded-2xl shadow-lg p-6 max-w-md">
        <h2 class="text-lg font-semibold text-gray-800 mb-4">
          Параметры токенов
        </h2>

        <div class="space-y-4">
          {/* Count */}
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Количество токенов
            </label>
            <input
              type="number"
              min="1"
              max="500"
              value={count()}
              onInput={(e) =>
                setCount(
                  Math.min(
                    500,
                    Math.max(1, parseInt(e.currentTarget.value) || 1),
                  ),
                )
              }
              class="w-32 px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
            />
          </div>

          {/* Roles */}
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Роли
            </label>
            <div class="space-y-2">
              <For each={availableRoles}>
                {(role) => (
                  <label
                    class={`flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 transition cursor-pointer ${((hasRole(ROLES.SUPERADMIN) && role.id === 2) || creating()) && "cursor-not-allowed"}`}
                  >
                    <input
                      type="checkbox"
                      checked={
                        selectedRoles().includes(role.id) ||
                        (hasRole(ROLES.SUPERADMIN) && role.id === 2)
                      }
                      onChange={() => toggleRole(role.id)}
                      class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500 cursor-pointer disabled:cursor-not-allowed"
                      disabled={
                        (hasRole(ROLES.SUPERADMIN) && role.id === 2) ||
                        creating()
                      }
                    />
                    <span class="text-gray-700">{role.name}</span>
                  </label>
                )}
              </For>
            </div>
          </div>

          <button
            onClick={handleCreate}
            disabled={creating()}
            class="w-full py-2.5 bg-blue-600 text-white rounded-xl hover:bg-blue-700 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed font-medium"
          >
            {creating()
              ? `Создание... ${progress().current}/${progress().total}`
              : "Создать токены"}
          </button>
        </div>
      </div>

      {/* List of created tokens */}
      <Show when={tokens().length > 0}>
        <div class="space-y-4">
          <div class="flex justify-between items-center">
            <h2 class="text-xl font-semibold text-gray-800">
              Созданные токены
            </h2>
            <button
              onClick={downloadTokensFile}
              class="px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition cursor-pointer font-medium"
            >
              📥 Скачать Markdown
            </button>
          </div>

          <div class="bg-white rounded-2xl shadow-lg overflow-hidden">
            <div class="divide-y divide-gray-100">
              <For each={tokens()}>
                {(item, index) => (
                  <div class="p-4 hover:bg-gray-50 transition flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                    <code class="text-sm font-mono break-all text-gray-600 bg-gray-50 px-3 py-1.5 rounded-lg">
                      {item.token}
                    </code>
                    <div class="flex gap-2">
                      <button
                        onClick={() =>
                          copyToClipboard(
                            `${window.location.protocol}//${window.location.host}/register?inviteToken=${item.token}`,
                            index(),
                          )
                        }
                        class="w-28 px-3 py-1.5 text-sm bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition cursor-pointer font-medium"
                      >
                        {buttonCopiedIndex() === index()
                          ? "✓ Скопировано"
                          : "📋 Копировать"}
                      </button>
                      <QRCodeButton
                        text={`${window.location.protocol}//${window.location.host}/register?inviteToken=${item.token}`}
                      />
                    </div>
                  </div>
                )}
              </For>
            </div>
          </div>
        </div>
      </Show>
    </div>
  );
};

export default InviteTokens;
