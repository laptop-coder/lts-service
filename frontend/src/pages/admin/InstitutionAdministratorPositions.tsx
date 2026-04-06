import { createSignal, onMount, For, Show } from "solid-js";
import { api } from "../../lib/api";
import { PERMISSIONS } from "../../lib/permissions";
import { usePermissions } from "../../lib/permissions";
import type { InstitutionAdministratorPosition } from "../../lib/types";

const InstitutionAdministratorPositions = () => {
  const [
    institutionAdministratorPositions,
    setInstitutionAdministratorPositions,
  ] = createSignal<InstitutionAdministratorPosition[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");
  const [
    newInstitutionAdministratorPositionName,
    setNewInstitutionAdministratorPositionName,
  ] = createSignal("");
  const [creating, setCreating] = createSignal(false);
  const [deletingId, setDeletingId] = createSignal<number | null>(null);

  const { hasPermission } = usePermissions();

  const loadInstitutionAdministratorPositions = async () => {
    try {
      const data = await api.get<{
        institutionAdministratorPositions: InstitutionAdministratorPosition[];
      }>("/institution_administrators/positions");
      setInstitutionAdministratorPositions(
        data.institutionAdministratorPositions,
      );
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

  const createInstitutionAdministratorPosition = async (e: Event) => {
    e.preventDefault();
    if (!newInstitutionAdministratorPositionName().trim()) return;

    setCreating(true);
    try {
      const formData = new URLSearchParams();
      formData.append("name", newInstitutionAdministratorPositionName().trim());

      await api.post("/institution_administrators/positions", formData);
      setNewInstitutionAdministratorPositionName("");
      await loadInstitutionAdministratorPositions();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка создания должности",
      );
    } finally {
      setCreating(false);
    }
  };

  const deleteInstitutionAdministratorPosition = async (id: number) => {
    if (!confirm("Удалить должность? Это действие нельзя отменить.")) return;

    setDeletingId(id);
    try {
      await api.delete(`/institution_administrators/positions/${id}`);
      await loadInstitutionAdministratorPositions();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка удаления должности",
      );
    } finally {
      setDeletingId(null);
    }
  };

  onMount(() => {
    loadInstitutionAdministratorPositions();
  });

  return (
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-bold">Управление должностями администрации</h1>
      </div>

      <Show when={error()}>
        <div class="bg-red-100 text-red-700 p-3 rounded-lg">{error()}</div>
      </Show>

      {/* Form for positions creating */}
      <Show
        when={hasPermission(
          PERMISSIONS.POSITION_INSTITUTION_ADMINISTRATOR_CREATE,
        )}
      >
        <form
          onSubmit={createInstitutionAdministratorPosition}
          class="bg-gray-50 p-4 rounded-lg space-y-3"
        >
          <h2 class="font-semibold">Создать новую должность</h2>
          <div class="flex gap-2">
            <input
              type="text"
              value={newInstitutionAdministratorPositionName()}
              onInput={(e) =>
                setNewInstitutionAdministratorPositionName(
                  e.currentTarget.value,
                )
              }
              placeholder="Название должности"
              class="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={creating()}
            />
            <button
              type="submit"
              disabled={
                creating() || !newInstitutionAdministratorPositionName().trim()
              }
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

      <Show
        when={!loading() && institutionAdministratorPositions().length === 0}
      >
        <div class="text-center text-gray-500 py-8">
          Нет должностей. Создайте первую.
        </div>
      </Show>

      <Show when={!loading() && institutionAdministratorPositions().length > 0}>
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
              <For each={institutionAdministratorPositions()}>
                {(institutionAdministratorPosition) => (
                  <tr class="hover:bg-gray-50">
                    <td class="px-4 py-3 text-sm text-gray-500">
                      {institutionAdministratorPosition.id}
                    </td>
                    <td class="px-4 py-3 font-medium">
                      {institutionAdministratorPosition.name}
                    </td>
                    <td class="px-4 py-3 text-right">
                      <Show
                        when={hasPermission(
                          PERMISSIONS.POSITION_INSTITUTION_ADMINISTRATOR_DELETE,
                        )}
                      >
                        <button
                          onClick={() =>
                            deleteInstitutionAdministratorPosition(
                              institutionAdministratorPosition.id,
                            )
                          }
                          disabled={
                            deletingId() === institutionAdministratorPosition.id
                          }
                          class="text-red-600 hover:text-red-800 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed"
                        >
                          {deletingId() === institutionAdministratorPosition.id
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

export default InstitutionAdministratorPositions;
