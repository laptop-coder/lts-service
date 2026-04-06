import { createSignal, onMount, For, Show } from "solid-js";
import { api } from "../../lib/api";
import { PERMISSIONS } from "../../lib/permissions";
import { usePermissions } from "../../lib/permissions";
import type { Room } from "../../lib/types";

const Rooms = () => {
  const [rooms, setRooms] = createSignal<Room[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");
  const [newRoomName, setNewRoomName] = createSignal("");
  const [creating, setCreating] = createSignal(false);
  const [deletingId, setDeletingId] = createSignal<number | null>(null);

  const { hasPermission } = usePermissions();

  const loadRooms = async () => {
    try {
      const data = await api.get<{ rooms: Room[] }>("/rooms");
      setRooms(data.rooms);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки кабинетов", // TODO: is it safety to print err message?
      );
    } finally {
      setLoading(false);
    }
  };

  const createRoom = async (e: Event) => {
    e.preventDefault();
    if (!newRoomName().trim()) return;

    setCreating(true);
    try {
      const formData = new URLSearchParams();
      formData.append("name", newRoomName().trim());

      await api.post("/rooms", formData);
      setNewRoomName("");
      await loadRooms();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка создания кабинета");
    } finally {
      setCreating(false);
    }
  };

  const deleteRoom = async (id: number) => {
    if (!confirm("Удалить кабинет? Это действие нельзя отменить.")) return;

    setDeletingId(id);
    try {
      await api.delete(`/rooms/${id}`);
      await loadRooms();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка удаления кабинета");
    } finally {
      setDeletingId(null);
    }
  };

  onMount(() => {
    loadRooms();
  });

  return (
    <div class="space-y-6 p-4">
      <div class="mb-6">
        <h1 class="text-3xl font-bold text-gray-800">Управление кабинетами</h1>
        <p class="text-gray-500 mt-1">Список аудиторий и помещений</p>
      </div>

      <Show when={error()}>
        <div class="bg-red-50 border border-red-200 text-red-600 p-3 rounded-xl">
          {error()}
        </div>
      </Show>

      {/* Form for rooms creating */}
      <Show when={hasPermission(PERMISSIONS.ROOM_CREATE)}>
        <div class="bg-white rounded-2xl shadow-lg p-6">
          <h2 class="text-lg font-semibold text-gray-800 mb-4">
            Добавить кабинет
          </h2>
          <form onSubmit={createRoom} class="flex gap-3">
            <input
              type="text"
              value={newRoomName()}
              onInput={(e) => setNewRoomName(e.currentTarget.value)}
              placeholder="Название кабинета"
              class="flex-1 px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition disabled:opacity-50"
              disabled={creating()}
            />
            <button
              type="submit"
              disabled={creating() || !newRoomName().trim()}
              class="px-5 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed font-medium"
            >
              {creating() ? "Создание..." : "Создать"}
            </button>
          </form>
        </div>
      </Show>

      {/* List of rooms */}
      <Show when={loading()}>
        <div class="flex justify-center py-16">
          <div class="text-gray-500">Загрузка...</div>
        </div>
      </Show>

      <Show when={!loading() && rooms().length === 0}>
        <div class="text-center py-16">
          <div class="text-5xl mb-3">🚪</div>
          <p class="text-gray-500">Нет кабинетов</p>
          <p class="text-gray-400 text-sm mt-1">Создайте первый</p>
        </div>
      </Show>

      <Show when={!loading() && rooms().length > 0}>
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
                <For each={rooms()}>
                  {(room) => (
                    <tr class="hover:bg-gray-50 transition">
                      <td class="px-6 py-4 text-sm text-gray-500 font-mono">
                        {room.id}
                      </td>
                      <td class="px-6 py-4 font-medium text-gray-800">
                        {room.name}
                      </td>
                      <td class="px-6 py-4 text-right">
                        <Show when={hasPermission(PERMISSIONS.ROOM_DELETE)}>
                          <button
                            onClick={() => deleteRoom(room.id)}
                            disabled={deletingId() === room.id}
                            class="text-red-600 hover:text-red-800 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed font-medium"
                          >
                            {deletingId() === room.id
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
    </div>
  );
};

export default Rooms;
