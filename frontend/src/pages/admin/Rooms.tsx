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
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-bold">Управление кабинетами</h1>
      </div>

      <Show when={error()}>
        <div class="bg-red-100 text-red-700 p-3 rounded-lg">{error()}</div>
      </Show>

      {/* Form for rooms creating */}
      <Show when={hasPermission(PERMISSIONS.ROOM_CREATE)}>
        <form onSubmit={createRoom} class="bg-gray-50 p-4 rounded-lg space-y-3">
          <h2 class="font-semibold">Создать новый кабинет</h2>
          <div class="flex gap-2">
            <input
              type="text"
              value={newRoomName()}
              onInput={(e) => setNewRoomName(e.currentTarget.value)}
              placeholder="Название кабинета"
              class="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={creating()}
            />
            <button
              type="submit"
              disabled={creating() || !newRoomName().trim()}
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {creating() ? "Создание..." : "Создать"}
            </button>
          </div>
        </form>
      </Show>

      {/* List of rooms */}
      <Show when={loading()}>
        <div class="text-center py-8">Загрузка...</div>
      </Show>

      <Show when={!loading() && rooms().length === 0}>
        <div class="text-center text-gray-500 py-8">
          Нет кабинетов. Создайте первый.
        </div>
      </Show>

      <Show when={!loading() && rooms().length > 0}>
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
              <For each={rooms()}>
                {(room) => (
                  <tr class="hover:bg-gray-50">
                    <td class="px-4 py-3 text-sm text-gray-500">{room.id}</td>
                    <td class="px-4 py-3 font-medium">{room.name}</td>
                    <td class="px-4 py-3 text-right">
                      <Show when={hasPermission(PERMISSIONS.ROOM_DELETE)}>
                        <button
                          onClick={() => deleteRoom(room.id)}
                          disabled={deletingId() === room.id}
                          class="text-red-600 hover:text-red-800 disabled:opacity-50"
                        >
                          {deletingId() === room.id ? "Удаление..." : "Удалить"}
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

export default Rooms;
