import { createSignal, onMount, For, Show } from "solid-js";
import { api } from "../../lib/api";
import { PERMISSIONS } from "../../lib/permissions";
import { usePermissions } from "../../lib/permissions";
import type { StaffPosition } from "../../lib/types";

const StaffPositions = () => {
  const [staffPositions, setStaffPositions] = createSignal<StaffPosition[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");
  const [newStaffPositionName, setNewStaffPositionName] = createSignal("");
  const [creating, setCreating] = createSignal(false);
  const [deletingId, setDeletingId] = createSignal<number | null>(null);

  const { hasPermission } = usePermissions();

  const loadStaffPositions = async () => {
    try {
      const data = await api.get<{ staffPositions: StaffPosition[] }>("/staff/positions");
      setStaffPositions(data.staffPositions);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки списка должностей",
      );
    } finally {
      setLoading(false);
    }
  };

  const createStaffPosition = async (e: Event) => {
    e.preventDefault();
    if (!newStaffPositionName().trim()) return;

    setCreating(true);
    try {
      const formData = new URLSearchParams();
      formData.append("name", newStaffPositionName().trim());

      await api.post("/staff/positions", formData);
      setNewStaffPositionName("");
      await loadStaffPositions();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка создания должности");
    } finally {
      setCreating(false);
    }
  };

  const deleteStaffPosition = async (id: number) => {
    if (!confirm("Удалить должность? Это действие нельзя отменить.")) return;

    setDeletingId(id);
    try {
      await api.delete(`/staff/positions/${id}`);
      await loadStaffPositions();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка удаления должности");
    } finally {
      setDeletingId(null);
    }
  };

  onMount(() => {
    loadStaffPositions();
  });

  return (
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-bold">Управление должностями сотрудников</h1>
      </div>

      <Show when={error()}>
        <div class="bg-red-100 text-red-700 p-3 rounded-lg">{error()}</div>
      </Show>

      {/* Form for positions creating */}
      <Show when={hasPermission(PERMISSIONS.POSITION_STAFF_CREATE)}>
        <form
          onSubmit={createStaffPosition}
          class="bg-gray-50 p-4 rounded-lg space-y-3"
        >
          <h2 class="font-semibold">Создать новую должность</h2>
          <div class="flex gap-2">
            <input
              type="text"
              value={newStaffPositionName()}
              onInput={(e) => setNewStaffPositionName(e.currentTarget.value)}
              placeholder="Название должности"
              class="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={creating()}
            />
            <button
              type="submit"
              disabled={creating() || !newStaffPositionName().trim()}
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed"
            >
              {creating() ? "Создание..." : "Создать"}
            </button>
          </div>
        </form>
      </Show>

      {/* List of positions */}
      <Show when={loading()}>
        <div class="text-center py-8">Загрузка...</div>
      </Show>

      <Show when={!loading() && staffPositions().length === 0}>
        <div class="text-center text-gray-500 py-8">
          Нет должностей. Создайте первую.
        </div>
      </Show>

      <Show when={!loading() && staffPositions().length > 0}>
        <div class="bg-white rounded-lg shadow overflow-hidden">
          <table class="w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">
                  ID
                </th>
                <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">
                  Название
                </th>
                <th class="px-4 py-3 text-right text-sm font-medium text-gray-500">
                  Действия
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <For each={staffPositions()}>
                {(staffPosition) => (
                  <tr class="hover:bg-gray-50">
                    <td class="px-4 py-3 text-sm text-gray-500">
                      {staffPosition.id}
                    </td>
                    <td class="px-4 py-3 font-medium">{staffPosition.name}</td>
                    <td class="px-4 py-3 text-right">
                      <Show when={hasPermission(PERMISSIONS.POSITION_STAFF_DELETE)}>
                        <button
                          onClick={() => deleteStaffPosition(staffPosition.id)}
                          disabled={deletingId() === staffPosition.id}
                          class="text-red-600 hover:text-red-800 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed"
                        >
                          {deletingId() === staffPosition.id
                            ? "Удаление..."
                            : "Удалить"}
                        </button>
                      </Show>
                    </td>
                  </tr>
                )}
              </For>
            </tbody>
          </table>
        </div>
      </Show>
    </div>
  );
};

export default StaffPositions;
