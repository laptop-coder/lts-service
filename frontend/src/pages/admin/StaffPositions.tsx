import { createSignal, createEffect, onMount, For, Show } from "solid-js";
import { api } from "../../lib/api";
import { PERMISSIONS } from "../../lib/permissions";
import { usePermissions } from "../../lib/permissions";
import type { StaffPosition } from "../../lib/types";
import Pagination from "../../components/Pagination";

const StaffPositions = () => {
  const [staffPositions, setStaffPositions] = createSignal<StaffPosition[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");
  const [newStaffPositionName, setNewStaffPositionName] = createSignal("");
  const [creating, setCreating] = createSignal(false);
  const [deletingId, setDeletingId] = createSignal<number | null>(null);
  const [page, setPage] = createSignal(0);
  const [hasMore, setHasMore] = createSignal(true);

  const { hasPermission } = usePermissions();

  let inputRef: HTMLInputElement | undefined;
  const focusInput = () => {
    if (inputRef) {
      inputRef.focus();
    }
  };

  createEffect(() => {
    focusInput();
  });

  const limit = 30;

  createEffect(() => {
    page();
    loadStaffPositions();
  });

  const loadStaffPositions = async () => {
    try {
      const data = await api.get<{ staffPositions: StaffPosition[] }>(
        `/staff/positions?limit=${limit + 1}&offset=${page() * limit}`,
      );
      setHasMore(data.staffPositions.length > limit);
      setStaffPositions(data.staffPositions.slice(0, limit));
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Ошибка загрузки списка должностей",
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
      setError(
        err instanceof Error ? err.message : "Ошибка создания должности",
      );
    } finally {
      setCreating(false);
      focusInput();
    }
  };

  const deleteStaffPosition = async (id: number) => {
    if (!confirm("Удалить должность? Это действие нельзя отменить.")) return;

    setDeletingId(id);
    try {
      await api.delete(`/staff/positions/${id}`);
      await loadStaffPositions();
      if (staffPositions().length === 0 && page() > 0) {
        setPage((prev) => prev - 1);
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка удаления должности",
      );
    } finally {
      setDeletingId(null);
      focusInput();
    }
  };

  onMount(() => {
    loadStaffPositions();
  });

  return (
    <div class="space-y-6 p-4">
      <div class="mb-6">
        <h1 class="text-3xl font-bold text-gray-800">
          Должности сотрудников ОУ
        </h1>
        <p class="text-gray-500 mt-1">Управление должностями сотрудников</p>
      </div>

      <Show when={error()}>
        <div class="bg-red-50 border border-red-200 text-red-600 p-3 rounded-xl">
          {error()}
        </div>
      </Show>

      {/* Form for positions creating */}
      <Show when={hasPermission(PERMISSIONS.POSITION_STAFF_CREATE)}>
        <div class="bg-white rounded-2xl shadow-lg p-6">
          <h2 class="text-lg font-semibold text-gray-800 mb-4">
            Добавить должность
          </h2>
          <form onSubmit={createStaffPosition} class="flex gap-3">
            <input
              ref={inputRef}
              type="text"
              value={newStaffPositionName()}
              onInput={(e) => setNewStaffPositionName(e.currentTarget.value)}
              placeholder="Название должности"
              class="flex-1 px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition disabled:opacity-50"
              disabled={creating()}
            />
            <button
              type="submit"
              disabled={creating() || !newStaffPositionName().trim()}
              class="px-5 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed font-medium"
            >
              {creating() ? "Создание..." : "Создать"}
            </button>
          </form>
        </div>
      </Show>

      {/* List of positions */}
      <Show when={loading()}>
        <div class="flex justify-center py-16">
          <div class="text-gray-500">Загрузка...</div>
        </div>
      </Show>

      <Show when={!loading() && staffPositions().length === 0}>
        <div class="text-center py-16">
          <div class="text-5xl mb-3">🧑‍💼</div>
          <p class="text-gray-500">Нет должностей</p>
          <p class="text-gray-400 text-sm mt-1">Создайте первую</p>
        </div>
      </Show>

      <Show when={!loading() && staffPositions().length > 0}>
        <div class="bg-white rounded-2xl shadow-lg overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th class="px-6 py-4 text-left text-sm font-semibold text-gray-600">
                    ID
                  </th>
                  <th class="px-6 py-4 text-left text-sm font-semibold text-gray-600">
                    Название
                  </th>
                  <th class="px-6 py-4 text-right text-sm font-semibold text-gray-600">
                    Действия
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <For each={staffPositions()}>
                  {(position) => (
                    <tr class="hover:bg-gray-50 transition">
                      <td class="px-6 py-4 text-sm text-gray-500 font-mono">
                        {position.id}
                      </td>
                      <td class="px-6 py-4 font-medium text-gray-800">
                        {position.name}
                      </td>
                      <td class="px-6 py-4 text-right">
                        <Show
                          when={hasPermission(
                            PERMISSIONS.POSITION_STAFF_DELETE,
                          )}
                        >
                          <button
                            onClick={() => deleteStaffPosition(position.id)}
                            disabled={deletingId() === position.id}
                            class="text-red-600 hover:text-red-800 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed font-medium"
                          >
                            {deletingId() === position.id
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
        </div>
      </Show>
      <Pagination
        page={page()}
        hasMore={hasMore()}
        onPrev={() => setPage((prev) => prev - 1)}
        onNext={() => setPage((prev) => prev + 1)}
      />
    </div>
  );
};

export default StaffPositions;
