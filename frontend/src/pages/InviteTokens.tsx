import { createSignal, Show, For } from "solid-js";
import { api } from "../lib/api";
import { PERMISSIONS } from "../lib/permissions";
import { usePermissions } from "../lib/permissions";
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

  const { hasPermission } = usePermissions();

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
      name: "Учитель",
      permission: PERMISSIONS.TOKEN_INVITE_USER_CREATE,
    },
    {
      id: 6,
      name: "Родитель",
      permission: PERMISSIONS.TOKEN_INVITE_USER_CREATE,
    },
    { id: 7, name: "Ученик", permission: PERMISSIONS.TOKEN_INVITE_USER_CREATE },
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
    <div class="space-y-6">
      <h1 class="text-2xl font-bold">Создание инвайт-токенов</h1>

      <Show when={error()}>
        <div class="bg-red-100 text-red-700 p-3 rounded-lg">{error()}</div>
      </Show>

      <div class="bg-white rounded-lg shadow p-6 space-y-4 max-w-md">
        {/* Count */}
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
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
            class="w-24 px-3 py-2 border rounded-md"
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
                <label class="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={selectedRoles().includes(role.id)}
                    onChange={() => toggleRole(role.id)}
                  />
                  <span>{role.name}</span>
                </label>
              )}
            </For>
          </div>
        </div>

        <button
          onClick={handleCreate}
          disabled={creating()}
          class="w-full py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed"
        >
          {creating()
            ? `Создание... ${progress().current}/${progress().total}`
            : "Создать токены"}
        </button>
      </div>

      {/* List of created tokens */}
      <Show when={tokens().length > 0}>
        <button
          onClick={downloadTokensFile}
          class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed"
        >
          Скачать Markdown
        </button>
        <div class="bg-white rounded-lg shadow p-6 space-y-3">
          <h2 class="text-lg font-semibold">Созданные токены</h2>
          <div class="space-y-2 max-h-96 overflow-y-auto">
            <For each={tokens()}>
              {(item, index) => (
                <div class="flex justify-between items-center p-2 bg-gray-50 rounded">
                  <code class="text-sm font-mono break-all flex-1">
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
                      class="px-3 py-1 text-sm bg-gray-200 rounded hover:bg-gray-300 transition cursor-pointer"
                    >
                      {buttonCopiedIndex() === index()
                        ? "Скопировано!"
                        : "Копировать"}
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
      </Show>
    </div>
  );
};

export default InviteTokens;
